package service

import (
	"baize/app/agent/dao"
	"baize/app/agent/model"
	"strings"
)

type ContextService struct {
	Dao *dao.ChatDao
}

func NewContextService() *ContextService {
	return &ContextService{Dao: dao.GetChatDao()}
}

func (s *ContextService) BuildMessages(session *model.ChatSession, currentMessage string) []OllamaMessage {
	messages := []OllamaMessage{{Role: "system", Content: s.buildSystemPrompt(session)}}
	hasCurrentMessage := false
	for _, item := range s.Dao.SelectRecentMessages(session.ID, 20) {
		if item.Role == "user" || item.Role == "assistant" || item.Role == "system" || item.Role == "tool" {
			messages = append(messages, OllamaMessage{Role: item.Role, Content: item.Content})
			if item.Role == "user" && item.Content == currentMessage {
				hasCurrentMessage = true
			}
		}
	}
	if !hasCurrentMessage {
		messages = append(messages, OllamaMessage{Role: "user", Content: currentMessage})
	}
	return messages
}

func (s *ContextService) buildSystemPrompt(session *model.ChatSession) string {
	var builder strings.Builder
	builder.WriteString("你是 IFS 国际货代智能助手，请始终使用简洁、专业的中文回答。你擅长出货计划、CBM 计算、Excel 分析、报价支持和柜型建议。\n\n")
	builder.WriteString("长期记忆：\n")
	memories := s.Dao.SelectMemories(session.UserID)
	if len(memories) == 0 {
		builder.WriteString("无\n")
	} else {
		for _, memory := range memories {
			builder.WriteString("- ")
			builder.WriteString(memory.MemoryKey)
			builder.WriteString("：")
			builder.WriteString(memory.MemoryValue)
			builder.WriteString("\n")
		}
	}
	builder.WriteString("\n会话摘要：\n")
	if strings.TrimSpace(session.Summary) == "" {
		builder.WriteString("无\n")
	} else {
		builder.WriteString(session.Summary)
		builder.WriteString("\n")
	}
	return builder.String()
}
