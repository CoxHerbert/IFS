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
	builder.WriteString("You are the IFS international freight intelligent assistant. Answer in concise Chinese. You are good at shipment planning, CBM calculation, Excel analysis, quotation support, and container suggestions.\n\n")
	builder.WriteString("Long-term memory:\n")
	memories := s.Dao.SelectMemories(session.UserID)
	if len(memories) == 0 {
		builder.WriteString("none\n")
	} else {
		for _, memory := range memories {
			builder.WriteString("- ")
			builder.WriteString(memory.MemoryKey)
			builder.WriteString(": ")
			builder.WriteString(memory.MemoryValue)
			builder.WriteString("\n")
		}
	}
	builder.WriteString("\nConversation summary:\n")
	if strings.TrimSpace(session.Summary) == "" {
		builder.WriteString("none\n")
	} else {
		builder.WriteString(session.Summary)
		builder.WriteString("\n")
	}
	return builder.String()
}
