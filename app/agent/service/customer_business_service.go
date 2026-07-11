package service

import (
	"baize/app/agent/protocol"
	"baize/app/agent/request"
	freightModels "baize/app/freight/models"
	freightService "baize/app/freight/service"
	"fmt"
	"strconv"
	"strings"
)

func tryCustomerBusiness(req *request.SendMessageRequest) *protocol.AgentResult {
	if req.Source!="customer" || req.CustomerID==0{return nil};message:=strings.TrimSpace(req.Message)
	if (strings.Contains(message,"付款")||strings.Contains(message,"凭证"))&&(strings.Contains(message,"申报")||strings.Contains(message,"提交")||strings.Contains(message,"上传")){return buildPaymentDeclarationDraft(message)}
	if (strings.Contains(message,"收款")||strings.Contains(message,"核销")||strings.Contains(message,"付款记录"))&&isQueryIntent(message){return queryCustomerReceipts(req)}
	if (strings.Contains(message,"出货")||strings.Contains(message,"货运计划"))&&isQueryIntent(message){return queryCustomerShipments(message,req)}
	return nil
}

func queryCustomerShipments(message string,req *request.SendMessageRequest)*protocol.AgentResult{
	q:=&freightModels.ShipmentPlanDQL{CustomerId:req.CustomerID};if shipmentNoPattern.MatchString(message){q.ShipmentNo=shipmentNoPattern.FindString(message)};q.Size=20;q.Limit=" limit 0,20"
	list,total:=freightService.GetShipmentService().SelectShipmentList(q);data:=make([]any,0,len(list));for _,v:=range list{data=append(data,map[string]any{"shipmentNo":v.ShipmentNo,"route":v.Pol+" → "+v.Pod,"plannedEtd":v.PlannedEtd,"plannedEta":v.PlannedEta,"status":v.Status,"paymentStatus":v.PaymentStatus,"paymentAmount":v.PaymentAmount})}
	r:=protocol.NewAgentResultV2("我的出货计划",fmt.Sprintf("查询到 %d 条出货计划。",*total),[]protocol.BlockItem{{Type:"table",Columns:[]protocol.TableColumn{{Label:"计划编号",Field:"shipmentNo"},{Label:"航线",Field:"route"},{Label:"ETD",Field:"plannedEtd"},{Label:"ETA",Field:"plannedEta"},{Label:"状态",Field:"status"},{Label:"付款状态",Field:"paymentStatus"},{Label:"已核销",Field:"paymentAmount"}},Data:data}});return &r
}

func queryCustomerReceipts(req *request.SendMessageRequest)*protocol.AgentResult{
	q:=&freightModels.ReceiptDQL{CustomerId:req.CustomerID};q.Size=20;q.Limit=" limit 0,20";list,total:=freightService.GetReceiptService().SelectList(q);data:=make([]any,0,len(list));for _,v:=range list{data=append(data,map[string]any{"receiptNo":v.ReceiptNo,"amount":v.Amount,"currency":v.Currency,"allocatedAmount":v.AllocatedAmount,"remaining":v.Amount-v.AllocatedAmount,"status":v.Status,"receiptTime":v.ReceiptTime})}
	r:=protocol.NewAgentResultV2("我的收款核销",fmt.Sprintf("查询到 %d 条收款记录。",*total),[]protocol.BlockItem{{Type:"table",Columns:[]protocol.TableColumn{{Label:"收款单号",Field:"receiptNo"},{Label:"金额",Field:"amount"},{Label:"币种",Field:"currency"},{Label:"已核销",Field:"allocatedAmount"},{Label:"剩余",Field:"remaining"},{Label:"状态",Field:"status"},{Label:"收款时间",Field:"receiptTime"}},Data:data}});return &r
}

func buildPaymentDeclarationDraft(message string)*protocol.AgentResult{
	amount:=0.0;if m:=explicitMoneyPattern.FindStringSubmatch(message);len(m)>1{amount,_=strconv.ParseFloat(m[1],64);if m[2]=="万"{amount*=10000}};shipmentNo:="";if shipmentNoPattern.MatchString(message){shipmentNo=shipmentNoPattern.FindString(message)}
	r:=protocol.NewAgentResultV2("提交付款申报","付款申报提交后由后台审核，提交本身不代表款项已到账。",[]protocol.BlockItem{{Type:"form",Title:"付款申报",FormCode:"customer_payment_declaration",SubmitAPI:"/api/agent/form/submit",Fields:[]protocol.FormField{{Field:"amount",Label:"付款金额",Component:"number",Required:true},{Field:"currency",Label:"币种",Component:"input",Required:true},{Field:"paymentTime",Label:"付款日期",Component:"date",Required:true},{Field:"shipmentNo",Label:"关联出货计划",Component:"input",Required:true},{Field:"voucher",Label:"付款凭证",Component:"upload",Required:true},{Field:"remark",Label:"备注",Component:"textarea"}},InitialValues:map[string]any{"amount":amount,"currency":"CNY","shipmentNo":shipmentNo}}});return &r
}
