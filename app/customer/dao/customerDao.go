package dao

import (
	"baize/app/common/datasource"
	"baize/app/constant/constants"
	"baize/app/customer/models"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var customerDaoImpl *customerDao

func init() {
	customerDaoImpl = &customerDao{
		selectCustomerSql: `select customer_id, customer_no, customer_name, company_name, contact_name, phone, email,
			status, remark, create_by, create_time, update_by, update_time `,
		fromCustomerSql: ` from customer`,
		selectContactSql: `select contact_id, customer_id, contact_name, position, phone, email, wechat,
			is_primary, status, remark, create_by, create_time, update_by, update_time `,
		fromContactSql: ` from customer_contact`,
		selectAccountSql: `select a.account_id, a.customer_id, c.customer_no, c.customer_name, c.company_name,
			a.username, a.real_name, a.phone, a.email, a.is_main, a.status, a.last_login_time,
			a.remark, a.create_by, a.create_time, a.update_by, a.update_time `,
		fromAccountSql: ` from customer_account a left join customer c on a.customer_id = c.customer_id`,
	}
}

type customerDao struct {
	selectCustomerSql string
	fromCustomerSql   string
	selectContactSql  string
	fromContactSql    string
	selectAccountSql  string
	fromAccountSql    string
}

func GetCustomerDao() *customerDao {
	return customerDaoImpl
}

func (dao *customerDao) SelectCustomerList(customer *models.CustomerDQL) (list []*models.CustomerVo, total *int64) {
	whereSql := ``
	if customer.CustomerName != "" {
		whereSql += " AND customer_name like concat('%', :customer_name, '%')"
	}
	if customer.CompanyName != "" {
		whereSql += " AND company_name like concat('%', :company_name, '%')"
	}
	if customer.ContactName != "" {
		whereSql += " AND contact_name like concat('%', :contact_name, '%')"
	}
	if customer.Phone != "" {
		whereSql += " AND phone like concat('%', :phone, '%')"
	}
	if customer.Email != "" {
		whereSql += " AND email like concat('%', :email, '%')"
	}
	if customer.Status != "" {
		whereSql += " AND status = :status"
	}
	if customer.BeginTime != "" {
		whereSql += " AND date_format(create_time,'%Y-%m-%d') >= :begin_time"
	}
	if customer.EndTime != "" {
		whereSql += " AND date_format(create_time,'%Y-%m-%d') <= :end_time"
	}
	if whereSql != "" {
		whereSql = " where " + whereSql[4:]
	}

	countRow, err := datasource.GetMasterDb().NamedQuery(constants.MysqlCount+dao.fromCustomerSql+whereSql, customer)
	if err != nil {
		panic(err)
	}
	total = new(int64)
	if countRow.Next() {
		countRow.Scan(total)
	}
	defer countRow.Close()

	list = make([]*models.CustomerVo, 0, customer.Size)
	if *total > customer.Offset {
		whereSql += " order by create_time desc"
		if customer.Limit != "" {
			whereSql += customer.Limit
		}
		rows, err := datasource.GetMasterDb().NamedQuery(dao.selectCustomerSql+dao.fromCustomerSql+whereSql, customer)
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			item := new(models.CustomerVo)
			if err := rows.StructScan(item); err != nil {
				panic(err)
			}
			list = append(list, item)
		}
		defer rows.Close()
	}
	return
}

func (dao *customerDao) SelectCustomerById(customerId int64) *models.CustomerVo {
	customer := new(models.CustomerVo)
	err := datasource.GetMasterDb().Get(customer, dao.selectCustomerSql+dao.fromCustomerSql+" where customer_id = ?", customerId)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		panic(err)
	}
	return customer
}

func (dao *customerDao) SelectCustomerOptions(keyword string) []*models.CustomerOptionVo {
	whereSql := " where status = '0'"
	args := make([]interface{}, 0, 2)
	if keyword != "" {
		whereSql += " and (customer_name like concat('%', ?, '%') or company_name like concat('%', ?, '%'))"
		args = append(args, keyword, keyword)
	}
	list := make([]*models.CustomerOptionVo, 0)
	err := datasource.GetMasterDb().Select(&list, `select customer_id, customer_no, customer_name, company_name from customer`+whereSql+` order by create_time desc limit 100`, args...)
	if err != nil {
		panic(err)
	}
	return list
}

func (dao *customerDao) InsertCustomer(customer *models.CustomerDML) {
	insertSQL := `insert into customer(
		customer_id, customer_no, customer_name, company_name, contact_name, phone, email,
		status, remark, create_by, create_time, update_by, update_time
	) values (
		:customer_id, :customer_no, :customer_name, :company_name, :contact_name, :phone, :email,
		:status, :remark, :create_by, now(), :update_by, now()
	)`
	if _, err := datasource.GetMasterDb().NamedExec(insertSQL, customer); err != nil {
		panic(err)
	}
}

func (dao *customerDao) UpdateCustomer(customer *models.CustomerDML) {
	updateSQL := `update customer set customer_name = :customer_name, company_name = :company_name,
		contact_name = :contact_name, phone = :phone, email = :email, status = :status,
		remark = :remark, update_by = :update_by, update_time = now()
		where customer_id = :customer_id`
	if _, err := datasource.GetMasterDb().NamedExec(updateSQL, customer); err != nil {
		panic(err)
	}
}

func (dao *customerDao) DeleteCustomerByIds(customerIds []int64) {
	query, args, err := sqlx.In("delete from customer where customer_id in (?)", customerIds)
	if err != nil {
		panic(err)
	}
	if _, err = datasource.GetMasterDb().Exec(query, args...); err != nil {
		panic(err)
	}
}

func (dao *customerDao) SelectContactList(contact *models.CustomerContactDQL) (list []*models.CustomerContactVo, total *int64) {
	whereSql := ``
	if contact.CustomerId != 0 {
		whereSql += " AND customer_id = :customer_id"
	}
	if contact.ContactName != "" {
		whereSql += " AND contact_name like concat('%', :contact_name, '%')"
	}
	if contact.Phone != "" {
		whereSql += " AND phone like concat('%', :phone, '%')"
	}
	if contact.Status != "" {
		whereSql += " AND status = :status"
	}
	if whereSql != "" {
		whereSql = " where " + whereSql[4:]
	}

	countRow, err := datasource.GetMasterDb().NamedQuery(constants.MysqlCount+dao.fromContactSql+whereSql, contact)
	if err != nil {
		panic(err)
	}
	total = new(int64)
	if countRow.Next() {
		countRow.Scan(total)
	}
	defer countRow.Close()

	list = make([]*models.CustomerContactVo, 0, contact.Size)
	if *total > contact.Offset {
		whereSql += " order by is_primary desc, create_time desc"
		if contact.Limit != "" {
			whereSql += contact.Limit
		}
		rows, err := datasource.GetMasterDb().NamedQuery(dao.selectContactSql+dao.fromContactSql+whereSql, contact)
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			item := new(models.CustomerContactVo)
			if err := rows.StructScan(item); err != nil {
				panic(err)
			}
			list = append(list, item)
		}
		defer rows.Close()
	}
	return
}

func (dao *customerDao) SelectContactById(contactId int64) *models.CustomerContactVo {
	contact := new(models.CustomerContactVo)
	err := datasource.GetMasterDb().Get(contact, dao.selectContactSql+dao.fromContactSql+" where contact_id = ?", contactId)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		panic(err)
	}
	return contact
}

func (dao *customerDao) InsertContact(contact *models.CustomerContactDML) {
	insertSQL := `insert into customer_contact(
		contact_id, customer_id, contact_name, position, phone, email, wechat, is_primary, status,
		remark, create_by, create_time, update_by, update_time
	) values (
		:contact_id, :customer_id, :contact_name, :position, :phone, :email, :wechat, :is_primary, :status,
		:remark, :create_by, now(), :update_by, now()
	)`
	if _, err := datasource.GetMasterDb().NamedExec(insertSQL, contact); err != nil {
		panic(err)
	}
}

func (dao *customerDao) UpdateContact(contact *models.CustomerContactDML) {
	updateSQL := `update customer_contact set contact_name = :contact_name, position = :position,
		phone = :phone, email = :email, wechat = :wechat, is_primary = :is_primary,
		status = :status, remark = :remark, update_by = :update_by, update_time = now()
		where contact_id = :contact_id`
	if _, err := datasource.GetMasterDb().NamedExec(updateSQL, contact); err != nil {
		panic(err)
	}
}

func (dao *customerDao) DeleteContactByIds(contactIds []int64) {
	query, args, err := sqlx.In("delete from customer_contact where contact_id in (?)", contactIds)
	if err != nil {
		panic(err)
	}
	if _, err = datasource.GetMasterDb().Exec(query, args...); err != nil {
		panic(err)
	}
}

func (dao *customerDao) SelectAccountList(account *models.CustomerAccountDQL) (list []*models.CustomerAccountVo, total *int64) {
	whereSql := ``
	if account.CustomerId != 0 {
		whereSql += " AND a.customer_id = :customer_id"
	}
	if account.Username != "" {
		whereSql += " AND a.username like concat('%', :username, '%')"
	}
	if account.RealName != "" {
		whereSql += " AND a.real_name like concat('%', :real_name, '%')"
	}
	if account.Phone != "" {
		whereSql += " AND a.phone like concat('%', :phone, '%')"
	}
	if account.Status != "" {
		whereSql += " AND a.status = :status"
	}
	if account.BeginTime != "" {
		whereSql += " AND date_format(a.create_time,'%Y-%m-%d') >= :begin_time"
	}
	if account.EndTime != "" {
		whereSql += " AND date_format(a.create_time,'%Y-%m-%d') <= :end_time"
	}
	if whereSql != "" {
		whereSql = " where " + whereSql[4:]
	}

	countRow, err := datasource.GetMasterDb().NamedQuery(constants.MysqlCount+dao.fromAccountSql+whereSql, account)
	if err != nil {
		panic(err)
	}
	total = new(int64)
	if countRow.Next() {
		countRow.Scan(total)
	}
	defer countRow.Close()

	list = make([]*models.CustomerAccountVo, 0, account.Size)
	if *total > account.Offset {
		whereSql += " order by a.create_time desc"
		if account.Limit != "" {
			whereSql += account.Limit
		}
		rows, err := datasource.GetMasterDb().NamedQuery(dao.selectAccountSql+dao.fromAccountSql+whereSql, account)
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			item := new(models.CustomerAccountVo)
			if err := rows.StructScan(item); err != nil {
				panic(err)
			}
			list = append(list, item)
		}
		defer rows.Close()
	}
	return
}

func (dao *customerDao) SelectAccountById(accountId int64) *models.CustomerAccountVo {
	account := new(models.CustomerAccountVo)
	err := datasource.GetMasterDb().Get(account, dao.selectAccountSql+dao.fromAccountSql+" where a.account_id = ?", accountId)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		panic(err)
	}
	return account
}

func (dao *customerDao) SelectAccountByUsername(username string) (account *models.CustomerAccountVo, password string) {
	row := struct {
		models.CustomerAccountVo
		Password string `db:"password"`
	}{}
	err := datasource.GetMasterDb().Get(&row, dao.selectAccountSql+", a.password "+dao.fromAccountSql+" where a.username = ?", username)
	if err == sql.ErrNoRows {
		return nil, ""
	} else if err != nil {
		panic(err)
	}
	return &row.CustomerAccountVo, row.Password
}

func (dao *customerDao) CheckAccountUsernameUnique(username string) int64 {
	var accountId int64
	err := datasource.GetMasterDb().Get(&accountId, "select account_id from customer_account where username = ? limit 1", username)
	if err == sql.ErrNoRows {
		return 0
	} else if err != nil {
		panic(err)
	}
	return accountId
}

func (dao *customerDao) InsertAccount(account *models.CustomerAccountDML) {
	insertSQL := `insert into customer_account(
		account_id, customer_id, username, password, real_name, phone, email, is_main, status,
		remark, create_by, create_time, update_by, update_time
	) values (
		:account_id, :customer_id, :username, :password, :real_name, :phone, :email, :is_main, :status,
		:remark, :create_by, now(), :update_by, now()
	)`
	if _, err := datasource.GetMasterDb().NamedExec(insertSQL, account); err != nil {
		panic(err)
	}
}

func (dao *customerDao) UpdateAccount(account *models.CustomerAccountDML) {
	updateSQL := `update customer_account set real_name = :real_name, phone = :phone, email = :email,
		is_main = :is_main, status = :status, remark = :remark, update_by = :update_by, update_time = now()
		where account_id = :account_id`
	if _, err := datasource.GetMasterDb().NamedExec(updateSQL, account); err != nil {
		panic(err)
	}
}

func (dao *customerDao) ResetAccountPassword(accountId int64, password string) {
	if _, err := datasource.GetMasterDb().Exec("update customer_account set password = ?, update_time = now() where account_id = ?", password, accountId); err != nil {
		panic(err)
	}
}

func (dao *customerDao) DeleteAccountByIds(accountIds []int64) {
	query, args, err := sqlx.In("delete from customer_account where account_id in (?)", accountIds)
	if err != nil {
		panic(err)
	}
	if _, err = datasource.GetMasterDb().Exec(query, args...); err != nil {
		panic(err)
	}
}

func (dao *customerDao) UpdateAccountLoginInfo(accountId int64) {
	if _, err := datasource.GetMasterDb().Exec("update customer_account set last_login_time = now() where account_id = ?", accountId); err != nil {
		panic(err)
	}
}
