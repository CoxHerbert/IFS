package request

type CreateSessionRequest struct {
	Title     string `json:"title"`
	ModelName string `json:"modelName"`
}

type UpdateSessionTitleRequest struct {
	Title string `json:"title" binding:"required"`
}

type SendMessageRequest struct {
	SessionID int64  `json:"sessionId" binding:"required"`
	Message   string `json:"message" binding:"required"`
	ModelName string `json:"modelName"`
}
