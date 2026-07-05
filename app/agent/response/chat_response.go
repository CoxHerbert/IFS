package response

import "baize/app/agent/protocol"

type SendMessageResponse struct {
	MessageID int64                `json:"messageId"`
	SessionID int64                `json:"sessionId"`
	Result    protocol.AgentResult `json:"result"`
}

type ModelOption struct {
	Label       string `json:"label"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
}
