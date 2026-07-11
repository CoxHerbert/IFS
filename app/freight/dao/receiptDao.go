package dao

import (
	"baize/app/common/datasource"
	"baize/app/constant/constants"
	"baize/app/freight/models"
	"database/sql"
	"fmt"
)

type receiptDao struct{}
var receiptDaoImpl = &receiptDao{}
func GetReceiptDao() *receiptDao { return receiptDaoImpl }

func (dao *receiptDao) SelectList(q *models.ReceiptDQL) ([]*models.ReceiptVo, *int64) {
	where := " where 1=1"
	if q.ReceiptNo != "" { where += " and r.receipt_no like concat('%', :receipt_no, '%')" }
	if q.CustomerName != "" { where += " and r.customer_name like concat('%', :customer_name, '%')" }
	if q.Status != "" { where += " and r.status = :status" }
	if q.CustomerId != 0 { where += " and r.customer_id = :customer_id" }
	if q.SalesUserId != 0 { where += " and c.sales_user_id = :sales_user_id" }
	from := ` from freight_receipt r left join customer c on c.customer_id=r.customer_id` + where
	total := new(int64)
	rows, err := datasource.GetMasterDb().NamedQuery(constants.MysqlCount + from, q); if err != nil { panic(err) }
	if rows.Next() { rows.Scan(total) }; rows.Close()
	list := make([]*models.ReceiptVo, 0, q.Size)
	if *total > q.Offset {
		sqlText := `select r.receipt_id,r.receipt_no,r.customer_id,r.customer_name,r.amount,r.currency,r.receipt_time,r.payment_method,
			r.status,r.voucher_url,r.voucher_name,r.remark,r.create_by,r.create_time,
			coalesce((select sum(a.allocated_amount) from freight_receipt_allocation a where a.receipt_id=r.receipt_id),0) allocated_amount` + from + " order by r.create_time desc" + q.Limit
		r, err := datasource.GetMasterDb().NamedQuery(sqlText, q); if err != nil { panic(err) }
		for r.Next() { vo := new(models.ReceiptVo); if err := r.StructScan(vo); err != nil { panic(err) }; list = append(list, vo) }; r.Close()
	}
	return list, total
}

func (dao *receiptDao) Allocate(receiptId int64, allocation *models.ReceiptAllocationDML) error {
	tx:=datasource.GetMasterDb().MustBegin(); var total,allocated float64
	if err:=tx.Get(&total,`select amount from freight_receipt where receipt_id=? for update`,receiptId);err!=nil{tx.Rollback();return fmt.Errorf("收款单不存在")}
	if err:=tx.Get(&allocated,`select coalesce(sum(allocated_amount),0) from freight_receipt_allocation where receipt_id=?`,receiptId);err!=nil{tx.Rollback();return err}
	if allocated+allocation.AllocatedAmount-total>0.001{tx.Rollback();return fmt.Errorf("核销金额超过收款单剩余可核销金额")}
	if _,err:=tx.NamedExec(`insert into freight_receipt_allocation(allocation_id,receipt_id,shipment_id,allocated_amount) values(:allocation_id,:receipt_id,:shipment_id,:allocated_amount)`,allocation);err!=nil{tx.Rollback();return fmt.Errorf("该收款单已核销过此出货计划")}
	status:="PARTIAL";if allocated+allocation.AllocatedAmount>=total-0.001{status="ALLOCATED"}
	if _,err:=tx.Exec(`update freight_receipt set status=?,update_time=now() where receipt_id=?`,status,receiptId);err!=nil{tx.Rollback();return err}
	if _,err:=tx.Exec(`update freight_shipment_plan p set p.payment_amount=(select coalesce(sum(x.allocated_amount),0) from freight_receipt_allocation x where x.shipment_id=p.shipment_id),p.payment_status='PARTIAL',p.update_time=now() where p.shipment_id=?`,allocation.ShipmentId);err!=nil{tx.Rollback();return err}
	return tx.Commit()
}

func (dao *receiptDao) SelectById(id int64) *models.ReceiptVo {
	vo := new(models.ReceiptVo)
	err := datasource.GetMasterDb().Get(vo, `select r.receipt_id,r.receipt_no,r.customer_id,r.customer_name,r.amount,r.currency,r.receipt_time,r.payment_method,
		r.status,r.voucher_url,r.voucher_name,r.remark,r.create_by,r.create_time,
		coalesce((select sum(a.allocated_amount) from freight_receipt_allocation a where a.receipt_id=r.receipt_id),0) allocated_amount
		from freight_receipt r where r.receipt_id=?`, id)
	if err == sql.ErrNoRows { return nil }; if err != nil { panic(err) }
	vo.Allocations = dao.SelectAllocations(id)
	return vo
}

func (dao *receiptDao) SelectAllocations(id int64) []*models.ReceiptAllocationVo {
	list := make([]*models.ReceiptAllocationVo, 0)
	err := datasource.GetMasterDb().Select(&list, `select a.allocation_id,a.receipt_id,a.shipment_id,p.shipment_no,a.allocated_amount
		from freight_receipt_allocation a join freight_shipment_plan p on p.shipment_id=a.shipment_id where a.receipt_id=? order by a.allocation_id`, id)
	if err != nil { panic(err) }; return list
}

func (dao *receiptDao) Insert(receipt *models.ReceiptDML, allocations []*models.ReceiptAllocationDML) {
	tx := datasource.GetMasterDb().MustBegin()
	_, err := tx.NamedExec(`insert into freight_receipt(receipt_id,receipt_no,customer_id,customer_name,amount,currency,receipt_time,payment_method,status,voucher_url,voucher_name,remark,create_by,create_time,update_by,update_time)
		values(:receipt_id,:receipt_no,:customer_id,:customer_name,:amount,:currency,:receipt_time,:payment_method,:status,:voucher_url,:voucher_name,:remark,:create_by,now(),:create_by,now())`, receipt)
	if err != nil { tx.Rollback(); panic(err) }
	for _, a := range allocations {
		if _, err = tx.NamedExec(`insert into freight_receipt_allocation(allocation_id,receipt_id,shipment_id,allocated_amount) values(:allocation_id,:receipt_id,:shipment_id,:allocated_amount)`, a); err != nil { tx.Rollback(); panic(err) }
		if _, err = tx.Exec(`update freight_shipment_plan p set p.payment_amount=(select coalesce(sum(x.allocated_amount),0) from freight_receipt_allocation x where x.shipment_id=p.shipment_id), p.payment_status='PARTIAL', p.update_time=now() where p.shipment_id=?`, a.ShipmentId); err != nil { tx.Rollback(); panic(err) }
	}
	if err = tx.Commit(); err != nil { panic(err) }
}

func (dao *receiptDao) Delete(id int64) {
	tx := datasource.GetMasterDb().MustBegin()
	shipmentIds := make([]int64, 0)
	if err := tx.Select(&shipmentIds, `select shipment_id from freight_receipt_allocation where receipt_id=?`, id); err != nil { tx.Rollback(); panic(err) }
	if _, err := tx.Exec(`delete from freight_receipt where receipt_id=?`, id); err != nil { tx.Rollback(); panic(err) }
	for _, shipmentId := range shipmentIds {
		if _, err := tx.Exec(`update freight_shipment_plan p set p.payment_amount=(select coalesce(sum(x.allocated_amount),0) from freight_receipt_allocation x where x.shipment_id=p.shipment_id), p.payment_status=case when (select coalesce(sum(x.allocated_amount),0) from freight_receipt_allocation x where x.shipment_id=p.shipment_id)=0 then 'UNPAID' else 'PARTIAL' end, p.update_time=now() where p.shipment_id=?`, shipmentId); err != nil { tx.Rollback(); panic(err) }
	}
	if err := tx.Commit(); err != nil { panic(err) }
}

func (dao *receiptDao) InsertDeclaration(item *models.PaymentDeclarationDML) {
	_,err:=datasource.GetMasterDb().NamedExec(`insert into freight_payment_declaration(declaration_id,declaration_no,customer_id,customer_name,shipment_id,shipment_no,amount,currency,payment_time,voucher_url,voucher_name,status,remark,create_by,create_time,update_by,update_time)
		values(:declaration_id,:declaration_no,:customer_id,:customer_name,:shipment_id,:shipment_no,:amount,:currency,:payment_time,:voucher_url,:voucher_name,:status,:remark,:create_by,now(),:create_by,now())`,item);if err!=nil{panic(err)}
}

func (dao *receiptDao) SelectDeclarationList(q *models.PaymentDeclarationDQL)([]*models.PaymentDeclarationVo,*int64){
	where:=" where 1=1";if q.DeclarationNo!=""{where+=" and d.declaration_no like concat('%',:declaration_no,'%')"};if q.CustomerName!=""{where+=" and d.customer_name like concat('%',:customer_name,'%')"};if q.Status!=""{where+=" and d.status=:status"};if q.SalesUserId!=0{where+=" and c.sales_user_id=:sales_user_id"};from:=" from freight_payment_declaration d left join customer c on c.customer_id=d.customer_id"+where
	total:=new(int64);rows,err:=datasource.GetMasterDb().NamedQuery(constants.MysqlCount+from,q);if err!=nil{panic(err)};if rows.Next(){rows.Scan(total)};rows.Close();list:=make([]*models.PaymentDeclarationVo,0,q.Size)
	if *total>q.Offset{r,err:=datasource.GetMasterDb().NamedQuery(`select d.declaration_id,d.declaration_no,d.customer_id,d.customer_name,d.shipment_id,d.shipment_no,d.amount,d.currency,d.payment_time,d.voucher_url,d.voucher_name,d.status,d.remark,d.create_by,d.create_time,d.review_by,d.review_time,d.review_remark`+from+" order by d.create_time desc"+q.Limit,q);if err!=nil{panic(err)};for r.Next(){vo:=new(models.PaymentDeclarationVo);if err:=r.StructScan(vo);err!=nil{panic(err)};list=append(list,vo)};r.Close()};return list,total
}

func (dao *receiptDao) SelectDeclarationById(id int64)*models.PaymentDeclarationVo{vo:=new(models.PaymentDeclarationVo);err:=datasource.GetMasterDb().Get(vo,`select declaration_id,declaration_no,customer_id,customer_name,shipment_id,shipment_no,amount,currency,payment_time,voucher_url,voucher_name,status,remark,create_by,create_time,review_by,review_time,review_remark from freight_payment_declaration where declaration_id=?`,id);if err==sql.ErrNoRows{return nil};if err!=nil{panic(err)};return vo}

func (dao *receiptDao) ApproveDeclaration(declarationId int64,receipt *models.ReceiptDML,allocation *models.ReceiptAllocationDML,reviewBy,remark string)error{
	tx:=datasource.GetMasterDb().MustBegin();result,err:=tx.Exec(`update freight_payment_declaration set status='CONFIRMED',review_by=?,review_time=now(),review_remark=?,update_by=?,update_time=now() where declaration_id=? and status='PENDING'`,reviewBy,remark,reviewBy,declarationId);if err!=nil{tx.Rollback();return err};affected,_:=result.RowsAffected();if affected==0{tx.Rollback();return fmt.Errorf("申报已处理，请勿重复审核")}
	if _,err=tx.NamedExec(`insert into freight_receipt(receipt_id,receipt_no,customer_id,customer_name,amount,currency,receipt_time,payment_method,status,voucher_url,voucher_name,remark,create_by,create_time,update_by,update_time) values(:receipt_id,:receipt_no,:customer_id,:customer_name,:amount,:currency,:receipt_time,:payment_method,:status,:voucher_url,:voucher_name,:remark,:create_by,now(),:create_by,now())`,receipt);err!=nil{tx.Rollback();return err}
	if _,err=tx.NamedExec(`insert into freight_receipt_allocation(allocation_id,receipt_id,shipment_id,allocated_amount) values(:allocation_id,:receipt_id,:shipment_id,:allocated_amount)`,allocation);err!=nil{tx.Rollback();return err}
	if _,err=tx.Exec(`update freight_shipment_plan p set p.payment_amount=(select coalesce(sum(x.allocated_amount),0) from freight_receipt_allocation x where x.shipment_id=p.shipment_id),p.payment_status='PARTIAL',p.update_by=?,p.update_time=now() where p.shipment_id=?`,reviewBy,allocation.ShipmentId);err!=nil{tx.Rollback();return err};return tx.Commit()
}

func (dao *receiptDao) RejectDeclaration(id int64,reviewBy,remark string)error{result,err:=datasource.GetMasterDb().Exec(`update freight_payment_declaration set status='REJECTED',review_by=?,review_time=now(),review_remark=?,update_by=?,update_time=now() where declaration_id=? and status='PENDING'`,reviewBy,remark,reviewBy,id);if err!=nil{return err};affected,_:=result.RowsAffected();if affected==0{return fmt.Errorf("申报已处理，请勿重复审核")};return nil}
