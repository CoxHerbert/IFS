package service

import (
	"baize/app/agent/dao"
	"baize/app/agent/protocol"
	"baize/app/agent/request"
	freightModels "baize/app/freight/models"
	freightService "baize/app/freight/service"
	"encoding/json"
	"fmt"
	"strconv"
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

	if req.FormCode == "save_shipment_plan" {
		return s.saveShipmentPlan(req)
	}

	result := protocol.NewAgentResultV2("出货信息已确认", "已收到出货信息，可以继续生成出货计划。", []protocol.BlockItem{
		{Type: "markdown", Content: "请继续上传货物明细或提交保存出货计划。"},
	})
	s.Dao.InsertMessage(req.SessionID, "assistant", result.Summary, mustJSON(result), "qwen2.5:7b")
	return result
}

func (s *FormActionService) saveShipmentPlan(req *request.FormSubmitRequest) protocol.AgentResult {
	importReq, err := buildShipmentImportReq(req)
	if err != nil {
		result := protocol.NewAgentResultV2("保存出货计划失败", err.Error(), []protocol.BlockItem{
			{Type: "error", Title: "保存失败", Content: err.Error()},
		})
		s.Dao.InsertMessage(req.SessionID, "assistant", result.Summary, mustJSON(result), "qwen2.5:7b")
		return result
	}

	operator := req.OperatorName
	if operator == "" {
		operator = "agent"
	}
	detail, err := freightService.GetShipmentService().ImportShipment(importReq, operator)
	if err != nil {
		result := protocol.NewAgentResultV2("保存出货计划失败", err.Error(), []protocol.BlockItem{
			{Type: "error", Title: "保存失败", Content: err.Error()},
		})
		s.Dao.InsertMessage(req.SessionID, "assistant", result.Summary, mustJSON(result), "qwen2.5:7b")
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
			Type:    "markdown",
			Content: "刷新后台出货计划列表，或重新请求 `/freight/shipment/list?pageNum=1&pageSize=10` 即可看到这条记录。",
		},
	})
	s.Dao.InsertMessage(req.SessionID, "assistant", result.Summary, mustJSON(result), "qwen2.5:7b")
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
		s.Dao.InsertMessage(req.SessionID, "assistant", result.Summary, mustJSON(result), "qwen2.5:7b")
		return result
	}

	result := protocol.NewAgentResultV2("请先上传货物明细", "上传 Excel 后，分析结果中会出现保存表单。", []protocol.BlockItem{
		{Type: "markdown", Content: "客户端发起时系统会自动使用当前登录账号绑定的客户；后台发起时先创建未绑定计划，后续在出货计划列表绑定客户。"},
	})
	s.Dao.InsertMessage(req.SessionID, "assistant", result.Summary, mustJSON(result), "qwen2.5:7b")
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
