package service

import (
	"baize/app/agent/dao"
	"baize/app/agent/model"
	"baize/app/agent/protocol"
	"baize/app/agent/request"
	"baize/app/agent/response"
	"baize/app/agent/skills"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"
)

const defaultModelName = "qwen2.5:7b"

type ChatService struct {
	Dao     *dao.ChatDao
	Context *ContextService
	Ollama  *OllamaService
}

var chatService = &ChatService{
	Dao:     dao.GetChatDao(),
	Context: NewContextService(),
	Ollama:  NewOllamaService(),
}

var quantityPattern = regexp.MustCompile(`(?i)(?:(?:\x{6570}\x{91CF})\s*(\d+)|qty\s*(\d+)|(\d+)\s*(?:\x{7BB1}|\x{4EF6}|ctn|ctns|carton|cartons))`)

var modelOptions = []*response.ModelOption{
	{Label: "Qwen 2.5 7B", Value: "qwen2.5:7b", Description: "默认模型，适合日常货运问答和出货分析。", Default: true},
	{Label: "Qwen 2.5 14B", Value: "qwen2.5:14b", Description: "更强的推理模型，适合复杂方案分析。"},
	{Label: "Llama 3.1 8B", Value: "llama3.1:8b", Description: "通用对话模型，可作为备选。"},
	{Label: "DeepSeek R1 8B", Value: "deepseek-r1:8b", Description: "偏推理场景，可用于复杂计算说明。"},
}

func GetChatService() *ChatService {
	return chatService
}

func (s *ChatService) ListModels() []*response.ModelOption {
	return modelOptions
}

func (s *ChatService) CreateSession(userID int64, req *request.CreateSessionRequest) *model.ChatSessionVO {
	modelName := normalizeModelName(req.ModelName)
	title := strings.TrimSpace(req.Title)
	if title == "" {
		title = "新对话"
	}
	id := s.Dao.InsertSession(userID, title, modelName)
	return &model.ChatSessionVO{ID: id, Title: title, ModelName: modelName}
}

func (s *ChatService) ListSessions(userID int64) []*model.ChatSessionVO {
	list := s.Dao.SelectSessions(userID)
	result := make([]*model.ChatSessionVO, 0, len(list))
	for _, item := range list {
		result = append(result, sessionVO(item))
	}
	return result
}

func (s *ChatService) ListMessages(userID int64, sessionID int64) []*model.ChatMessageVO {
	if s.Dao.SelectSession(sessionID, userID) == nil {
		return []*model.ChatMessageVO{}
	}
	list := s.Dao.SelectMessages(sessionID)
	result := make([]*model.ChatMessageVO, 0, len(list))
	for _, item := range list {
		result = append(result, messageVO(item))
	}
	return result
}

func (s *ChatService) UpdateSessionTitle(userID int64, sessionID int64, req *request.UpdateSessionTitleRequest) error {
	if sessionID == 0 {
		return fmt.Errorf("缺少会话 ID")
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return fmt.Errorf("标题不能为空")
	}
	if len([]rune(title)) > 80 {
		return fmt.Errorf("标题不能超过 80 个字符")
	}
	if !s.Dao.UpdateSessionTitle(sessionID, userID, title) {
		return fmt.Errorf("未找到对应会话")
	}
	return nil
}

func (s *ChatService) DeleteSession(userID int64, sessionID int64) error {
	if sessionID == 0 {
		return fmt.Errorf("缺少会话 ID")
	}
	if !s.Dao.DeleteSession(sessionID, userID) {
		return fmt.Errorf("未找到对应会话")
	}
	return nil
}

func (s *ChatService) Send(userID int64, req *request.SendMessageRequest) (*response.SendMessageResponse, error) {
	session := s.Dao.SelectSession(req.SessionID, userID)
	if session == nil {
		return nil, fmt.Errorf("未找到对应会话")
	}
	modelName := normalizeModelName(req.ModelName)
	if strings.TrimSpace(req.ModelName) == "" {
		modelName = normalizeModelName(session.ModelName)
	}

	s.Dao.InsertMessage(session.ID, "user", req.Message, "", modelName)
	if formResult := s.tryFormBlock(req.Message); formResult != nil {
		messageID := s.saveAssistant(session.ID, formResult.Summary, *formResult, modelName)
		return &response.SendMessageResponse{MessageID: messageID, SessionID: session.ID, Result: *formResult}, nil
	}

	result := s.tryLocalSkill(req.Message)
	if result == nil {
		answer, err := s.Ollama.Chat(modelName, s.Context.BuildMessages(session, req.Message))
		if err != nil {
			errorResult := protocol.NewErrorResult("Ollama 调用失败：" + err.Error())
			messageID := s.saveAssistant(session.ID, errorResult.Summary, errorResult, modelName)
			return &response.SendMessageResponse{MessageID: messageID, SessionID: session.ID, Result: errorResult}, nil
		}
		normalResult := protocol.NewAgentResult("智能助手回复", answer, []protocol.BlockItem{{Type: "markdown", Content: answer}})
		messageID := s.saveAssistant(session.ID, answer, normalResult, modelName)
		return &response.SendMessageResponse{MessageID: messageID, SessionID: session.ID, Result: normalResult}, nil
	}

	messageID := s.saveAssistant(session.ID, result.Summary, *result, modelName)
	return &response.SendMessageResponse{MessageID: messageID, SessionID: session.ID, Result: *result}, nil
}

func (s *ChatService) tryFormBlock(message string) *protocol.AgentResult {
	lower := strings.ToLower(message)
	intent := strings.Contains(message, "出货") ||
		strings.Contains(message, "出一票货") ||
		strings.Contains(message, "报价") ||
		strings.Contains(message, "发美国") ||
		strings.Contains(message, "出货计划") ||
		strings.Contains(lower, "shipment") ||
		strings.Contains(lower, "quote")
	if !intent {
		return nil
	}

	initialValues := map[string]any{}
	if strings.Contains(message, "宁波") {
		initialValues["pol"] = "宁波"
	}
	if strings.Contains(message, "洛杉矶") || strings.Contains(strings.ToLower(message), "los angeles") || strings.Contains(message, "美国") {
		initialValues["pod"] = "洛杉矶"
	}

	result := protocol.NewAgentResultV2("请补充出货信息", "为了生成出货计划，请补充以下信息。", []protocol.BlockItem{
		{
			Type:      "form",
			Title:     "出货信息表单",
			FormCode:  "shipment_info",
			SubmitAPI: "/api/agent/form/submit",
			Fields: []protocol.FormField{
				{Field: "pol", Label: "起运港", Component: "input", Required: true, Placeholder: "例如：宁波"},
				{Field: "pod", Label: "目的港", Component: "input", Required: true, Placeholder: "例如：洛杉矶"},
				{Field: "cargoReadyDate", Label: "货好时间", Component: "date"},
				{Field: "tradeTerm", Label: "贸易条款", Component: "select", Required: true, Options: []protocol.FormOption{
					{Label: "FOB", Value: "FOB"},
					{Label: "EXW", Value: "EXW"},
					{Label: "CIF", Value: "CIF"},
					{Label: "DDP", Value: "DDP"},
				}},
				{Field: "remark", Label: "备注", Component: "textarea"},
			},
			InitialValues: initialValues,
		},
	})
	return &result
}

func (s *ChatService) AnalyzeShipment(userID int64, sessionID int64, req *request.ShipmentAnalyzeRequest) (*response.SendMessageResponse, error) {
	session := s.Dao.SelectSession(sessionID, userID)
	if session == nil {
		return nil, fmt.Errorf("未找到对应会话")
	}
	modelName := normalizeModelName(req.ModelName)
	if strings.TrimSpace(req.ModelName) == "" {
		modelName = normalizeModelName(session.ModelName)
	}

	fileName := strings.TrimSpace(req.FileName)
	if fileName == "" {
		fileName = "Excel 文件"
	}
	userContent := fmt.Sprintf("上传出货文件：%s。前端已解析 %d 行货物，总箱数 %d，总体积 %s CBM，建议柜型 %s。",
		fileName,
		len(req.CargoList),
		req.Summary.TotalQty,
		formatNumber(req.Summary.TotalCBM),
		req.Summary.ContainerSuggestion,
	)
	s.Dao.InsertMessage(session.ID, "user", userContent, "", modelName)

	result := shipmentAnalyzeService.Analyze(req)
	messageID := s.saveAssistant(session.ID, result.Summary, result, modelName)
	return &response.SendMessageResponse{MessageID: messageID, SessionID: session.ID, Result: result}, nil
}

func (s *ChatService) AnalyzeShipmentFile(userID int64, sessionID int64, file *multipart.FileHeader, modelName string) (*response.SendMessageResponse, error) {
	session := s.Dao.SelectSession(sessionID, userID)
	if session == nil {
		return nil, fmt.Errorf("未找到对应会话")
	}
	if strings.TrimSpace(modelName) == "" {
		modelName = normalizeModelName(session.ModelName)
	} else {
		modelName = normalizeModelName(modelName)
	}

	parsed, err := skills.ParseShipmentExcel(file)
	if err != nil {
		return nil, err
	}
	req := shipmentRequestFromParsed(parsed, modelName)

	userContent := fmt.Sprintf("上传出货文件：%s。后端已解析 %d 行货物，总箱数 %d，总体积 %s CBM，建议柜型 %s。",
		parsed.FileName,
		len(parsed.CargoList),
		parsed.TotalQty,
		formatNumber(parsed.TotalCBM),
		parsed.ContainerSuggestion,
	)
	s.Dao.InsertMessage(session.ID, "user", userContent, "", modelName)

	result := shipmentAnalyzeService.Analyze(req)
	messageID := s.saveAssistant(session.ID, result.Summary, result, modelName)
	return &response.SendMessageResponse{MessageID: messageID, SessionID: session.ID, Result: result}, nil
}

func shipmentRequestFromParsed(parsed *skills.ParsedShipmentResult, modelName string) *request.ShipmentAnalyzeRequest {
	list := make([]request.StandardCargoItem, 0, len(parsed.CargoList))
	for _, item := range parsed.CargoList {
		list = append(list, request.StandardCargoItem{
			RowIndex:    item.RowIndex,
			SKU:         item.SKU,
			ProductName: item.ProductName,
			Qty:         item.Qty,
			Length:      item.Length,
			Width:       item.Width,
			Height:      item.Height,
			Weight:      item.Weight,
			CBM:         item.CBM,
			Raw:         item.Raw,
		})
	}
	return &request.ShipmentAnalyzeRequest{
		Summary: request.ShipmentSummary{
			TotalQty:            parsed.TotalQty,
			TotalCBM:            parsed.TotalCBM,
			ContainerSuggestion: parsed.ContainerSuggestion,
		},
		CargoList: list,
		ModelName: modelName,
		FileName:  parsed.FileName,
	}
}

func (s *ChatService) tryLocalSkill(message string) *protocol.AgentResult {
	dim, hasDimensions, err := skills.ParseDimensions(message)
	if err != nil {
		result := protocol.NewErrorResult(err.Error())
		return &result
	}
	quantity, hasQuantity := parseQuantity(message)
	if !hasDimensions || !hasQuantity {
		return nil
	}
	totalCBM := skills.CalculateCBM(dim, quantity)
	container := skills.PlanContainer(totalCBM)
	summary := fmt.Sprintf("\u8be5\u6279\u8d27\u7269\u603b\u4f53\u79ef\u4e3a %s CBM\uff0c\u5efa\u8bae\u67dc\u578b\uff1a%s\u3002", formatNumber(totalCBM), container)
	result := protocol.NewAgentResult("\u8ba1\u7b97\u7ed3\u679c", summary, []protocol.BlockItem{
		{
			Type: "metrics",
			Items: []protocol.MetricItem{
				{Label: "\u957f", Value: formatNumber(dim.Length) + " cm"},
				{Label: "\u5bbd", Value: formatNumber(dim.Width) + " cm"},
				{Label: "\u9ad8", Value: formatNumber(dim.Height) + " cm"},
				{Label: "\u7bb1\u6570", Value: quantity},
				{Label: "\u603b\u4f53\u79ef", Value: formatNumber(totalCBM) + " CBM"},
				{Label: "\u5efa\u8bae\u67dc\u578b", Value: container},
			},
		},
		{Type: "markdown", Content: "\u8ba1\u7b97\u516c\u5f0f\uff1a\u957f \u00d7 \u5bbd \u00d7 \u9ad8 \u00f7 1000000 \u00d7 \u7bb1\u6570"},
	})
	return &result
}

func (s *ChatService) saveAssistant(sessionID int64, content string, result protocol.AgentResult, modelName string) int64 {
	blockJSON, _ := json.Marshal(result)
	return s.Dao.InsertMessage(sessionID, "assistant", content, string(blockJSON), modelName)
}

func parseQuantity(message string) (int, bool) {
	match := quantityPattern.FindStringSubmatch(message)
	if len(match) == 0 {
		return 0, false
	}
	for i := 1; i < len(match); i++ {
		if match[i] == "" {
			continue
		}
		quantity, err := strconv.Atoi(match[i])
		return quantity, err == nil && quantity > 0
	}
	return 0, false
}

func normalizeModelName(modelName string) string {
	modelName = strings.TrimSpace(modelName)
	if modelName == "" {
		return defaultModelName
	}
	return modelName
}

func sessionVO(item *model.ChatSession) *model.ChatSessionVO {
	return &model.ChatSessionVO{
		ID:        item.ID,
		Title:     item.Title,
		ModelName: item.ModelName,
		Summary:   item.Summary,
		UpdatedAt: item.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func messageVO(item *model.ChatMessage) *model.ChatMessageVO {
	vo := &model.ChatMessageVO{
		ID:        item.ID,
		SessionID: item.SessionID,
		Role:      item.Role,
		Content:   item.Content,
		CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if strings.TrimSpace(item.BlockResult) != "" {
		var block any
		if json.Unmarshal([]byte(item.BlockResult), &block) == nil {
			vo.BlockResult = block
		}
	}
	return vo
}

func formatNumber(value float64) string {
	if value == float64(int64(value)) {
		return strconv.FormatInt(int64(value), 10)
	}
	return strconv.FormatFloat(value, 'f', -1, 64)
}
