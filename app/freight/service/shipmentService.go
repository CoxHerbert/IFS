package service

import (
	"baize/app/freight/dao"
	"baize/app/freight/models"
	"baize/app/utils/snowflake"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"
)

var shipmentServiceImpl *shipmentService

func init() {
	shipmentServiceImpl = &shipmentService{shipmentDao: dao.GetShipmentDao()}
}

type shipmentService struct {
	shipmentDao interface {
		InsertShipment(plan *models.ShipmentPlanDML, cargoList []*models.CargoDML, containers []*models.ContainerPlanDML)
		SelectShipmentList(query *models.ShipmentPlanDQL) (list []*models.ShipmentPlanVo, total *int64)
		SelectShipmentById(shipmentId int64) *models.ShipmentPlanVo
		SelectShipmentByToken(token string) *models.ShipmentPlanVo
		SelectCargoList(shipmentId int64) []*models.CargoVo
		SelectContainerList(shipmentId int64) []*models.ContainerPlanVo
		SelectOrderByShipmentId(shipmentId int64) *models.ShipmentOrderVo
		UpdateShipmentStatus(update *models.ShipmentStatusUpdateDML)
		InsertShipmentOrder(order *models.ShipmentOrderDML)
		DeleteShipmentByIds(shipmentIds []int64)
	}
}

type containerCapacity struct {
	Type      string
	VolumeCbm float64
	WeightKg  float64
}

var capacities = []containerCapacity{
	{Type: "20GP", VolumeCbm: 28, WeightKg: 21700},
	{Type: "40GP", VolumeCbm: 58, WeightKg: 26500},
	{Type: "40HQ", VolumeCbm: 68, WeightKg: 26500},
}

var shipmentStatuses = []*models.ShipmentStatusStep{
	{Value: "10", Label: "计划已创建"},
	{Value: "20", Label: "出货计划已确认"},
	{Value: "30", Label: "等待客户发货"},
	{Value: "40", Label: "已提货/已送仓"},
	{Value: "50", Label: "仓库已收货"},
	{Value: "60", Label: "已入仓/码头进仓"},
	{Value: "70", Label: "订舱处理中"},
	{Value: "80", Label: "舱位已确认"},
	{Value: "90", Label: "报关资料已收齐"},
	{Value: "100", Label: "报关已放行"},
	{Value: "110", Label: "已装柜"},
	{Value: "120", Label: "已进港/码头放行"},
	{Value: "130", Label: "船舶已开航"},
	{Value: "140", Label: "目的港已到港"},
	{Value: "150", Label: "目的港清关中"},
	{Value: "160", Label: "目的港已清关"},
	{Value: "170", Label: "已派送/已签收"},
	{Value: "900", Label: "异常处理中"},
}

func GetShipmentService() *shipmentService {
	return shipmentServiceImpl
}

func (service *shipmentService) ImportShipment(req *models.ShipmentImportReq, username string) (*models.ShipmentDetailVo, error) {
	if req.CustomerId == 0 || len(req.CargoList) == 0 {
		return nil, errors.New("请选择客户并导入货物明细")
	}
	shipmentId := snowflake.GenID()
	plan := &models.ShipmentPlanDML{
		ShipmentId:   shipmentId,
		ShipmentNo:   fmt.Sprintf("SP%s%d", time.Now().Format("20060102"), shipmentId%1000000),
		OrderNo:      req.OrderNo,
		CustomerId:   req.CustomerId,
		CustomerName: req.CustomerName,
		Pol:          req.Pol,
		Pod:          req.Pod,
		PlannedEtd:   req.PlannedEtd,
		PlannedEta:   req.PlannedEta,
		Status:       "10",
		ShareToken:   genShareToken(),
		Remark:       req.Remark,
		CreateBy:     username,
		UpdateBy:     username,
	}

	cargoList := make([]*models.CargoDML, 0, len(req.CargoList))
	for _, item := range req.CargoList {
		if item == nil || item.CargoName == "" {
			continue
		}
		volume := item.VolumeCbm
		if volume == 0 && item.LengthCm > 0 && item.WidthCm > 0 && item.HeightCm > 0 && item.Cartons > 0 {
			volume = item.LengthCm * item.WidthCm * item.HeightCm / 1000000 * float64(item.Cartons)
		}
		cargo := &models.CargoDML{
			CargoId:     snowflake.GenID(),
			ShipmentId:  shipmentId,
			Sku:         item.Sku,
			CargoName:   item.CargoName,
			PackageType: item.PackageType,
			Quantity:    item.Quantity,
			Cartons:     item.Cartons,
			WeightKg:    item.WeightKg,
			VolumeCbm:   round2(volume),
			LengthCm:    item.LengthCm,
			WidthCm:     item.WidthCm,
			HeightCm:    item.HeightCm,
		}
		plan.TotalWeight += cargo.WeightKg
		plan.TotalVolume += cargo.VolumeCbm
		plan.TotalCartons += cargo.Cartons
		cargoList = append(cargoList, cargo)
	}
	if len(cargoList) == 0 {
		return nil, errors.New("货物明细不能为空")
	}
	plan.TotalWeight = round2(plan.TotalWeight)
	plan.TotalVolume = round2(plan.TotalVolume)
	containers := service.RecommendContainers(shipmentId, plan.TotalVolume, plan.TotalWeight, req.PreferredType)
	service.shipmentDao.InsertShipment(plan, cargoList, containers)
	return service.SelectShipmentDetail(shipmentId), nil
}

func (service *shipmentService) RecommendContainers(shipmentId int64, totalVolume, totalWeight float64, preferredType string) []*models.ContainerPlanDML {
	capacity := capacities[len(capacities)-1]
	for _, item := range capacities {
		if preferredType == item.Type {
			capacity = item
			break
		}
	}
	if preferredType == "" {
		for _, item := range capacities {
			if totalVolume <= item.VolumeCbm && totalWeight <= item.WeightKg {
				capacity = item
				break
			}
		}
	}
	quantity := int64(math.Ceil(math.Max(totalVolume/capacity.VolumeCbm, totalWeight/capacity.WeightKg)))
	if quantity < 1 {
		quantity = 1
	}
	maxVolume := capacity.VolumeCbm * float64(quantity)
	maxWeight := capacity.WeightKg * float64(quantity)
	loadRate := math.Max(totalVolume/maxVolume, totalWeight/maxWeight) * 100
	return []*models.ContainerPlanDML{{
		ContainerPlanId: snowflake.GenID(),
		ShipmentId:      shipmentId,
		ContainerType:   capacity.Type,
		Quantity:        quantity,
		MaxVolume:       round2(maxVolume),
		MaxWeight:       round2(maxWeight),
		UsedVolume:      round2(totalVolume),
		UsedWeight:      round2(totalWeight),
		LoadRate:        round2(loadRate),
		Remark:          "系统按体积/重量瓶颈自动推荐",
	}}
}

func (service *shipmentService) SelectShipmentList(query *models.ShipmentPlanDQL) (list []*models.ShipmentPlanVo, total *int64) {
	return service.shipmentDao.SelectShipmentList(query)
}

func (service *shipmentService) SelectShipmentDetail(shipmentId int64) *models.ShipmentDetailVo {
	plan := service.shipmentDao.SelectShipmentById(shipmentId)
	if plan == nil {
		return nil
	}
	return service.buildDetail(plan)
}

func (service *shipmentService) SelectShareDetail(token string) *models.ShipmentDetailVo {
	plan := service.shipmentDao.SelectShipmentByToken(token)
	if plan == nil {
		return nil
	}
	return service.buildDetail(plan)
}

func (service *shipmentService) UpdateShipmentStatus(shipmentId int64, req *models.ShipmentStatusUpdateReq, username string) error {
	if !validStatus(req.Status) {
		return errors.New("出货状态不正确")
	}
	if service.shipmentDao.SelectShipmentById(shipmentId) == nil {
		return errors.New("出货计划不存在")
	}
	service.shipmentDao.UpdateShipmentStatus(&models.ShipmentStatusUpdateDML{
		ShipmentId: shipmentId,
		Status:     req.Status,
		ActualEtd:  req.ActualEtd,
		ActualEta:  req.ActualEta,
		Remark:     req.Remark,
		UpdateBy:   username,
	})
	return nil
}

func (service *shipmentService) ConfirmShipment(shipmentId int64, username string) (*models.ShipmentOrderVo, error) {
	plan := service.shipmentDao.SelectShipmentById(shipmentId)
	if plan == nil {
		return nil, errors.New("出货计划不存在")
	}
	orderStatus := plan.Status
	if shouldPromoteToConfirmed(plan.Status) {
		orderStatus = "20"
		service.shipmentDao.UpdateShipmentStatus(&models.ShipmentStatusUpdateDML{
			ShipmentId: shipmentId,
			Status:     orderStatus,
			UpdateBy:   username,
		})
	}
	if order := service.shipmentDao.SelectOrderByShipmentId(shipmentId); order != nil {
		return order, nil
	}
	order := &models.ShipmentOrderDML{
		OrderId:    snowflake.GenID(),
		ShipmentId: shipmentId,
		OrderNo:    fmt.Sprintf("SO%s%d", time.Now().Format("20060102"), shipmentId%1000000),
		Status:     orderStatus,
		CreateBy:   username,
		UpdateBy:   username,
	}
	service.shipmentDao.InsertShipmentOrder(order)
	return service.shipmentDao.SelectOrderByShipmentId(shipmentId), nil
}

func (service *shipmentService) DeleteShipmentByIds(shipmentIds []int64) {
	service.shipmentDao.DeleteShipmentByIds(shipmentIds)
}

func (service *shipmentService) buildDetail(plan *models.ShipmentPlanVo) *models.ShipmentDetailVo {
	return &models.ShipmentDetailVo{
		Plan:       plan,
		CargoList:  service.shipmentDao.SelectCargoList(plan.ShipmentId),
		Containers: service.shipmentDao.SelectContainerList(plan.ShipmentId),
		Order:      service.shipmentDao.SelectOrderByShipmentId(plan.ShipmentId),
		StatusFlow: buildStatusFlow(plan.Status),
	}
}

func buildStatusFlow(current string) []*models.ShipmentStatusStep {
	flow := make([]*models.ShipmentStatusStep, 0, len(shipmentStatuses))
	if current == "900" {
		for _, status := range shipmentStatuses {
			if status.Value != "900" {
				continue
			}
			step := *status
			step.Active = true
			flow = append(flow, &step)
		}
		return flow
	}
	currentValue, _ := strconv.Atoi(current)
	for _, status := range shipmentStatuses {
		if status.Value == "900" {
			continue
		}
		statusValue, _ := strconv.Atoi(status.Value)
		step := *status
		step.Active = statusValue <= currentValue
		flow = append(flow, &step)
	}
	return flow
}

func validStatus(status string) bool {
	for _, item := range shipmentStatuses {
		if item.Value == status {
			return true
		}
	}
	return false
}

func genShareToken() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%d", snowflake.GenID())
	}
	return hex.EncodeToString(buf)
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}

func shouldPromoteToConfirmed(status string) bool {
	if status == "" || status == "900" {
		return false
	}
	currentValue, err := strconv.Atoi(status)
	if err != nil {
		return false
	}
	return currentValue < 20
}
