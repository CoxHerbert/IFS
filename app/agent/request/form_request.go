package request

type FormSubmitRequest struct {
	SessionID    int64          `json:"sessionId" binding:"required"`
	FormCode     string         `json:"formCode" binding:"required"`
	Values       map[string]any `json:"values" binding:"required"`
	Source       string         `json:"-"`
	CustomerID   int64          `json:"-"`
	CustomerName string         `json:"-"`
	OperatorName string         `json:"-"`
}

type ActionExecuteRequest struct {
	SessionID  int64          `json:"sessionId" binding:"required"`
	ActionCode string         `json:"actionCode" binding:"required"`
	Payload    map[string]any `json:"payload"`
}
