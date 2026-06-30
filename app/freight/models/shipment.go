package models

import (
	"baize/app/common/baize/baizeUnix"
	"baize/app/common/commonModels"
)

type ShipmentPlanDQL struct {
	ShipmentNo   string `form:"shipmentNo" db:"shipment_no"`
	OrderNo      string `form:"orderNo" db:"order_no"`
	CustomerId   int64  `form:"customerId" db:"customer_id"`
	CustomerName string `form:"customerName" db:"customer_name"`
	Pol          string `form:"pol" db:"pol"`
	Pod          string `form:"pod" db:"pod"`
	Status       string `form:"status" db:"status"`
	BeginTime    string `form:"beginTime" db:"begin_time"`
	EndTime      string `form:"endTime" db:"end_time"`
	commonModels.BaseEntityDQL
}

type ShipmentImportReq struct {
	CustomerId    int64             `json:"customerId,string" binding:"required"`
	CustomerName  string            `json:"customerName"`
	OrderNo       string            `json:"orderNo"`
	Pol           string            `json:"pol"`
	Pod           string            `json:"pod"`
	PlannedEtd    string            `json:"plannedEtd"`
	PlannedEta    string            `json:"plannedEta"`
	Remark        string            `json:"remark"`
	CargoList     []*CargoImportReq `json:"cargoList" binding:"required"`
	PreferredType string            `json:"preferredType"`
}

type CargoImportReq struct {
	Sku         string  `json:"sku"`
	CargoName   string  `json:"cargoName" binding:"required"`
	PackageType string  `json:"packageType"`
	Quantity    int64   `json:"quantity"`
	Cartons     int64   `json:"cartons"`
	WeightKg    float64 `json:"weightKg"`
	VolumeCbm   float64 `json:"volumeCbm"`
	LengthCm    float64 `json:"lengthCm"`
	WidthCm     float64 `json:"widthCm"`
	HeightCm    float64 `json:"heightCm"`
}

type ShipmentPlanDML struct {
	ShipmentId   int64   `json:"shipmentId,string" db:"shipment_id"`
	ShipmentNo   string  `json:"shipmentNo" db:"shipment_no"`
	OrderNo      string  `json:"orderNo" db:"order_no"`
	CustomerId   int64   `json:"customerId,string" db:"customer_id"`
	CustomerName string  `json:"customerName" db:"customer_name"`
	Pol          string  `json:"pol" db:"pol"`
	Pod          string  `json:"pod" db:"pod"`
	PlannedEtd   string  `json:"plannedEtd" db:"planned_etd"`
	PlannedEta   string  `json:"plannedEta" db:"planned_eta"`
	Status       string  `json:"status" db:"status"`
	TotalWeight  float64 `json:"totalWeight" db:"total_weight"`
	TotalVolume  float64 `json:"totalVolume" db:"total_volume"`
	TotalCartons int64   `json:"totalCartons" db:"total_cartons"`
	ShareToken   string  `json:"shareToken" db:"share_token"`
	Remark       string  `json:"remark" db:"remark"`
	CreateBy     string  `json:"createBy" db:"create_by"`
	UpdateBy     string  `json:"updateBy" db:"update_by"`
}

type ShipmentPlanVo struct {
	ShipmentId   int64                `json:"shipmentId,string" db:"shipment_id"`
	ShipmentNo   string               `json:"shipmentNo" db:"shipment_no"`
	OrderNo      string               `json:"orderNo" db:"order_no"`
	CustomerId   int64                `json:"customerId,string" db:"customer_id"`
	CustomerName string               `json:"customerName" db:"customer_name"`
	Pol          string               `json:"pol" db:"pol"`
	Pod          string               `json:"pod" db:"pod"`
	PlannedEtd   string               `json:"plannedEtd" db:"planned_etd"`
	PlannedEta   string               `json:"plannedEta" db:"planned_eta"`
	ActualEtd    string               `json:"actualEtd" db:"actual_etd"`
	ActualEta    string               `json:"actualEta" db:"actual_eta"`
	Status       string               `json:"status" db:"status"`
	TotalWeight  float64              `json:"totalWeight" db:"total_weight"`
	TotalVolume  float64              `json:"totalVolume" db:"total_volume"`
	TotalCartons int64                `json:"totalCartons" db:"total_cartons"`
	ShareToken   string               `json:"shareToken" db:"share_token"`
	Remark       string               `json:"remark" db:"remark"`
	CreateBy     string               `json:"createBy" db:"create_by"`
	CreateTime   *baizeUnix.BaiZeTime `json:"createTime" db:"create_time"`
	UpdateBy     string               `json:"updateBy" db:"update_by"`
	UpdateTime   *baizeUnix.BaiZeTime `json:"updateTime" db:"update_time"`
}

type CargoDML struct {
	CargoId     int64   `json:"cargoId,string" db:"cargo_id"`
	ShipmentId  int64   `json:"shipmentId,string" db:"shipment_id"`
	Sku         string  `json:"sku" db:"sku"`
	CargoName   string  `json:"cargoName" db:"cargo_name"`
	PackageType string  `json:"packageType" db:"package_type"`
	Quantity    int64   `json:"quantity" db:"quantity"`
	Cartons     int64   `json:"cartons" db:"cartons"`
	WeightKg    float64 `json:"weightKg" db:"weight_kg"`
	VolumeCbm   float64 `json:"volumeCbm" db:"volume_cbm"`
	LengthCm    float64 `json:"lengthCm" db:"length_cm"`
	WidthCm     float64 `json:"widthCm" db:"width_cm"`
	HeightCm    float64 `json:"heightCm" db:"height_cm"`
}

type CargoVo struct {
	CargoDML
}

type ContainerPlanDML struct {
	ContainerPlanId int64   `json:"containerPlanId,string" db:"container_plan_id"`
	ShipmentId      int64   `json:"shipmentId,string" db:"shipment_id"`
	ContainerType   string  `json:"containerType" db:"container_type"`
	Quantity        int64   `json:"quantity" db:"quantity"`
	MaxVolume       float64 `json:"maxVolume" db:"max_volume"`
	MaxWeight       float64 `json:"maxWeight" db:"max_weight"`
	UsedVolume      float64 `json:"usedVolume" db:"used_volume"`
	UsedWeight      float64 `json:"usedWeight" db:"used_weight"`
	LoadRate        float64 `json:"loadRate" db:"load_rate"`
	Remark          string  `json:"remark" db:"remark"`
}

type ContainerPlanVo struct {
	ContainerPlanDML
}

type ShipmentOrderDML struct {
	OrderId    int64  `json:"orderId,string" db:"order_id"`
	ShipmentId int64  `json:"shipmentId,string" db:"shipment_id"`
	OrderNo    string `json:"orderNo" db:"order_no"`
	Status     string `json:"status" db:"status"`
	CreateBy   string `json:"createBy" db:"create_by"`
	UpdateBy   string `json:"updateBy" db:"update_by"`
}

type ShipmentOrderVo struct {
	OrderId    int64                `json:"orderId,string" db:"order_id"`
	ShipmentId int64                `json:"shipmentId,string" db:"shipment_id"`
	OrderNo    string               `json:"orderNo" db:"order_no"`
	Status     string               `json:"status" db:"status"`
	CreateTime *baizeUnix.BaiZeTime `json:"createTime" db:"create_time"`
}

type ShipmentDetailVo struct {
	Plan       *ShipmentPlanVo       `json:"plan"`
	CargoList  []*CargoVo            `json:"cargoList"`
	Containers []*ContainerPlanVo    `json:"containers"`
	Order      *ShipmentOrderVo      `json:"order"`
	StatusFlow []*ShipmentStatusStep `json:"statusFlow"`
}

type ShipmentStatusStep struct {
	Value  string `json:"value"`
	Label  string `json:"label"`
	Active bool   `json:"active"`
}

type ShipmentStatusUpdateReq struct {
	Status    string `json:"status" binding:"required"`
	ActualEtd string `json:"actualEtd"`
	ActualEta string `json:"actualEta"`
	Remark    string `json:"remark"`
}

type ShipmentStatusUpdateDML struct {
	ShipmentId int64  `db:"shipment_id"`
	Status     string `db:"status"`
	ActualEtd  string `db:"actual_etd"`
	ActualEta  string `db:"actual_eta"`
	Remark     string `db:"remark"`
	UpdateBy   string `db:"update_by"`
}

type ShareInfoVo struct {
	ShareUrl string `json:"shareUrl"`
	Token    string `json:"token"`
}
