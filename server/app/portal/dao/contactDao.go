package dao

import (
	"baize/app/common/datasource"
	"baize/app/constant/constants"
	"baize/app/portal/models"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var contactDaoImpl *contactDao

func init() {
	contactDaoImpl = &contactDao{
		selectContactSql: `select contact_id, lead_no, contact_name, company_name, phone, email, route,
			cargo_info, message, source, status, ip_addr, user_agent, remark,
			create_by, create_time, update_by, update_time `,
		fromContactSql: ` from portal_contact`,
	}
}

type contactDao struct {
	selectContactSql string
	fromContactSql   string
}

func GetContactDao() *contactDao {
	return contactDaoImpl
}

func (contactDao *contactDao) InsertContact(contact *models.ContactDML) {
	insertSQL := `insert into portal_contact(
		contact_id, lead_no, contact_name, company_name, phone, email, route,
		cargo_info, message, source, status, ip_addr, user_agent,
		create_by, create_time, update_by, update_time
	) values (
		:contact_id, :lead_no, :contact_name, :company_name, :phone, :email, :route,
		:cargo_info, :message, :source, :status, :ip_addr, :user_agent,
		:create_by, now(), :update_by, now()
	)`
	_, err := datasource.GetMasterDb().NamedExec(insertSQL, contact)
	if err != nil {
		panic(err)
	}
}

func (contactDao *contactDao) SelectContactList(contact *models.ContactDQL) (list []*models.ContactVo, total *int64) {
	whereSql := ``
	if contact.ContactName != "" {
		whereSql += " AND contact_name like concat('%', :contact_name, '%')"
	}
	if contact.CompanyName != "" {
		whereSql += " AND company_name like concat('%', :company_name, '%')"
	}
	if contact.Phone != "" {
		whereSql += " AND phone like concat('%', :phone, '%')"
	}
	if contact.Email != "" {
		whereSql += " AND email like concat('%', :email, '%')"
	}
	if contact.Route != "" {
		whereSql += " AND route like concat('%', :route, '%')"
	}
	if contact.Source != "" {
		whereSql += " AND source = :source"
	}
	if contact.Status != "" {
		whereSql += " AND status = :status"
	}
	if contact.BeginTime != "" {
		whereSql += " AND date_format(create_time,'%Y-%m-%d') >= :begin_time"
	}
	if contact.EndTime != "" {
		whereSql += " AND date_format(create_time,'%Y-%m-%d') <= :end_time"
	}
	if whereSql != "" {
		whereSql = " where " + whereSql[4:]
	}

	countRow, err := datasource.GetMasterDb().NamedQuery(constants.MysqlCount+contactDao.fromContactSql+whereSql, contact)
	if err != nil {
		panic(err)
	}
	total = new(int64)
	if countRow.Next() {
		countRow.Scan(total)
	}
	defer countRow.Close()

	list = make([]*models.ContactVo, 0, contact.Size)
	if *total > contact.Offset {
		whereSql += " order by create_time desc"
		if contact.Limit != "" {
			whereSql += contact.Limit
		}
		listRows, err := datasource.GetMasterDb().NamedQuery(contactDao.selectContactSql+contactDao.fromContactSql+whereSql, contact)
		if err != nil {
			panic(err)
		}
		for listRows.Next() {
			contactVo := new(models.ContactVo)
			if err := listRows.StructScan(contactVo); err != nil {
				panic(err)
			}
			list = append(list, contactVo)
		}
		defer listRows.Close()
	}
	return
}

func (contactDao *contactDao) SelectContactById(contactId int64) (contact *models.ContactVo) {
	contact = new(models.ContactVo)
	err := datasource.GetMasterDb().Get(contact, contactDao.selectContactSql+contactDao.fromContactSql+" where contact_id = ?", contactId)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		panic(err)
	}
	return
}

func (contactDao *contactDao) UpdateContact(contact *models.ContactDML) {
	updateSQL := `update portal_contact set update_by = :update_by, update_time = now()`
	if contact.Status != "" {
		updateSQL += ", status = :status"
	}
	updateSQL += ", remark = :remark where contact_id = :contact_id"
	_, err := datasource.GetMasterDb().NamedExec(updateSQL, contact)
	if err != nil {
		panic(err)
	}
}

func (contactDao *contactDao) DeleteContactByIds(contactIds []int64) {
	query, args, err := sqlx.In("delete from portal_contact where contact_id in (?)", contactIds)
	if err != nil {
		panic(err)
	}
	_, err = datasource.GetMasterDb().Exec(query, args...)
	if err != nil {
		panic(err)
	}
}
