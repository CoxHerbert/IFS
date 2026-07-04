package skills

import (
	"encoding/csv"
	"fmt"
	"math"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type ShipmentRow map[string]any

type ParsedShipmentItem struct {
	RowIndex    int
	SKU         string
	ProductName string
	Qty         int
	Length      float64
	Width       float64
	Height      float64
	Weight      float64
	CBM         float64
	Raw         ShipmentRow
}

type ParsedShipmentResult struct {
	FileName            string
	TotalQty            int
	TotalCBM            float64
	ContainerSuggestion string
	CargoList           []ParsedShipmentItem
}

var nonNumberPattern = regexp.MustCompile(`[^\d.\-]`)

func ParseShipmentExcel(file *multipart.FileHeader) (*ParsedShipmentResult, error) {
	rows, err := readTabularRows(file)
	if err != nil {
		return nil, err
	}
	items := normalizeRows(rows)
	if len(items) == 0 {
		return nil, fmt.Errorf("未识别到有效货物数据，请确认 Excel 包含数量、尺寸或 CBM")
	}

	var totalQty int
	var totalCBM float64
	for _, item := range items {
		totalQty += item.Qty
		totalCBM += item.CBM
	}
	totalCBM = round3(totalCBM)

	return &ParsedShipmentResult{
		FileName:            file.Filename,
		TotalQty:            totalQty,
		TotalCBM:            totalCBM,
		ContainerSuggestion: PlanContainer(totalCBM),
		CargoList:           items,
	}, nil
}

func readTabularRows(file *multipart.FileHeader) ([][]string, error) {
	reader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == ".csv" {
		csvReader := csv.NewReader(reader)
		csvReader.FieldsPerRecord = -1
		return csvReader.ReadAll()
	}
	if ext == ".xls" {
		return nil, fmt.Errorf("暂不支持旧版 .xls 文件，请先另存为 .xlsx 或 CSV 后再上传")
	}

	workbook, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, fmt.Errorf("Excel 文件读取失败：%w", err)
	}

	sheetMap := workbook.GetSheetMap()
	indexes := make([]int, 0, len(sheetMap))
	for index := range sheetMap {
		indexes = append(indexes, index)
	}
	sort.Ints(indexes)

	for _, index := range indexes {
		rows := rowsBySheetName(workbook, sheetMap[index])
		if len(rows) > 0 && !allRowsEmpty(rows) {
			return rows, nil
		}
	}

	for index := 1; index <= workbook.SheetCount+5; index++ {
		rows := rowsBySheetName(workbook, workbook.GetSheetName(index))
		if len(rows) > 0 && !allRowsEmpty(rows) {
			return rows, nil
		}
	}

	return nil, fmt.Errorf("Excel 中没有可读取的数据工作表，请确认文件为 .xlsx 且至少有一个非空工作表")
}

func rowsBySheetName(workbook *excelize.File, sheetName string) [][]string {
	if strings.TrimSpace(sheetName) == "" {
		return nil
	}
	return workbook.GetRows(sheetName)
}

func normalizeRows(rows [][]string) []ParsedShipmentItem {
	headerIndex := findHeaderIndex(rows)
	if headerIndex < 0 {
		return []ParsedShipmentItem{}
	}

	headers := normalizeHeaders(rows[headerIndex])
	items := make([]ParsedShipmentItem, 0)
	for i := headerIndex + 1; i < len(rows); i++ {
		row := rows[i]
		if isEmptyRow(row) {
			continue
		}
		raw := buildRawRow(headers, row)
		item, ok := buildShipmentItem(i+1, headers, row, raw)
		if ok {
			items = append(items, item)
		}
	}
	return items
}

func findHeaderIndex(rows [][]string) int {
	for i, row := range rows {
		score := 0
		for _, cell := range row {
			if fieldName(cell) != "" {
				score++
			}
		}
		if score >= 2 {
			return i
		}
	}
	return -1
}

func normalizeHeaders(row []string) []string {
	headers := make([]string, len(row))
	used := map[string]int{}
	for i, header := range row {
		field := fieldName(header)
		if field == "" {
			field = fmt.Sprintf("col%d", i+1)
		}
		used[field]++
		if used[field] > 1 {
			field = fmt.Sprintf("%s%d", field, used[field])
		}
		headers[i] = field
	}
	return headers
}

func buildRawRow(headers []string, row []string) ShipmentRow {
	raw := ShipmentRow{}
	for i, header := range headers {
		if i < len(row) {
			value := strings.TrimSpace(row[i])
			if value != "" {
				raw[header] = value
			}
		}
	}
	return raw
}

func buildShipmentItem(rowIndex int, headers []string, row []string, raw ShipmentRow) (ParsedShipmentItem, bool) {
	values := map[string]string{}
	for i, header := range headers {
		if i < len(row) {
			values[header] = strings.TrimSpace(row[i])
		}
	}

	qty := parseInt(values["qty"])
	if qty <= 0 {
		qty = 1
	}

	length := parseNumber(values["length"])
	width := parseNumber(values["width"])
	height := parseNumber(values["height"])
	if (length <= 0 || width <= 0 || height <= 0) && values["dimension"] != "" {
		if dim, ok, err := ParseDimensions(values["dimension"]); err == nil && ok {
			length = dim.Length
			width = dim.Width
			height = dim.Height
		}
	}

	cbm := parseNumber(values["cbm"])
	if cbm <= 0 && length > 0 && width > 0 && height > 0 {
		cbm = CalculateCBM(Dimensions{Length: length, Width: width, Height: height, Unit: "cm"}, qty)
	}
	if cbm <= 0 {
		return ParsedShipmentItem{}, false
	}

	return ParsedShipmentItem{
		RowIndex:    rowIndex,
		SKU:         values["sku"],
		ProductName: firstNonEmpty(values["productName"], values["name"]),
		Qty:         qty,
		Length:      round3(length),
		Width:       round3(width),
		Height:      round3(height),
		Weight:      round3(parseNumber(values["weight"])),
		CBM:         round3(cbm),
		Raw:         raw,
	}, true
}

func fieldName(header string) string {
	key := strings.ToLower(strings.TrimSpace(header))
	key = strings.ReplaceAll(key, " ", "")
	key = strings.ReplaceAll(key, "_", "")
	key = strings.ReplaceAll(key, "-", "")
	key = strings.ReplaceAll(key, "（cm）", "")
	key = strings.ReplaceAll(key, "(cm)", "")

	switch {
	case key == "sku" || key == "货号" || key == "型号" || key == "编码" || key == "产品编码" || key == "itemno":
		return "sku"
	case key == "品名" || key == "产品" || key == "产品名称" || key == "货物" || key == "货物名称" || key == "名称" || key == "product" || key == "productname" || key == "name" || strings.Contains(key, "品名"):
		return "productName"
	case key == "数量" || key == "箱数" || key == "件数" || key == "qty" || key == "quantity" || key == "ctn" || key == "ctns" || key == "carton" || key == "cartons" || strings.Contains(key, "数量") || strings.Contains(key, "箱数"):
		return "qty"
	case key == "长" || key == "长度" || key == "length" || key == "l" || strings.Contains(key, "长cm"):
		return "length"
	case key == "宽" || key == "宽度" || key == "width" || key == "w" || strings.Contains(key, "宽cm"):
		return "width"
	case key == "高" || key == "高度" || key == "height" || key == "h" || strings.Contains(key, "高cm"):
		return "height"
	case key == "尺寸" || key == "规格" || key == "外箱尺寸" || key == "包装尺寸" || key == "箱规" || key == "dimension" || key == "dimensions" || key == "size" || strings.Contains(key, "尺寸") || strings.Contains(key, "规格"):
		return "dimension"
	case key == "重量" || key == "毛重" || key == "单重" || key == "weight" || key == "gw" || key == "grossweight" || strings.Contains(key, "重量") || strings.Contains(key, "毛重"):
		return "weight"
	case key == "体积" || key == "方数" || key == "cbm" || key == "总体积" || key == "总方数" || strings.Contains(key, "cbm") || strings.Contains(key, "体积") || strings.Contains(key, "方数"):
		return "cbm"
	default:
		return ""
	}
}

func parseInt(value string) int {
	number := parseNumber(value)
	if number <= 0 {
		return 0
	}
	return int(math.Round(number))
}

func parseNumber(value string) float64 {
	value = strings.TrimSpace(strings.ReplaceAll(value, ",", ""))
	if value == "" {
		return 0
	}
	value = nonNumberPattern.ReplaceAllString(value, "")
	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return number
}

func isEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

func allRowsEmpty(rows [][]string) bool {
	for _, row := range rows {
		if !isEmptyRow(row) {
			return false
		}
	}
	return true
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func round3(value float64) float64 {
	return math.Round(value*1000) / 1000
}
