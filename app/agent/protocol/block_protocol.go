package protocol

type AgentResult struct {
	Version string      `json:"version"`
	Type    string      `json:"type"`
	Title   string      `json:"title"`
	Summary string      `json:"summary"`
	Blocks  []BlockItem `json:"blocks"`
}

type BlockItem struct {
	Type          string         `json:"type"`
	Title         string         `json:"title,omitempty"`
	Content       string         `json:"content,omitempty"`
	Items         []MetricItem   `json:"items,omitempty"`
	Columns       []TableColumn  `json:"columns,omitempty"`
	Data          []any          `json:"data,omitempty"`
	Name          string         `json:"name,omitempty"`
	URL           string         `json:"url,omitempty"`
	FormCode      string         `json:"formCode,omitempty"`
	SubmitAPI     string         `json:"submitApi,omitempty"`
	Fields        []FormField    `json:"fields,omitempty"`
	InitialValues map[string]any `json:"initialValues,omitempty"`
	ActionCode    string         `json:"actionCode,omitempty"`
	Label         string         `json:"label,omitempty"`
	Payload       map[string]any `json:"payload,omitempty"`
}

type MetricItem struct {
	Label string `json:"label"`
	Value any    `json:"value"`
}

type TableColumn struct {
	Label string `json:"label"`
	Field string `json:"field"`
}

type FormField struct {
	Field       string       `json:"field"`
	Label       string       `json:"label"`
	Component   string       `json:"component"`
	Required    bool         `json:"required,omitempty"`
	Placeholder string       `json:"placeholder,omitempty"`
	Options     []FormOption `json:"options,omitempty"`
}

type FormOption struct {
	Label string `json:"label"`
	Value any    `json:"value"`
}

func NewAgentResult(title, summary string, blocks []BlockItem) AgentResult {
	return AgentResult{
		Version: "1.0",
		Type:    "agent_result",
		Title:   title,
		Summary: summary,
		Blocks:  blocks,
	}
}

func NewAgentResultV2(title, summary string, blocks []BlockItem) AgentResult {
	return AgentResult{
		Version: "2.0",
		Type:    "agent_result",
		Title:   title,
		Summary: summary,
		Blocks:  blocks,
	}
}

func NewErrorResult(message string) AgentResult {
	return NewAgentResult("Error", message, []BlockItem{
		{
			Type:    "error",
			Title:   "Error",
			Content: message,
		},
	})
}
