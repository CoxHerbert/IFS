package dao

import (
	"baize/app/common/datasource"
	"baize/app/constant/constants"
	"baize/app/notification/models"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var notificationDaoImpl *notificationDao

func init() {
	notificationDaoImpl = &notificationDao{
		selectSql: `select notification_id, user_id, title, content, biz_type, biz_id, read_flag, create_by, create_time, read_time `,
		fromSql:   ` from sys_notification`,
	}
}

type notificationDao struct {
	selectSql string
	fromSql   string
}

func GetNotificationDao() *notificationDao {
	return notificationDaoImpl
}

func (dao *notificationDao) Insert(item *models.NotificationDML) {
	_, err := datasource.GetMasterDb().NamedExec(`insert into sys_notification(
		notification_id, user_id, title, content, biz_type, biz_id, read_flag, create_by, create_time
	) values (
		:notification_id, :user_id, :title, :content, :biz_type, :biz_id, :read_flag, :create_by, now()
	)`, item)
	if err != nil {
		panic(err)
	}
}

func (dao *notificationDao) SelectList(userId int64, query *models.NotificationDQL) (list []*models.NotificationVo, total *int64) {
	whereSql := " where user_id = :user_id"
	params := map[string]any{
		"title":     query.Title,
		"biz_type":  query.BizType,
		"user_id":   userId,
		"read_flag": query.ReadFlag,
	}
	if query.Title != "" {
		whereSql += " and title like concat('%', :title, '%')"
	}
	if query.BizType != "" {
		whereSql += " and biz_type = :biz_type"
	}
	if query.ReadFlag != "" {
		whereSql += " and read_flag = :read_flag"
	}

	countRow, err := datasource.GetMasterDb().NamedQuery(constants.MysqlCount+dao.fromSql+whereSql, params)
	if err != nil {
		panic(err)
	}
	total = new(int64)
	if countRow.Next() {
		countRow.Scan(total)
	}
	defer countRow.Close()

	list = make([]*models.NotificationVo, 0, query.Size)
	if *total > query.Offset {
		sqlText := dao.selectSql + dao.fromSql + whereSql + " order by read_flag asc, create_time desc"
		if query.Limit != "" {
			sqlText += query.Limit
		}
		rows, err := datasource.GetMasterDb().NamedQuery(sqlText, params)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		for rows.Next() {
			vo := new(models.NotificationVo)
			if err := rows.StructScan(vo); err != nil {
				panic(err)
			}
			list = append(list, vo)
		}
	}
	return
}

func (dao *notificationDao) CountUnread(userId int64) int64 {
	var total int64
	err := datasource.GetMasterDb().Get(&total, `select count(1) from sys_notification where user_id = ? and read_flag = '0'`, userId)
	if err == sql.ErrNoRows {
		return 0
	}
	if err != nil {
		panic(err)
	}
	return total
}

func (dao *notificationDao) MarkRead(notificationId int64, userId int64) bool {
	result, err := datasource.GetMasterDb().Exec(
		`update sys_notification set read_flag = '1', read_time = now() where notification_id = ? and user_id = ? and read_flag = '0'`,
		notificationId,
		userId,
	)
	if err != nil {
		panic(err)
	}
	affected, _ := result.RowsAffected()
	return affected > 0
}

func (dao *notificationDao) MarkAllRead(userId int64) int64 {
	result, err := datasource.GetMasterDb().Exec(
		`update sys_notification set read_flag = '1', read_time = now() where user_id = ? and read_flag = '0'`,
		userId,
	)
	if err != nil {
		panic(err)
	}
	affected, _ := result.RowsAffected()
	return affected
}

func (dao *notificationDao) DeleteByIds(notificationIds []int64, userId int64) int64 {
	if len(notificationIds) == 0 {
		return 0
	}
	query, args, err := sqlx.In("delete from sys_notification where user_id = ? and notification_id in (?)", userId, notificationIds)
	if err != nil {
		panic(err)
	}
	query = datasource.GetMasterDb().Rebind(query)
	result, err := datasource.GetMasterDb().Exec(query, args...)
	if err != nil {
		panic(err)
	}
	affected, _ := result.RowsAffected()
	return affected
}
