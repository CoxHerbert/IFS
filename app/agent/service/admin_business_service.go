package service

import (
	"baize/app/agent/protocol"
	"baize/app/agent/request"
	customerModels "baize/app/customer/models"
	customerService "baize/app/customer/service"
	freightModels "baize/app/freight/models"
	freightService "baize/app/freight/service"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var moneyIntentPattern = regexp.MustCompile(`(?i)(\d+(?:\.\d+)?)\s*(万|元|rmb|cny)?`)
var explicitMoneyPattern = regexp.MustCompile(`(?i)(\d+(?:\.\d+)?)\s*(万|元|rmb|cny)`)
var shipmentNoPattern = regexp.MustCompile(`(?i)\b(?:SP|SH|IFS)[A-Z0-9_-]+\b`)
var receiptNoPattern = regexp.MustCompile(`(?i)\bRC[A-Z0-9_-]+\b`)

func tryAdminBusiness(req *request.SendMessageRequest) *protocol.AgentResult {
	if req.Source != "admin" { return nil }
	message := strings.TrimSpace(req.Message)
	if strings.Contains(message,"核销") && strings.Contains(message,"收款") { if !hasAgentPermission(req,"freight:receipt:add") { return permissionDenied() }; return buildAllocationDraft(message) }
	if strings.Contains(message,"收款") && (strings.Contains(message,"新增") || strings.Contains(message,"收到") || strings.Contains(message,"登记")) { if !hasAgentPermission(req,"freight:receipt:add") { return permissionDenied() }; return buildReceiptDraft(message) }
	if (strings.Contains(message,"出货计划") || strings.Contains(message,"出货单")) && isQueryIntent(message) { if !hasAgentPermission(req,"freight:shipment:list") { return permissionDenied() }; return queryShipments(message,req) }
	if strings.Contains(message,"收款") && isQueryIntent(message) { if !hasAgentPermission(req,"freight:receipt:list") { return permissionDenied() }; return queryReceipts(message,req) }
	if strings.Contains(message,"客户") && isQueryIntent(message) { if !hasAgentPermission(req,"customer:customer:list") { return permissionDenied() }; return queryCustomers(message,req) }
	return nil
}

func buildAllocationDraft(message string) *protocol.AgentResult {
	receiptNo,shipmentNo:="",""; if receiptNoPattern.MatchString(message){receiptNo=receiptNoPattern.FindString(message)};if shipmentNoPattern.MatchString(message){shipmentNo=shipmentNoPattern.FindString(message)}
	amount:=0.0; matches:=explicitMoneyPattern.FindAllStringSubmatch(message,-1); if len(matches)>0 { m:=matches[0]; amount,_=strconv.ParseFloat(m[1],64);if len(m)>2&&m[2]=="万"{amount*=10000} }
	r:=protocol.NewAgentResultV2("确认追加核销","这是写账操作，请核对后提交。",[]protocol.BlockItem{{Type:"form",Title:"收款核销确认",FormCode:"admin_receipt_allocate",SubmitAPI:"/agent/chat/form/submit",Fields:[]protocol.FormField{{Field:"receiptNo",Label:"收款单号",Component:"input",Required:true},{Field:"shipmentNo",Label:"出货计划",Component:"input",Required:true},{Field:"allocatedAmount",Label:"核销金额",Component:"number",Required:true}},InitialValues:map[string]any{"receiptNo":receiptNo,"shipmentNo":shipmentNo,"allocatedAmount":amount}}});return &r
}

func hasAgentPermission(req *request.SendMessageRequest, permission string) bool { for _,p:=range req.Permissions { if p=="*:*:*" || p==permission { return true } }; return false }
func permissionDenied() *protocol.AgentResult { r:=protocol.NewErrorResult("没有执行该业务操作的权限"); return &r }

func isQueryIntent(s string) bool { return strings.Contains(s,"查") || strings.Contains(s,"哪些") || strings.Contains(s,"多少") || strings.Contains(s,"最近") || strings.Contains(s,"列表") }
func cleanKeyword(s string) string {
	r:=strings.NewReplacer("查询","","查一下","","查","","客户","","出货计划","","出货单","","收款单","","收款","","最近的","","最近","","有哪些","","列表","")
	return strings.Trim(strings.TrimSpace(r.Replace(s)),"：:，,。？? ")
}

func queryCustomers(message string, req *request.SendMessageRequest) *protocol.AgentResult {
	q:=&customerModels.CustomerDQL{CustomerName:cleanKeyword(message)}; q.Size=10; q.Limit=" limit 0,10"
	if !req.CanManageAll { q.SalesUserId=req.OperatorID }
	list,total:=customerService.GetCustomerService().SelectCustomerList(q); data:=make([]any,0,len(list))
	for _,v:=range list { data=append(data,map[string]any{"customerNo":v.CustomerNo,"customerName":v.CustomerName,"companyName":v.CompanyName,"salesUserName":v.SalesUserName,"status":v.Status}) }
	r:=protocol.NewAgentResultV2("客户查询",fmt.Sprintf("查询到 %d 条客户记录。",*total),[]protocol.BlockItem{{Type:"table",Columns:[]protocol.TableColumn{{Label:"客户编号",Field:"customerNo"},{Label:"客户名称",Field:"customerName"},{Label:"公司",Field:"companyName"},{Label:"业务员",Field:"salesUserName"},{Label:"状态",Field:"status"}},Data:data}}); return &r
}

func queryShipments(message string, req *request.SendMessageRequest) *protocol.AgentResult {
	q:=&freightModels.ShipmentPlanDQL{}; keyword:=cleanKeyword(message); if shipmentNoPattern.MatchString(message) { q.ShipmentNo=shipmentNoPattern.FindString(message) } else { q.CustomerName=keyword }
	if !req.CanManageAll { q.SalesUserId=req.OperatorID }; q.Size=10; q.Limit=" limit 0,10"
	list,total:=freightService.GetShipmentService().SelectShipmentList(q); data:=make([]any,0,len(list))
	for _,v:=range list { data=append(data,map[string]any{"shipmentNo":v.ShipmentNo,"customerName":v.CustomerName,"route":v.Pol+" → "+v.Pod,"status":v.Status,"paymentStatus":v.PaymentStatus,"paymentAmount":v.PaymentAmount}) }
	r:=protocol.NewAgentResultV2("出货计划查询",fmt.Sprintf("查询到 %d 条出货计划。",*total),[]protocol.BlockItem{
		{Type:"table",Columns:[]protocol.TableColumn{{Label:"计划编号",Field:"shipmentNo"},{Label:"客户",Field:"customerName"},{Label:"航线",Field:"route"},{Label:"状态",Field:"status"},{Label:"收款状态",Field:"paymentStatus"},{Label:"已核销",Field:"paymentAmount"}},Data:data},
		{Type:"navigate",Label:"打开出货计划列表",URL:"/freight/shipment"},
	}); return &r
}

func queryReceipts(message string, req *request.SendMessageRequest) *protocol.AgentResult {
	q:=&freightModels.ReceiptDQL{CustomerName:cleanKeyword(message)}; if !req.CanManageAll { q.SalesUserId=req.OperatorID }; q.Size=10; q.Limit=" limit 0,10"
	list,total:=freightService.GetReceiptService().SelectList(q); data:=make([]any,0,len(list))
	for _,v:=range list { data=append(data,map[string]any{"receiptNo":v.ReceiptNo,"customerName":v.CustomerName,"amount":v.Amount,"currency":v.Currency,"allocatedAmount":v.AllocatedAmount,"status":v.Status,"receiptTime":v.ReceiptTime}) }
	r:=protocol.NewAgentResultV2("收款单查询",fmt.Sprintf("查询到 %d 条收款记录。",*total),[]protocol.BlockItem{{Type:"table",Columns:[]protocol.TableColumn{{Label:"收款单号",Field:"receiptNo"},{Label:"客户",Field:"customerName"},{Label:"金额",Field:"amount"},{Label:"币种",Field:"currency"},{Label:"已核销",Field:"allocatedAmount"},{Label:"状态",Field:"status"},{Label:"收款时间",Field:"receiptTime"}},Data:data}}); return &r
}

func buildReceiptDraft(message string) *protocol.AgentResult {
	amount:=0.0; if m:=moneyIntentPattern.FindStringSubmatch(message);len(m)>1 { amount,_=strconv.ParseFloat(m[1],64); if len(m)>2&&m[2]=="万" { amount*=10000 } }
	shipmentNo:=""; if shipmentNoPattern.MatchString(message) { shipmentNo=shipmentNoPattern.FindString(message) }
	initial:=map[string]any{"amount":amount,"currency":"CNY","paymentMethod":"BANK_TRANSFER","shipmentNos":shipmentNo,"allocationAmounts":amount}
	r:=protocol.NewAgentResultV2("确认新增收款","这是写账操作，请核对后提交。多计划时，计划编号和核销金额都用英文逗号分隔。",[]protocol.BlockItem{{Type:"form",Title:"新增收款确认",FormCode:"admin_receipt_create",SubmitAPI:"/agent/chat/form/submit",Fields:[]protocol.FormField{
		{Field:"customerName",Label:"客户名称",Component:"input",Required:true},{Field:"amount",Label:"收款金额",Component:"number",Required:true},{Field:"currency",Label:"币种",Component:"input",Required:true},{Field:"paymentMethod",Label:"收款方式",Component:"select",Required:true,Options:[]protocol.FormOption{{Label:"银行转账",Value:"BANK_TRANSFER"},{Label:"现金",Value:"CASH"},{Label:"其他",Value:"OTHER"}}},{Field:"shipmentNos",Label:"出货计划",Component:"input",Placeholder:"多个编号用英文逗号分隔"},{Field:"allocationAmounts",Label:"核销金额",Component:"input",Placeholder:"与出货计划一一对应，用英文逗号分隔"},{Field:"remark",Label:"备注",Component:"textarea"}},InitialValues:initial}}); return &r
}
