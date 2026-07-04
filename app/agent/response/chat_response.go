package response

import "baize/app/agent/protocol"

type SendMessageResponse struct {
	MessageID int64                `json:"messageId"`
	SessionID int64                `json:"sessionId"`
	Result    protocol.AgentResult `json:"result"`
}
