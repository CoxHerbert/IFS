package controller

import (
	"baize/app/agent/protocol"
	"baize/app/agent/request"
	"baize/app/agent/service"
	"baize/app/agent/skills"

	"github.com/gin-gonic/gin"
)

var shipmentAnalyzeService = service.GetShipmentAnalyzeService()

func AnalyzeShipment(c *gin.Context) {
	if file, err := c.FormFile("file"); err == nil {
		parsed, err := skills.ParseShipmentExcel(file)
		if err != nil {
			c.JSON(400, protocol.NewErrorResult(err.Error()))
			return
		}
		req := shipmentAnalyzeRequestFromParsed(parsed, c.PostForm("modelName"))
		c.JSON(200, shipmentAnalyzeService.Analyze(req))
		return
	}

	req := new(request.ShipmentAnalyzeRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, protocol.NewErrorResult("invalid shipment analyze payload"))
		return
	}
	c.JSON(200, shipmentAnalyzeService.Analyze(req))
}

func shipmentAnalyzeRequestFromParsed(parsed *skills.ParsedShipmentResult, modelName string) *request.ShipmentAnalyzeRequest {
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
