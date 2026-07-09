package service

import (
	"baize/app/constant/dataScopeAspect"
	customerModels "baize/app/customer/models"
	customerService "baize/app/customer/service"
	"baize/app/freight/dao"
	"baize/app/freight/models"
	notificationService "baize/app/notification/service"
	"baize/app/system/models/loginModels"
	"baize/app/utils/admin"
	"baize/app/utils/snowflake"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
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
		UpdateShipmentCustomer(shipmentId int64, customerId int64, customerName string, salesUserId int64, salesUserName string, updateBy string)
		InsertShipmentOrder(order *models.ShipmentOrderDML)
		DeleteShipmentByIds(shipmentIds []int64)
	}
}

type containerCapacity struct {
	Type      string
	VolumeCbm float64
	SafeCbm   float64
	WeightKg  float64
	LengthCm  float64
	WidthCm   float64
	HeightCm  float64
}

var capacities = []containerCapacity{
	{Type: "20GP", VolumeCbm: 33.2, SafeCbm: 28, WeightKg: 21700, LengthCm: 589, WidthCm: 235, HeightCm: 239},
	{Type: "40GP", VolumeCbm: 67.7, SafeCbm: 58, WeightKg: 26500, LengthCm: 1203, WidthCm: 235, HeightCm: 239},
	{Type: "40HQ", VolumeCbm: 76.4, SafeCbm: 68, WeightKg: 26500, LengthCm: 1203, WidthCm: 235, HeightCm: 269},
}

var placementColors = []string{"#2563eb", "#16a34a", "#f97316", "#9333ea", "#0f766e", "#dc2626", "#ca8a04", "#0891b2"}

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
	{Value: "130", Label: "船舶已开船"},
	{Value: "140", Label: "目的港已到港"},
	{Value: "150", Label: "目的港清关中"},
	{Value: "160", Label: "目的港已清关"},
	{Value: "170", Label: "已派送/已签收"},
	{Value: "900", Label: "异常处理中"},
}

func GetShipmentService() *shipmentService {
	return shipmentServiceImpl
}

func (service *shipmentService) ImportShipment(req *models.ShipmentImportReq, username string, operatorUserId int64) (*models.ShipmentDetailVo, error) {
	if len(req.CargoList) == 0 {
		return nil, errors.New("请导入货物明细")
	}
	shipmentId := snowflake.GenID()
	customer := customerService.GetCustomerService().SelectCustomerById(req.CustomerId)
	if customer != nil && req.CustomerName == "" {
		req.CustomerName = customer.CustomerName
	}
	plan := &models.ShipmentPlanDML{
		ShipmentId:    shipmentId,
		ShipmentNo:    fmt.Sprintf("SP%s%d", time.Now().Format("20060102"), shipmentId%1000000),
		OrderNo:       req.OrderNo,
		CustomerId:    req.CustomerId,
		CustomerName:  req.CustomerName,
		SalesUserId:   customerSalesUserId(customer),
		SalesUserName: customerSalesUserName(customer),
		Pol:           req.Pol,
		Pod:           req.Pod,
		PlannedEtd:    req.PlannedEtd,
		PlannedEta:    req.PlannedEta,
		Status:        "10",
		ShareToken:    genShareToken(),
		Remark:        req.Remark,
		CreateBy:      username,
		UpdateBy:      username,
	}

	cargoItems, summary, err := service.normalizeCargoList(req.CargoList)
	if err != nil {
		return nil, err
	}

	cargoList := make([]*models.CargoDML, 0, len(cargoItems))
	for _, item := range cargoItems {
		cargoList = append(cargoList, &models.CargoDML{
			CargoId:     snowflake.GenID(),
			ShipmentId:  shipmentId,
			Sku:         item.Sku,
			CargoName:   item.CargoName,
			PackageType: item.PackageType,
			Quantity:    item.Quantity,
			Cartons:     item.Cartons,
			WeightKg:    item.WeightKg,
			VolumeCbm:   item.VolumeCbm,
			LengthCm:    item.LengthCm,
			WidthCm:     item.WidthCm,
			HeightCm:    item.HeightCm,
		})
	}

	plan.TotalWeight = summary.TotalWeight
	plan.TotalVolume = summary.TotalVolume
	plan.TotalCartons = summary.TotalCartons
	containers := service.RecommendContainers(shipmentId, plan.TotalVolume, plan.TotalWeight, req.PreferredType)
	service.shipmentDao.InsertShipment(plan, cargoList, containers)
	notificationService.GetNotificationService().NotifyShipmentCreated(plan, username, operatorUserId)
	return service.SelectShipmentDetail(shipmentId), nil
}

func (service *shipmentService) EstimateShipment(req *models.ShipmentEstimateReq) (*models.ShipmentEstimateVo, error) {
	if len(req.CargoList) == 0 {
		return nil, errors.New("货物明细不能为空")
	}
	cargoList, summary, err := service.normalizeCargoList(req.CargoList)
	if err != nil {
		return nil, err
	}
	containers := service.buildContainerPreview(summary.TotalVolume, summary.TotalWeight, req.PreferredType, req, cargoList)
	lcl := buildLclSuggestion(summary.TotalVolume, req.PreferredType, req)
	loadingPlan := buildLoadingPlan(containers)
	return &models.ShipmentEstimateVo{
		Summary:             summary,
		NormalizedCargoList: cargoList,
		Containers:          containers,
		Lcl:                 lcl,
		Recommendation:      buildRecommendation(containers, lcl, req.PreferredType),
		LoadingPlan:         loadingPlan,
		Warnings:            collectEstimateWarnings(containers),
	}, nil
}

func (service *shipmentService) RecommendContainers(shipmentId int64, totalVolume, totalWeight float64, preferredType string) []*models.ContainerPlanDML {
	container := calculateContainerPlan(totalVolume, totalWeight, preferredType)
	container.ContainerPlanId = snowflake.GenID()
	container.ShipmentId = shipmentId
	container.Remark = "系统按体积和重量自动推荐"
	return []*models.ContainerPlanDML{&container}
}

func (service *shipmentService) buildContainerPreview(totalVolume, totalWeight float64, preferredType string, req *models.ShipmentEstimateReq, cargoList []*models.CargoVo) []*models.ContainerPlanVo {
	if preferredType == "LCL" || (preferredType == "" && totalVolume > 0 && totalVolume < 15) {
		return []*models.ContainerPlanVo{}
	}
	candidates := make([]containerCapacity, 0, len(capacities))
	for _, item := range capacities {
		if preferredType == "" || preferredType == item.Type {
			candidates = append(candidates, item)
		}
	}
	if len(candidates) == 0 {
		candidates = append(candidates, capacities...)
	}
	plans := make([]*models.ContainerPlanVo, 0, len(candidates))
	for _, capacity := range candidates {
		plan := calculateContainerPlan(totalVolume, totalWeight, capacity.Type)
		plan.Remark = "按安全装载体积、重量和单件尺寸生成的整柜测算"
		vo := &models.ContainerPlanVo{ContainerPlanDML: plan}
		enrichContainerPlan(vo, capacity, req, cargoList)
		plans = append(plans, vo)
	}
	return plans
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

func (service *shipmentService) UpdateShipmentStatus(shipmentId int64, req *models.ShipmentStatusUpdateReq, username string, operatorUserId int64, canManageAll bool) error {
	if !validStatus(req.Status) {
		return errors.New("出货状态不正确")
	}
	plan := service.shipmentDao.SelectShipmentById(shipmentId)
	if plan == nil {
		return errors.New("出货计划不存在")
	}
	if !canOperateShipment(plan, operatorUserId, canManageAll) {
		return errors.New("无权维护该客户的出货计划")
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

func (service *shipmentService) UpdateShipmentCustomer(shipmentId int64, req *models.ShipmentCustomerBindReq, username string) error {
	if shipmentId == 0 || req == nil || req.CustomerId == 0 {
		return errors.New("请选择要绑定的客户")
	}
	if service.shipmentDao.SelectShipmentById(shipmentId) == nil {
		return errors.New("出货计划不存在")
	}
	customer := customerService.GetCustomerService().SelectCustomerById(req.CustomerId)
	if customer != nil && req.CustomerName == "" {
		req.CustomerName = customer.CustomerName
	}
	service.shipmentDao.UpdateShipmentCustomer(shipmentId, req.CustomerId, req.CustomerName, customerSalesUserId(customer), customerSalesUserName(customer), username)
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

func (service *shipmentService) CanOperateShipment(shipmentId int64, operatorUserId int64, canManageAll bool) bool {
	return canOperateShipment(service.shipmentDao.SelectShipmentById(shipmentId), operatorUserId, canManageAll)
}

func CanManageAllShipments(user *loginModels.User) bool {
	if user == nil {
		return false
	}
	if admin.IsAdmin(user.UserId) {
		return true
	}
	for _, role := range user.Roles {
		if role != nil && role.DataScope == dataScopeAspect.DataScopeAll {
			return true
		}
	}
	return false
}

func canOperateShipment(plan *models.ShipmentPlanVo, operatorUserId int64, canManageAll bool) bool {
	if plan == nil {
		return false
	}
	if canManageAll {
		return true
	}
	return operatorUserId != 0 && plan.SalesUserId == operatorUserId
}

func customerSalesUserId(customer *customerModels.CustomerVo) int64 {
	if customer == nil {
		return 0
	}
	return customer.SalesUserId
}

func customerSalesUserName(customer *customerModels.CustomerVo) string {
	if customer == nil {
		return ""
	}
	return customer.SalesUserName
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

func (service *shipmentService) normalizeCargoList(items []*models.CargoImportReq) ([]*models.CargoVo, *models.ShipmentEstimateSummaryVo, error) {
	cargoList := make([]*models.CargoVo, 0, len(items))
	summary := &models.ShipmentEstimateSummaryVo{}
	for _, item := range items {
		if item == nil || strings.TrimSpace(item.CargoName) == "" {
			continue
		}
		volume := item.VolumeCbm
		if volume == 0 && item.LengthCm > 0 && item.WidthCm > 0 && item.HeightCm > 0 && item.Cartons > 0 {
			volume = item.LengthCm * item.WidthCm * item.HeightCm / 1000000 * float64(item.Cartons)
		}
		cargo := &models.CargoVo{
			CargoDML: models.CargoDML{
				Sku:         strings.TrimSpace(item.Sku),
				CargoName:   strings.TrimSpace(item.CargoName),
				PackageType: strings.TrimSpace(item.PackageType),
				Quantity:    item.Quantity,
				Cartons:     item.Cartons,
				WeightKg:    round2(item.WeightKg),
				VolumeCbm:   round2(volume),
				LengthCm:    round2(item.LengthCm),
				WidthCm:     round2(item.WidthCm),
				HeightCm:    round2(item.HeightCm),
			},
		}
		summary.LineCount++
		summary.TotalQuantity += cargo.Quantity
		summary.TotalCartons += cargo.Cartons
		summary.TotalWeight += cargo.WeightKg
		summary.TotalVolume += cargo.VolumeCbm
		cargoList = append(cargoList, cargo)
	}
	if len(cargoList) == 0 {
		return nil, nil, errors.New("货物明细不能为空")
	}
	summary.TotalWeight = round2(summary.TotalWeight)
	summary.TotalVolume = round2(summary.TotalVolume)
	return cargoList, summary, nil
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

func calculateContainerPlan(totalVolume, totalWeight float64, preferredType string) models.ContainerPlanDML {
	if preferredType == "LCL" || (preferredType == "" && totalVolume > 0 && totalVolume < 15) {
		return models.ContainerPlanDML{
			ContainerType: "LCL",
			Quantity:      1,
			MaxVolume:     15,
			MaxWeight:     3000,
			UsedVolume:    round2(totalVolume),
			UsedWeight:    round2(totalWeight),
			LoadRate:      round2(safeDivide(totalVolume, 15) * 100),
		}
	}
	capacity := capacities[len(capacities)-1]
	for _, item := range capacities {
		if preferredType == item.Type {
			capacity = item
			break
		}
	}
	if preferredType == "" {
		for _, item := range capacities {
			if totalVolume <= item.SafeCbm && totalWeight <= item.WeightKg {
				capacity = item
				break
			}
		}
	}
	quantity := int64(math.Ceil(math.Max(safeDivide(totalVolume, capacity.SafeCbm), safeDivide(totalWeight, capacity.WeightKg))))
	if quantity < 1 {
		quantity = 1
	}
	maxVolume := capacity.SafeCbm * float64(quantity)
	maxWeight := capacity.WeightKg * float64(quantity)
	loadRate := math.Max(safeDivide(totalVolume, maxVolume), safeDivide(totalWeight, maxWeight)) * 100
	return models.ContainerPlanDML{
		ContainerType: capacity.Type,
		Quantity:      quantity,
		MaxVolume:     round2(maxVolume),
		MaxWeight:     round2(maxWeight),
		UsedVolume:    round2(totalVolume),
		UsedWeight:    round2(totalWeight),
		LoadRate:      round2(loadRate),
	}
}

func buildLclSuggestion(totalVolume float64, preferredType string, req *models.ShipmentEstimateReq) *models.ShipmentEstimateLclVo {
	recommended := preferredType == "LCL" || (preferredType == "" && totalVolume > 0 && totalVolume < 15)
	remark := "当前体积建议优先整柜测算"
	if recommended {
		remark = "当前体积适合优先按散货拼箱评估"
	}
	rate := defaultIfZero(req.LclRate, 85)
	minCharge := defaultIfZero(req.LclMinCharge, 300)
	extraFees := req.ExtraFees
	totalCost := math.Max(totalVolume*rate, minCharge) + extraFees
	return &models.ShipmentEstimateLclVo{
		Recommended: recommended,
		TotalVolume: round2(totalVolume),
		RatePerCbm:  round2(rate),
		MinCharge:   round2(minCharge),
		ExtraFees:   round2(extraFees),
		TotalCost:   round2(totalCost),
		Remark:      remark,
	}
}

func enrichContainerPlan(plan *models.ContainerPlanVo, capacity containerCapacity, req *models.ShipmentEstimateReq, cargoList []*models.CargoVo) {
	plan.InternalLengthCm = capacity.LengthCm
	plan.InternalWidthCm = capacity.WidthCm
	plan.InternalHeightCm = capacity.HeightCm
	plan.SafeVolume = capacity.SafeCbm * float64(plan.Quantity)
	plan.EffectiveVolume = round2(effectiveVolume(cargoList))
	plan.UnitCost = containerRate(capacity.Type, req)
	plan.ExtraFees = round2(req.ExtraFees)
	plan.TotalCost = round2(plan.UnitCost*float64(plan.Quantity) + plan.ExtraFees)
	plan.Warnings = validateCargoForContainer(cargoList, capacity)
	if plan.EffectiveVolume > plan.SafeVolume {
		plan.Warnings = append(plan.Warnings, "货物折算装载体积超过该柜型安全容量，实际装柜风险较高")
	}
	plan.RiskLevel = riskLevel(plan.LoadRate, len(plan.Warnings))
	plan.Remark = plan.Remark + "；风险等级：" + plan.RiskLevel
	plan.Placements = simulatePlacements(cargoList, capacity, 120)
}

func effectiveVolume(cargoList []*models.CargoVo) float64 {
	total := 0.0
	for _, item := range cargoList {
		total += item.VolumeCbm / loadingEfficiency(item.PackageType)
	}
	return total
}

func loadingEfficiency(packageType string) float64 {
	value := strings.ToLower(packageType)
	switch {
	case strings.Contains(value, "托盘"), strings.Contains(value, "pallet"):
		return 0.78
	case strings.Contains(value, "不可堆"), strings.Contains(value, "易碎"), strings.Contains(value, "fragile"):
		return 0.62
	case strings.Contains(value, "机械"), strings.Contains(value, "设备"), strings.Contains(value, "machine"):
		return 0.65
	default:
		return 0.88
	}
}

func validateCargoForContainer(cargoList []*models.CargoVo, capacity containerCapacity) []string {
	warnings := make([]string, 0)
	for _, item := range cargoList {
		if item.LengthCm == 0 || item.WidthCm == 0 || item.HeightCm == 0 {
			continue
		}
		if !canFitItem(item.LengthCm, item.WidthCm, item.HeightCm, capacity) {
			warnings = append(warnings, fmt.Sprintf("%s 单件尺寸 %.0f×%.0f×%.0fcm 可能无法放入 %s", item.CargoName, item.LengthCm, item.WidthCm, item.HeightCm, capacity.Type))
		}
	}
	return warnings
}

func canFitItem(length, width, height float64, capacity containerCapacity) bool {
	options := [][3]float64{
		{length, width, height}, {length, height, width}, {width, length, height},
		{width, height, length}, {height, length, width}, {height, width, length},
	}
	for _, item := range options {
		if item[0] <= capacity.LengthCm && item[1] <= capacity.WidthCm && item[2] <= capacity.HeightCm {
			return true
		}
	}
	return false
}

func riskLevel(loadRate float64, warningCount int) string {
	switch {
	case warningCount > 0 || loadRate >= 94:
		return "高"
	case loadRate >= 82:
		return "中"
	default:
		return "低"
	}
}

func containerRate(containerType string, req *models.ShipmentEstimateReq) float64 {
	switch containerType {
	case "20GP":
		return defaultIfZero(req.Rate20GP, 1800)
	case "40GP":
		return defaultIfZero(req.Rate40GP, 2600)
	case "40HQ":
		return defaultIfZero(req.Rate40HQ, 2800)
	default:
		return 0
	}
}

func defaultIfZero(value, fallback float64) float64 {
	if value > 0 {
		return value
	}
	return fallback
}

func simulatePlacements(cargoList []*models.CargoVo, capacity containerCapacity, maxPlacements int) []*models.LoadingPlacementVo {
	placements := make([]*models.LoadingPlacementVo, 0)
	x, y, rowDepth := 0.0, 0.0, 0.0
	for index, cargo := range cargoList {
		length, width, height := chooseFloorOrientation(cargo.LengthCm, cargo.WidthCm, cargo.HeightCm, capacity)
		if length <= 0 || width <= 0 || height <= 0 || cargo.Cartons <= 0 {
			continue
		}
		stackable := !strings.Contains(cargo.PackageType, "不可堆") && !strings.Contains(strings.ToLower(cargo.PackageType), "fragile")
		layers := int64(1)
		if stackable {
			layers = int64(math.Max(1, math.Floor(capacity.HeightCm/height)))
		}
		stackHeight := height * float64(layers)
		remaining := cargo.Cartons
		for remaining > 0 && len(placements) < maxPlacements {
			if y+width > capacity.WidthCm {
				x += rowDepth
				y = 0
				rowDepth = 0
			}
			if x+length > capacity.LengthCm {
				return placements
			}
			qty := layers
			if remaining < qty {
				qty = remaining
			}
			placements = append(placements, &models.LoadingPlacementVo{
				CargoName: cargo.CargoName,
				Sku:       cargo.Sku,
				Color:     placementColors[index%len(placementColors)],
				Quantity:  qty,
				X:         round2(x),
				Y:         round2(y),
				Z:         0,
				Length:    round2(length),
				Width:     round2(width),
				Height:    round2(stackHeight),
				Remark:    fmt.Sprintf("%d 件/叠", qty),
			})
			remaining -= qty
			y += width
			if length > rowDepth {
				rowDepth = length
			}
		}
	}
	return placements
}

func chooseFloorOrientation(length, width, height float64, capacity containerCapacity) (float64, float64, float64) {
	options := [][3]float64{{length, width, height}, {width, length, height}}
	for _, item := range options {
		if item[0] <= capacity.LengthCm && item[1] <= capacity.WidthCm && item[2] <= capacity.HeightCm {
			return item[0], item[1], item[2]
		}
	}
	return length, width, height
}

func buildLoadingPlan(containers []*models.ContainerPlanVo) *models.LoadingPlanVo {
	if len(containers) == 0 {
		return nil
	}
	best := containers[0]
	for _, item := range containers {
		if item.RiskLevel != "高" && (best.RiskLevel == "高" || item.TotalCost < best.TotalCost) {
			best = item
		}
	}
	return &models.LoadingPlanVo{
		ContainerType:    best.ContainerType,
		Quantity:         best.Quantity,
		InternalLengthCm: best.InternalLengthCm,
		InternalWidthCm:  best.InternalWidthCm,
		InternalHeightCm: best.InternalHeightCm,
		ViewScale:        1,
		Utilization:      best.LoadRate,
		Placements:       best.Placements,
	}
}

func buildRecommendation(containers []*models.ContainerPlanVo, lcl *models.ShipmentEstimateLclVo, preferredType string) *models.ShipmentRecommendationVo {
	if preferredType == "LCL" || len(containers) == 0 {
		return &models.ShipmentRecommendationVo{Mode: "LCL", Title: "推荐散货拼箱", Reason: lcl.Remark, RiskLevel: "低", Confidence: "中"}
	}
	best := containers[0]
	for _, item := range containers {
		if item.RiskLevel == "高" && best.RiskLevel != "高" {
			continue
		}
		if best.RiskLevel == "高" || item.TotalCost < best.TotalCost {
			best = item
		}
	}
	mode := "FCL"
	title := fmt.Sprintf("推荐 %d×%s", best.Quantity, best.ContainerType)
	reason := "整柜方案在当前体积下更适合进一步报价"
	saving := 0.0
	if lcl.TotalCost > 0 && best.TotalCost > 0 {
		saving = round2(math.Abs(lcl.TotalCost - best.TotalCost))
		if lcl.TotalCost < best.TotalCost && preferredType == "" {
			mode = "LCL"
			title = "推荐散货拼箱"
			reason = "散货估算成本低于当前整柜方案"
		} else {
			reason = "整柜成本与装载空间更适合当前货量"
		}
	}
	confidence := "高"
	if best.RiskLevel == "中" {
		confidence = "中"
	}
	if best.RiskLevel == "高" {
		confidence = "低"
	}
	return &models.ShipmentRecommendationVo{Mode: mode, Title: title, Reason: reason, Saving: saving, RiskLevel: best.RiskLevel, Confidence: confidence}
}

func collectEstimateWarnings(containers []*models.ContainerPlanVo) []string {
	seen := map[string]bool{}
	warnings := make([]string, 0)
	for _, container := range containers {
		for _, warning := range container.Warnings {
			if !seen[warning] {
				seen[warning] = true
				warnings = append(warnings, warning)
			}
		}
	}
	return warnings
}

func safeDivide(dividend, divisor float64) float64 {
	if divisor == 0 {
		return 0
	}
	return dividend / divisor
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
