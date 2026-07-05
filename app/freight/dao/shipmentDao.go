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
			planned_etd, planned_eta, actual_etd, actual_eta, status, total_weight, total_volume, total_cartons,
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
		status, total_weight, total_volume, total_cartons, share_token, remark, create_by, create_time, update_by, update_time
	) values (
		:shipment_id, :shipment_no, :order_no, :customer_id, :customer_name, :sales_user_id, :sales_user_name, :pol, :pod, :planned_etd, :planned_eta,
		:status, :total_weight, :total_volume, :total_cartons, :share_token, :remark, :create_by, now(), :update_by, now()
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

func (dao *shipmentDao) UpdateShipmentStatus(update *models.ShipmentStatusUpdateDML) {
	updateSQL := `update freight_shipment_plan set status = :status, update_by = :update_by, update_time = now()`
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
