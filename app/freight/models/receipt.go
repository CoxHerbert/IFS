package models

import (
	"baize/app/common/baize/baizeUnix"
	"baize/app/common/commonModels"
)

type ReceiptDQL struct {
	ReceiptNo string `form:"receiptNo" db:"receipt_no"`
	CustomerName string `form:"customerName" db:"customer_name"`
	Status string `form:"status" db:"status"`
	CustomerId int64 `form:"customerId" db:"customer_id"`
	SalesUserId int64 `form:"salesUserId" db:"sales_user_id"`
	commonModels.BaseEntityDQL
}

type ReceiptAllocationReq struct {
	ShipmentId int64 `json:"shipmentId,string"`
	AllocatedAmount float64 `json:"allocatedAmount"`
}

type ReceiptCreateReq struct {
	CustomerId int64 `json:"customerId,string"`
	CustomerName string `json:"customerName"`
	Amount float64 `json:"amount"`
	Currency string `json:"currency"`
	ReceiptTime string `json:"receiptTime"`
	PaymentMethod string `json:"paymentMethod"`
	Remark string `json:"remark"`
	Allocations []*ReceiptAllocationReq `json:"allocations"`
}

type ReceiptDML struct {
	ReceiptId int64 `json:"receiptId,string" db:"receipt_id"`
	ReceiptNo string `json:"receiptNo" db:"receipt_no"`
	CustomerId int64 `json:"customerId,string" db:"customer_id"`
	CustomerName string `json:"customerName" db:"customer_name"`
	Amount float64 `json:"amount" db:"amount"`
	Currency string `json:"currency" db:"currency"`
	ReceiptTime string `json:"receiptTime" db:"receipt_time"`
	PaymentMethod string `json:"paymentMethod" db:"payment_method"`
	Status string `json:"status" db:"status"`
	VoucherUrl string `json:"voucherUrl" db:"voucher_url"`
	VoucherName string `json:"voucherName" db:"voucher_name"`
	Remark string `json:"remark" db:"remark"`
	CreateBy string `json:"createBy" db:"create_by"`
}

type ReceiptVo struct {
	ReceiptDML
	AllocatedAmount float64 `json:"allocatedAmount" db:"allocated_amount"`
	CreateTime *baizeUnix.BaiZeTime `json:"createTime" db:"create_time"`
	Allocations []*ReceiptAllocationVo `json:"allocations,omitempty"`
}

type ReceiptAllocationDML struct {
	AllocationId int64 `db:"allocation_id"`
	ReceiptId int64 `db:"receipt_id"`
	ShipmentId int64 `db:"shipment_id"`
	AllocatedAmount float64 `db:"allocated_amount"`
}

type ReceiptAllocationVo struct {
	AllocationId int64 `json:"allocationId,string" db:"allocation_id"`
	ReceiptId int64 `json:"receiptId,string" db:"receipt_id"`
	ShipmentId int64 `json:"shipmentId,string" db:"shipment_id"`
	ShipmentNo string `json:"shipmentNo" db:"shipment_no"`
	AllocatedAmount float64 `json:"allocatedAmount" db:"allocated_amount"`
}

type PaymentDeclarationDML struct {
	DeclarationId int64 `json:"declarationId,string" db:"declaration_id"`
	DeclarationNo string `json:"declarationNo" db:"declaration_no"`
	CustomerId int64 `json:"customerId,string" db:"customer_id"`
	CustomerName string `json:"customerName" db:"customer_name"`
	ShipmentId int64 `json:"shipmentId,string" db:"shipment_id"`
	ShipmentNo string `json:"shipmentNo" db:"shipment_no"`
	Amount float64 `json:"amount" db:"amount"`
	Currency string `json:"currency" db:"currency"`
	PaymentTime string `json:"paymentTime" db:"payment_time"`
	VoucherUrl string `json:"voucherUrl" db:"voucher_url"`
	VoucherName string `json:"voucherName" db:"voucher_name"`
	Status string `json:"status" db:"status"`
	Remark string `json:"remark" db:"remark"`
	CreateBy string `json:"createBy" db:"create_by"`
}

type PaymentDeclarationDQL struct {
	DeclarationNo string `form:"declarationNo" db:"declaration_no"`
	CustomerName string `form:"customerName" db:"customer_name"`
	Status string `form:"status" db:"status"`
	SalesUserId int64 `form:"salesUserId" db:"sales_user_id"`
	commonModels.BaseEntityDQL
}

type PaymentDeclarationVo struct {
	PaymentDeclarationDML
	ReviewBy string `json:"reviewBy" db:"review_by"`
	ReviewTime *baizeUnix.BaiZeTime `json:"reviewTime" db:"review_time"`
	ReviewRemark string `json:"reviewRemark" db:"review_remark"`
	CreateTime *baizeUnix.BaiZeTime `json:"createTime" db:"create_time"`
}

type PaymentDeclarationReviewReq struct { Remark string `json:"remark"` }
