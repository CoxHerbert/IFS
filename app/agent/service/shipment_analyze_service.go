package service

import (
	"baize/app/agent/protocol"
	"baize/app/agent/request"
	"encoding/json"
	"fmt"
)

type ShipmentAnalyzeService struct {
	Ollama *OllamaService
}

var shipmentAnalyzeService = &ShipmentAnalyzeService{Ollama: NewOllamaService()}

func GetShipmentAnalyzeService() *ShipmentAnalyzeService {
	return shipmentAnalyzeService
}

func (s *ShipmentAnalyzeService) Analyze(req *request.ShipmentAnalyzeRequest) protocol.AgentResult {
	modelName := normalizeModelName(req.ModelName)
	ruleAdvice := buildRoutingAdvice(req.Summary.TotalCBM, req.Summary.ContainerSuggestion)
	modelMarkdown, err := s.Ollama.Chat(modelName, []OllamaMessage{
		{Role: "system", Content: buildShipmentAnalyzeSystemPrompt()},
		{Role: "user", Content: buildShipmentAnalyzeUserPrompt(req)},
	})
	if err != nil {
		modelMarkdown = "Ollama 调用失败：" + err.Error()
		return protocol.NewAgentResult("出货计划分析失败", modelMarkdown, []protocol.BlockItem{
			{Type: "error", Title: "错误", Content: modelMarkdown},
		})
	}

	markdown := ruleAdvice + "\n\n" + modelMarkdown
	summary := fmt.Sprintf("本次共 %d 箱，总体积 %s CBM，建议方案：%s。",
		req.Summary.TotalQty,
		formatNumber(req.Summary.TotalCBM),
		req.Summary.ContainerSuggestion,
	)

	return protocol.NewAgentResult("出货计划分析结果", summary, []protocol.BlockItem{
		{
			Type: "metrics",
			Items: []protocol.MetricItem{
				{Label: "总箱数", Value: req.Summary.TotalQty},
				{Label: "总体积", Value: formatNumber(req.Summary.TotalCBM) + " CBM"},
				{Label: "建议方案", Value: req.Summary.ContainerSuggestion},
			},
		},
		buildShipmentTableBlock(req.CargoList),
		buildSaveShipmentFormBlock(req),
		{Type: "markdown", Content: markdown},
	})
}

func buildRoutingAdvice(totalCBM float64, suggestion string) string {
	if totalCBM > 0 && totalCBM < 15 {
		return fmt.Sprintf("### 规则判断\n当前总体积为 %s CBM，属于小票货，优先建议 **LCL 拼箱**。除非客户明确要求包柜、货值/时效/敏感货属性需要单独控货，通常不建议使用整柜。", formatNumber(totalCBM))
	}
	return fmt.Sprintf("### 规则判断\n当前总体积为 %s CBM，系统建议方案为 **%s**。", formatNumber(totalCBM), suggestion)
}

func buildShipmentAnalyzeSystemPrompt() string {
	return "你是 IFS 国际货运智能助手。后端已经完成 Excel 解析、标准化和 CBM 计算。\n\n必须遵守：\n1. 不要重新计算 CBM。\n2. totalCBM < 15 时，默认优先建议 LCL 拼箱，不要建议 20GP/40GP/40HQ，除非用户明确要求包柜。\n3. 15 <= totalCBM <= 28 时，可以建议 20GP，但也要提醒可按价格比较 LCL 与 FCL。\n4. 只基于输入里的 ContainerSuggestion 做分析，不要编造柜型容量。\n5. 返回简洁中文 markdown。"
}

func buildShipmentAnalyzeUserPrompt(req *request.ShipmentAnalyzeRequest) string {
	payload, _ := json.MarshalIndent(req, "", "  ")
	return "以下是后端解析并计算后的结构化出货数据，请基于这些结果给出建议。注意：如果 containerSuggestion 是 LCL 拼箱，请围绕拼箱说明，不要反向推荐整柜：\n\n```json\n" + string(payload) + "\n```"
}

func buildShipmentTableBlock(list []request.StandardCargoItem) protocol.BlockItem {
	data := make([]any, 0, len(list))
	for _, item := range list {
		data = append(data, map[string]any{
			"rowIndex":    item.RowIndex,
			"sku":         item.SKU,
			"productName": item.ProductName,
			"qty":         item.Qty,
			"length":      formatNumber(item.Length),
			"width":       formatNumber(item.Width),
			"height":      formatNumber(item.Height),
			"weight":      formatNumber(item.Weight),
			"cbm":         formatNumber(item.CBM),
		})
	}
	return protocol.BlockItem{
		Type:  "table",
		Title: "服务端识别的货物明细",
		Columns: []protocol.TableColumn{
			{Label: "行号", Field: "rowIndex"},
			{Label: "货号", Field: "sku"},
			{Label: "品名", Field: "productName"},
			{Label: "数量", Field: "qty"},
			{Label: "长(cm)", Field: "length"},
			{Label: "宽(cm)", Field: "width"},
			{Label: "高(cm)", Field: "height"},
			{Label: "重量", Field: "weight"},
			{Label: "CBM", Field: "cbm"},
		},
		Data: data,
	}
}

func buildSaveShipmentFormBlock(req *request.ShipmentAnalyzeRequest) protocol.BlockItem {
	return protocol.BlockItem{
		Type:      "form",
		Title:     "保存为正式出货计划",
		FormCode:  "save_shipment_plan",
		SubmitAPI: "/api/agent/form/submit",
		Fields: []protocol.FormField{
			{Field: "orderNo", Label: "客户订单号", Component: "input"},
			{Field: "pol", Label: "起运港", Component: "input"},
			{Field: "pod", Label: "目的港", Component: "input"},
			{Field: "plannedEtd", Label: "计划开船", Component: "date"},
			{Field: "plannedEta", Label: "计划到港", Component: "date"},
			{Field: "preferredType", Label: "偏好方案", Component: "select", Options: []protocol.FormOption{
				{Label: "自动推荐", Value: ""},
				{Label: "LCL 拼箱", Value: "LCL"},
				{Label: "20GP", Value: "20GP"},
				{Label: "40GP", Value: "40GP"},
				{Label: "40HQ", Value: "40HQ"},
			}},
			{Field: "remark", Label: "备注", Component: "textarea"},
		},
		InitialValues: map[string]any{
			"fileName":      req.FileName,
			"cargoList":     buildShipmentImportCargo(req.CargoList),
			"preferredType": preferredTypeFromSuggestion(req.Summary.ContainerSuggestion),
			"remark":        "由 IFS Agent 解析 Excel 后生成；客户归属按发起端规则处理",
		},
	}
}

func preferredTypeFromSuggestion(suggestion string) string {
	switch suggestion {
	case "LCL 拼箱":
		return "LCL"
	case "1×20GP":
		return "20GP"
	case "1×40GP":
		return "40GP"
	case "1×40HQ":
		return "40HQ"
	default:
		return ""
	}
}

func buildShipmentImportCargo(list []request.StandardCargoItem) []map[string]any {
	cargoList := make([]map[string]any, 0, len(list))
	for _, item := range list {
		name := item.ProductName
		if name == "" {
			name = item.SKU
		}
		if name == "" {
			name = "Excel 第 " + formatNumber(float64(item.RowIndex)) + " 行货物"
		}
		cargoList = append(cargoList, map[string]any{
			"sku":       item.SKU,
			"cargoName": name,
			"quantity":  item.Qty,
			"cartons":   item.Qty,
			"weightKg":  item.Weight,
			"volumeCbm": item.CBM,
			"lengthCm":  item.Length,
			"widthCm":   item.Width,
			"heightCm":  item.Height,
		})
	}
	return cargoList
}
