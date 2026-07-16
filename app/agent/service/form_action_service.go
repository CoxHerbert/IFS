package service

import (
	"baize/app/agent/dao"
	"baize/app/agent/protocol"
	"baize/app/agent/request"
	customerService "baize/app/customer/service"
	freightModels "baize/app/freight/models"
	freightService "baize/app/freight/service"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type FormActionService struct {
	Dao *dao.ChatDao
}

var formActionService = &FormActionService{Dao: dao.GetChatDao()}

func GetFormActionService() *FormActionService {
	return formActionService
}

func (s *FormActionService) SubmitForm(req *request.FormSubmitRequest) protocol.AgentResult {
	payload, _ := json.Marshal(req.Values)
	s.Dao.InsertFormSubmission(req.SessionID, req.FormCode, string(payload))

	switch req.FormCode {
	case "save_shipment_plan":
		return s.saveShipmentPlan(req)
	case "admin_receipt_create":
		return s.createReceipt(req)
	case "admin_receipt_allocate":
		return s.allocateReceipt(req)
	case "customer_payment_declaration":
		return s.createPaymentDeclaration(req)
	}

	result := protocol.NewAgentResultV2("出货信息已确认", "已收到出货信息，可以继续生成出货计划。", []protocol.BlockItem{
		{Type: "markdown", Content: "请继续上传货物明细或提交保存出货计划。"},
	})
	s.Dao.InsertMessage(req.SessionID, "assistant", result.Summary, mustJSON(result), configuredDefaultModel())
	return result
}

func (s *FormActionService) createPaymentDeclaration(req *request.FormSubmitRequest) protocol.AgentResult {
	if req.Source != "customer" || req.CustomerID == 0 {
		return protocol.NewErrorResult("该操作仅允许登录客户提交")
	}

	shipmentNo := strings.TrimSpace(valueOrEmpty(req.Values["shipmentNo"]))
	query := &freightModels.ShipmentPlanDQL{ShipmentNo: shipmentNo, CustomerId: req.CustomerID}
	query.Size = 2
	query.Limit = " limit 0,2"
	plans, _ := freightService.GetShipmentService().SelectShipmentList(query)
	if len(plans) != 1 || !strings.EqualFold(plans[0].ShipmentNo, shipmentNo) {
		return protocol.NewErrorResult("关联出货计划不存在")
	}

	item := &freightModels.PaymentDeclarationDML{
		CustomerId:   req.CustomerID,
		CustomerName: req.CustomerName,
		ShipmentId:   plans[0].ShipmentId,
		Amount:       toFloat64(req.Values["amount"]),
		Currency:     valueOrEmpty(req.Values["currency"]),
		PaymentTime:  valueOrEmpty(req.Values["paymentTime"]),
		VoucherUrl:   req.VoucherURL,
		VoucherName:  req.VoucherName,
		Remark:       valueOrEmpty(req.Values["remark"]),
		CreateBy:     req.OperatorName,
	}
	if err := freightService.GetReceiptService().CreatePaymentDeclaration(item); err != nil {
		return protocol.NewErrorResult(err.Error())
	}

	result := protocol.NewAgentResultV2("付款申报已提交", "后台审核确认到账后才会生成正式收款与核销。", []protocol.BlockItem{
		{Type: "metrics", Items: []protocol.MetricItem{
			{Label: "申报单号", Value: item.DeclarationNo},
			{Label: "出货计划", Value: item.ShipmentNo},
			{Label: "申报金额", Value: fmt.Sprintf("%s %.2f", item.Currency, item.Amount)},
			{Label: "状态", Value: "待审核"},
		}},
	})
	s.Dao.InsertMessage(req.SessionID, "assistant", result.Summary, mustJSON(result), configuredDefaultModel())
	return result
}

func (s *FormActionService) allocateReceipt(req *request.FormSubmitRequest) protocol.AgentResult {
	if req.Source != "admin" || !hasPermissionValue(req.Permissions, "freight:receipt:add") {
		return protocol.NewErrorResult("该操作仅允许有收款新增权限的管理端用户执行")
	}

	receiptNo := strings.TrimSpace(valueOrEmpty(req.Values["receiptNo"]))
	receiptQuery := &freightModels.ReceiptDQL{ReceiptNo: receiptNo}
	receiptQuery.Size = 2
	receiptQuery.Limit = " limit 0,2"
	receipts, _ := freightService.GetReceiptService().SelectList(receiptQuery)
	if len(receipts) != 1 || !strings.EqualFold(receipts[0].ReceiptNo, receiptNo) {
		return protocol.NewErrorResult("收款单不存在")
	}

	shipmentNo := strings.TrimSpace(valueOrEmpty(req.Values["shipmentNo"]))
	shipmentQuery := &freightModels.ShipmentPlanDQL{ShipmentNo: shipmentNo, CustomerId: receipts[0].CustomerId}
	shipmentQuery.Size = 2
	shipmentQuery.Limit = " limit 0,2"
	plans, _ := freightService.GetShipmentService().SelectShipmentList(shipmentQuery)
	if len(plans) != 1 || !strings.EqualFold(plans[0].ShipmentNo, shipmentNo) {
		return protocol.NewErrorResult("出货计划不存在或不属于收款客户")
	}

	vo, err := freightService.GetReceiptService().Allocate(receipts[0].ReceiptId, plans[0].ShipmentId, toFloat64(req.Values["allocatedAmount"]))
	if err != nil {
		return protocol.NewErrorResult(err.Error())
	}

	result := protocol.NewAgentResultV2("核销完成", "收款单已追加核销。", []protocol.BlockItem{
		{Type: "metrics", Items: []protocol.MetricItem{
			{Label: "收款单号", Value: vo.ReceiptNo},
			{Label: "本次出货计划", Value: shipmentNo},
			{Label: "累计已核销", Value: vo.AllocatedAmount},
			{Label: "核销状态", Value: vo.Status},
		}},
	})
	s.Dao.InsertMessage(req.SessionID, "assistant", result.Summary, mustJSON(result), configuredDefaultModel())
	return result
}

func (s *FormActionService) createReceipt(req *request.FormSubmitRequest) protocol.AgentResult {
	if req.Source != "admin" || !hasPermissionValue(req.Permissions, "freight:receipt:add") {
		return protocol.NewErrorResult("该操作仅允许在管理端执行")
	}

	customerName := strings.TrimSpace(valueOrEmpty(req.Values["customerName"]))
	options := customerService.GetCustomerService().SelectCustomerOptions(customerName)
	if len(options) != 1 {
		return protocol.NewErrorResult("客户名称必须唯一匹配，请输入更准确的客户名称")
	}

	amount := toFloat64(req.Values["amount"])
	shipmentNos := splitCSV(valueOrEmpty(req.Values["shipmentNos"]))
	allocationAmounts := splitCSV(valueOrEmpty(req.Values["allocationAmounts"]))
	allocations := make([]*freightModels.ReceiptAllocationReq, 0, len(shipmentNos))
	for i, no := range shipmentNos {
		query := &freightModels.ShipmentPlanDQL{ShipmentNo: no, CustomerId: options[0].CustomerId}
		query.Size = 2
		query.Limit = " limit 0,2"
		plans, _ := freightService.GetShipmentService().SelectShipmentList(query)
		if len(plans) != 1 || !strings.EqualFold(plans[0].ShipmentNo, no) {
			return protocol.NewErrorResult("未找到客户名下的出货计划：" + no)
		}
		allocated := amount
		if i < len(allocationAmounts) {
			allocated = toFloat64(allocationAmounts[i])
		}
		allocations = append(allocations, &freightModels.ReceiptAllocationReq{ShipmentId: plans[0].ShipmentId, AllocatedAmount: allocated})
	}

	vo, err := freightService.GetReceiptService().Create(&freightModels.ReceiptCreateReq{
		CustomerId:     options[0].CustomerId,
		CustomerName:   options[0].CustomerName,
		Amount:         amount,
		Currency:       valueOrEmpty(req.Values["currency"]),
		PaymentMethod:  valueOrEmpty(req.Values["paymentMethod"]),
		Remark:         valueOrEmpty(req.Values["remark"]),
		Allocations:    allocations,
	}, "", "", req.OperatorName)
	if err != nil {
		return protocol.NewErrorResult(err.Error())
	}

	result := protocol.NewAgentResultV2("收款单已创建", "收款已登记并完成指定核销。", []protocol.BlockItem{
		{Type: "metrics", Items: []protocol.MetricItem{
			{Label: "收款单号", Value: vo.ReceiptNo},
			{Label: "客户", Value: vo.CustomerName},
			{Label: "收款金额", Value: fmt.Sprintf("%s %.2f", vo.Currency, vo.Amount)},
			{Label: "已核销", Value: vo.AllocatedAmount},
			{Label: "核销状态", Value: vo.Status},
		}},
	})
	s.Dao.InsertMessage(req.SessionID, "assistant", result.Summary, mustJSON(result), configuredDefaultModel())
	return result
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if text := strings.TrimSpace(part); text != "" {
			result = append(result, text)
		}
	}
	return result
}

func hasPermissionValue(values []string, target string) bool {
	for _, value := range values {
		if value == "*:*:*" || value == target {
			return true
		}
	}
	return false
}

func (s *FormActionService) saveShipmentPlan(req *request.FormSubmitRequest) protocol.AgentResult {
	importReq, err := buildShipmentImportReq(req)
	if err != nil {
		result := protocol.NewAgentResultV2("保存出货计划失败", err.Error(), []protocol.BlockItem{
			{Type: "error", Title: "保存失败", Content: err.Error()},
		})
		s.Dao.InsertMessage(req.SessionID, "assistant", result.Summary, mustJSON(result), configuredDefaultModel())
		return result
	}

	operator := req.OperatorName
	if operator == "" {
		operator = "agent"
	}
	detail, err := freightService.GetShipmentService().ImportShipment(importReq, operator, req.OperatorID)
	if err != nil {
		result := protocol.NewAgentResultV2("保存出货计划失败", err.Error(), []protocol.BlockItem{
			{Type: "error", Title: "保存失败", Content: err.Error()},
		})
		s.Dao.InsertMessage(req.SessionID, "assistant", result.Summary, mustJSON(result), configuredDefaultModel())
		return result
	}

	customerValue := detail.Plan.CustomerName
	if customerValue == "" {
		customerValue = "未绑定客户"
	}
	summary := "出货计划已保存，可在后台「货代业务 / 出货计划」中查看。"
	if detail.Plan.CustomerId == 0 {
		summary = "出货计划已保存，当前未绑定客户，可在后台出货计划列表中后续绑定。"
	}

	result := protocol.NewAgentResultV2("出货计划已保存", summary, []protocol.BlockItem{
		{
			Type: "metrics",
			Items: []protocol.MetricItem{
				{Label: "计划编号", Value: detail.Plan.ShipmentNo},
				{Label: "客户", Value: customerValue},
				{Label: "总箱数", Value: detail.Plan.TotalCartons},
				{Label: "总体积", Value: formatNumber(detail.Plan.TotalVolume) + " CBM"},
				{Label: "状态", Value: detail.Plan.Status},
			},
		},
		{
			Type:    "link",
			Title:   "查看出货计划",
			Content: "出货计划已生成，可打开列表查看并继续维护。",
			Label:   "打开出货计划列表",
			URL:     "/freight/shipment",
		},
	})
	s.Dao.InsertMessage(req.SessionID, "assistant", result.Summary, mustJSON(result), configuredDefaultModel())
	return result
}

func buildShipmentImportReq(req *request.FormSubmitRequest) (*freightModels.ShipmentImportReq, error) {
	cargoList := buildCargoImportList(req.Values["cargoList"])
	if len(cargoList) == 0 {
		return nil, fmt.Errorf("缺少货物明细，无法保存为出货计划")
	}

	customerID := req.CustomerID
	customerName := req.CustomerName
	// 兼容历史表单值；新表单不再要求用户在对话里选择客户。
	if customerID == 0 {
		customerID = toInt64(req.Values["customerId"])
	}
	if customerName == "" {
		customerName = valueOrEmpty(req.Values["customerName"])
	}

	return &freightModels.ShipmentImportReq{
		CustomerId:    customerID,
		CustomerName:  customerName,
		OrderNo:       valueOrEmpty(req.Values["orderNo"]),
		Pol:           valueOrEmpty(req.Values["pol"]),
		Pod:           valueOrEmpty(req.Values["pod"]),
		PlannedEtd:    valueOrEmpty(req.Values["plannedEtd"]),
		PlannedEta:    valueOrEmpty(req.Values["plannedEta"]),
		Remark:        valueOrEmpty(req.Values["remark"]),
		PreferredType: valueOrEmpty(req.Values["preferredType"]),
		CargoList:     cargoList,
	}, nil
}

func buildCargoImportList(value any) []*freightModels.CargoImportReq {
	rows, ok := value.([]any)
	if !ok {
		return nil
	}
	list := make([]*freightModels.CargoImportReq, 0, len(rows))
	for _, row := range rows {
		item, ok := row.(map[string]any)
		if !ok {
			continue
		}
		cargoName := valueOrEmpty(item["cargoName"])
		if cargoName == "" {
			continue
		}
		list = append(list, &freightModels.CargoImportReq{
			Sku:       valueOrEmpty(item["sku"]),
			CargoName: cargoName,
			Quantity:  toInt64(item["quantity"]),
			Cartons:   toInt64(item["cartons"]),
			WeightKg:  toFloat64(item["weightKg"]),
			VolumeCbm: toFloat64(item["volumeCbm"]),
			LengthCm:  toFloat64(item["lengthCm"]),
			WidthCm:   toFloat64(item["widthCm"]),
			HeightCm:  toFloat64(item["heightCm"]),
		})
	}
	return list
}

func (s *FormActionService) ExecuteAction(req *request.ActionExecuteRequest) protocol.AgentResult {
	if req.ActionCode != "generate_shipment_plan" {
		result := protocol.NewAgentResultV2("操作暂不支持", "当前 action 暂不支持。", []protocol.BlockItem{
			{Type: "error", Title: "操作失败", Content: "unsupported action: " + req.ActionCode},
		})
		s.Dao.InsertMessage(req.SessionID, "assistant", result.Summary, mustJSON(result), configuredDefaultModel())
		return result
	}

	result := protocol.NewAgentResultV2("请先上传货物明细", "上传 Excel 后，分析结果中会出现保存表单。", []protocol.BlockItem{
		{Type: "markdown", Content: "客户端发起时系统会自动使用当前登录账号绑定的客户；后台发起时先创建未绑定计划，后续在出货计划列表绑定客户。"},
	})
	s.Dao.InsertMessage(req.SessionID, "assistant", result.Summary, mustJSON(result), configuredDefaultModel())
	return result
}

func valueOrEmpty(value any) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

func toInt64(value any) int64 {
	switch v := value.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	case json.Number:
		n, _ := v.Int64()
		return n
	case string:
		n, _ := strconv.ParseInt(v, 10, 64)
		return n
	default:
		return 0
	}
}

func toFloat64(value any) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case json.Number:
		n, _ := v.Float64()
		return n
	case string:
		n, _ := strconv.ParseFloat(v, 64)
		return n
	default:
		return 0
	}
}

func mustJSON(value any) string {
	payload, _ := json.Marshal(value)
	return string(payload)
}
