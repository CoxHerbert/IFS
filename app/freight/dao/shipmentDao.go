package dao

import (
	"baize/app/common/datasource"
	"baize/app/constant/constants"
	"baize/app/freight/models"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var shipmentDaoImpl *shipmentDao

func init() {
	shipmentDaoImpl = &shipmentDao{
		selectPlanSql: `select shipment_id, shipment_no, order_no, customer_id, customer_name, sales_user_id, sales_user_name, pol, pod,
			planned_etd, planned_eta, actual_etd, actual_eta, status, payment_status, payment_amount, total_weight, total_volume, total_cartons,
			share_token, remark, create_by, create_time, update_by, update_time `,
		fromPlanSql: ` from freight_shipment_plan`,
	}
}

type shipmentDao struct {
	selectPlanSql string
	fromPlanSql   string
}

func GetShipmentDao() *shipmentDao {
	return shipmentDaoImpl
}

func (dao *shipmentDao) InsertShipment(plan *models.ShipmentPlanDML, cargoList []*models.CargoDML, containers []*models.ContainerPlanDML) {
	tx := datasource.GetMasterDb().MustBegin()
	_, err := tx.NamedExec(`insert into freight_shipment_plan(
		shipment_id, shipment_no, order_no, customer_id, customer_name, sales_user_id, sales_user_name, pol, pod, planned_etd, planned_eta,
		status, payment_status, payment_amount, total_weight, total_volume, total_cartons, share_token, remark, create_by, create_time, update_by, update_time
	) values (
		:shipment_id, :shipment_no, :order_no, :customer_id, :customer_name, :sales_user_id, :sales_user_name, :pol, :pod, :planned_etd, :planned_eta,
		:status, :payment_status, :payment_amount, :total_weight, :total_volume, :total_cartons, :share_token, :remark, :create_by, now(), :update_by, now()
	)`, plan)
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	for _, cargo := range cargoList {
		_, err = tx.NamedExec(`insert into freight_shipment_cargo(
			cargo_id, shipment_id, sku, cargo_name, package_type, quantity, cartons, weight_kg, volume_cbm, length_cm, width_cm, height_cm
		) values (
			:cargo_id, :shipment_id, :sku, :cargo_name, :package_type, :quantity, :cartons, :weight_kg, :volume_cbm, :length_cm, :width_cm, :height_cm
		)`, cargo)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}
	for _, container := range containers {
		_, err = tx.NamedExec(`insert into freight_container_plan(
			container_plan_id, shipment_id, container_type, quantity, max_volume, max_weight, used_volume, used_weight, load_rate, remark
		) values (
			:container_plan_id, :shipment_id, :container_type, :quantity, :max_volume, :max_weight, :used_volume, :used_weight, :load_rate, :remark
		)`, container)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}
	if err = tx.Commit(); err != nil {
		panic(err)
	}
}

func (dao *shipmentDao) SelectShipmentList(query *models.ShipmentPlanDQL) (list []*models.ShipmentPlanVo, total *int64) {
	whereSql := ""
	if query.ShipmentNo != "" {
		whereSql += " AND shipment_no like concat('%', :shipment_no, '%')"
	}
	if query.OrderNo != "" {
		whereSql += " AND order_no like concat('%', :order_no, '%')"
	}
	if query.CustomerId != 0 {
		whereSql += " AND customer_id = :customer_id"
	}
	if query.CustomerName != "" {
		whereSql += " AND customer_name like concat('%', :customer_name, '%')"
	}
	if query.SalesUserId != 0 {
		whereSql += " AND sales_user_id = :sales_user_id"
	}
	if query.Pol != "" {
		whereSql += " AND pol like concat('%', :pol, '%')"
	}
	if query.Pod != "" {
		whereSql += " AND pod like concat('%', :pod, '%')"
	}
	if query.Status != "" {
		whereSql += " AND status = :status"
	}
	if query.BeginTime != "" {
		whereSql += " AND date_format(create_time,'%Y-%m-%d') >= :begin_time"
	}
	if query.EndTime != "" {
		whereSql += " AND date_format(create_time,'%Y-%m-%d') <= :end_time"
	}
	if whereSql != "" {
		whereSql = " where " + whereSql[5:]
	}

	countRow, err := datasource.GetMasterDb().NamedQuery(constants.MysqlCount+dao.fromPlanSql+whereSql, query)
	if err != nil {
		panic(err)
	}
	total = new(int64)
	if countRow.Next() {
		countRow.Scan(total)
	}
	defer countRow.Close()

	list = make([]*models.ShipmentPlanVo, 0, query.Size)
	if *total > query.Offset {
		whereSql += " order by create_time desc"
		if query.Limit != "" {
			whereSql += query.Limit
		}
		rows, err := datasource.GetMasterDb().NamedQuery(dao.selectPlanSql+dao.fromPlanSql+whereSql, query)
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			vo := new(models.ShipmentPlanVo)
			if err := rows.StructScan(vo); err != nil {
				panic(err)
			}
			list = append(list, vo)
		}
		defer rows.Close()
	}
	return
}

func (dao *shipmentDao) SelectShipmentById(shipmentId int64) *models.ShipmentPlanVo {
	plan := new(models.ShipmentPlanVo)
	err := datasource.GetMasterDb().Get(plan, dao.selectPlanSql+dao.fromPlanSql+" where shipment_id = ?", shipmentId)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		panic(err)
	}
	return plan
}

func (dao *shipmentDao) SelectShipmentByToken(token string) *models.ShipmentPlanVo {
	plan := new(models.ShipmentPlanVo)
	err := datasource.GetMasterDb().Get(plan, dao.selectPlanSql+dao.fromPlanSql+" where share_token = ?", token)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		panic(err)
	}
	return plan
}

func (dao *shipmentDao) SelectCargoList(shipmentId int64) []*models.CargoVo {
	list := make([]*models.CargoVo, 0)
	err := datasource.GetMasterDb().Select(&list, `select cargo_id, shipment_id, sku, cargo_name, package_type, quantity,
		cartons, weight_kg, volume_cbm, length_cm, width_cm, height_cm
		from freight_shipment_cargo where shipment_id = ? order by cargo_id`, shipmentId)
	if err != nil {
		panic(err)
	}
	return list
}

func (dao *shipmentDao) SelectContainerList(shipmentId int64) []*models.ContainerPlanVo {
	list := make([]*models.ContainerPlanVo, 0)
	err := datasource.GetMasterDb().Select(&list, `select container_plan_id, shipment_id, container_type, quantity,
		max_volume, max_weight, used_volume, used_weight, load_rate, remark
		from freight_container_plan where shipment_id = ? order by container_plan_id`, shipmentId)
	if err != nil {
		panic(err)
	}
	return list
}

func (dao *shipmentDao) SelectOrderByShipmentId(shipmentId int64) *models.ShipmentOrderVo {
	order := new(models.ShipmentOrderVo)
	err := datasource.GetMasterDb().Get(order, `select order_id, shipment_id, order_no, status, create_time
		from freight_shipment_order where shipment_id = ?`, shipmentId)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		panic(err)
	}
	return order
}

func (dao *shipmentDao) SelectPaymentList(shipmentId int64) []*models.ShipmentPaymentVo {
	list := make([]*models.ShipmentPaymentVo, 0)
	err := datasource.GetMasterDb().Select(&list, `select r.receipt_id payment_id, a.shipment_id, a.allocated_amount amount, r.currency,
		r.receipt_time payment_time, r.payment_method, r.voucher_url, r.voucher_name, r.remark, r.create_by, r.create_time
		from freight_receipt_allocation a join freight_receipt r on r.receipt_id=a.receipt_id
		where a.shipment_id = ? order by r.receipt_time desc, r.create_time desc`, shipmentId)
	if err != nil { panic(err) }
	return list
}

func (dao *shipmentDao) InsertPayment(payment *models.ShipmentPaymentDML) {
	tx := datasource.GetMasterDb().MustBegin()
	if _, err := tx.NamedExec(`insert into freight_shipment_payment(payment_id, shipment_id, amount, currency, payment_time,
		payment_method, voucher_url, voucher_name, remark, create_by, create_time)
		values(:payment_id, :shipment_id, :amount, :currency, :payment_time, :payment_method, :voucher_url, :voucher_name, :remark, :create_by, now())`, payment); err != nil {
		tx.Rollback(); panic(err)
	}
	if _, err := tx.Exec(`update freight_shipment_plan p set p.payment_amount =
		(select coalesce(sum(amount), 0) from freight_shipment_payment where shipment_id = p.shipment_id),
		p.payment_status = 'PARTIAL', p.update_by = ?, p.update_time = now() where p.shipment_id = ?`, payment.CreateBy, payment.ShipmentId); err != nil {
		tx.Rollback(); panic(err)
	}
	if err := tx.Commit(); err != nil { panic(err) }
}

func (dao *shipmentDao) DeletePayment(paymentId, shipmentId int64, updateBy string) bool {
	tx := datasource.GetMasterDb().MustBegin()
	result, err := tx.Exec(`delete from freight_shipment_payment where payment_id = ? and shipment_id = ?`, paymentId, shipmentId)
	if err != nil { tx.Rollback(); panic(err) }
	affected, _ := result.RowsAffected()
	if affected == 0 { tx.Rollback(); return false }
	_, err = tx.Exec(`update freight_shipment_plan p set p.payment_amount =
		(select coalesce(sum(amount), 0) from freight_shipment_payment where shipment_id = p.shipment_id),
		p.payment_status = case when (select coalesce(sum(amount), 0) from freight_shipment_payment where shipment_id = p.shipment_id) = 0 then 'UNPAID' else 'PARTIAL' end,
		p.update_by = ?, p.update_time = now() where p.shipment_id = ?`, updateBy, shipmentId)
	if err != nil { tx.Rollback(); panic(err) }
	if err = tx.Commit(); err != nil { panic(err) }
	return true
}

func (dao *shipmentDao) UpdateShipmentStatus(update *models.ShipmentStatusUpdateDML) {
	updateSQL := `update freight_shipment_plan set status = :status, payment_status = :payment_status, payment_amount = :payment_amount, update_by = :update_by, update_time = now()`
	if update.ActualEtd != "" {
		updateSQL += ", actual_etd = :actual_etd"
	}
	if update.ActualEta != "" {
		updateSQL += ", actual_eta = :actual_eta"
	}
	if update.Remark != "" {
		updateSQL += ", remark = :remark"
	}
	updateSQL += " where shipment_id = :shipment_id"
	_, err := datasource.GetMasterDb().NamedExec(updateSQL, update)
	if err != nil {
		panic(err)
	}
}

func (dao *shipmentDao) UpdateShipmentCustomer(shipmentId int64, customerId int64, customerName string, salesUserId int64, salesUserName string, updateBy string) {
	_, err := datasource.GetMasterDb().Exec(
		`update freight_shipment_plan set customer_id = ?, customer_name = ?, sales_user_id = ?, sales_user_name = ?, update_by = ?, update_time = now() where shipment_id = ?`,
		customerId,
		customerName,
		salesUserId,
		salesUserName,
		updateBy,
		shipmentId,
	)
	if err != nil {
		panic(err)
	}
}

func (dao *shipmentDao) InsertShipmentOrder(order *models.ShipmentOrderDML) {
	_, err := datasource.GetMasterDb().NamedExec(`insert into freight_shipment_order(
		order_id, shipment_id, order_no, status, create_by, create_time, update_by, update_time
	) values (
		:order_id, :shipment_id, :order_no, :status, :create_by, now(), :update_by, now()
	)`, order)
	if err != nil {
		panic(err)
	}
}

func (dao *shipmentDao) DeleteShipmentByIds(shipmentIds []int64) {
	query, args, err := sqlx.In("delete from freight_shipment_plan where shipment_id in (?)", shipmentIds)
	if err != nil {
		panic(err)
	}
	_, err = datasource.GetMasterDb().Exec(query, args...)
	if err != nil {
		panic(err)
	}
}
