package request

type ShipmentAnalyzeRequest struct {
	Summary   ShipmentSummary     `json:"summary"`
	CargoList []StandardCargoItem `json:"cargoList"`
	ModelName string              `json:"modelName"`
	FileName  string              `json:"fileName"`
}

type ShipmentSummary struct {
	TotalQty            int     `json:"totalQty"`
	TotalCBM            float64 `json:"totalCBM"`
	ContainerSuggestion string  `json:"containerSuggestion"`
}

type StandardCargoItem struct {
	RowIndex    int     `json:"rowIndex"`
	SKU         string  `json:"sku"`
	ProductName string  `json:"productName"`
	Qty         int     `json:"qty"`
	Length      float64 `json:"length"`
	Width       float64 `json:"width"`
	Height      float64 `json:"height"`
	Weight      float64 `json:"weight"`
	CBM         float64 `json:"cbm"`
	Raw         any     `json:"raw,omitempty"`
}
