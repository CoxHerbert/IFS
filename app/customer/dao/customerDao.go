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
			sales_user_id, sales_user_name, status, remark, create_by, create_time, update_by, update_time `,
		fromCustomerSql: ` from customer`,
		selectContactSql: `select contact_id, customer_id, contact_name, position, phone, email, wechat,
			is_primary, status, remark, create_by, create_time, update_by, update_time `,
		fromContactSql: ` from customer_contact`,
		selectAccountSql: `select a.account_id, a.customer_id, c.customer_no, c.customer_name, c.company_name,
			a.username, a.real_name, a.phone, a.email, a.is_main, a.status,
			(
				select ifnull(group_concat(distinct r.role_name order by r.role_sort separator ','), '')
				from customer_workspace_account_role ar
				left join customer_workspace_role r on ar.role_id = r.role_id and r.del_flag = '0'
				where ar.account_id = a.account_id
			) as role_names,
			a.last_login_time,
			a.remark, a.create_by, a.create_time, a.update_by, a.update_time `,
		fromAccountSql: ` from customer_account a left join customer c on a.customer_id = c.customer_id`,
		selectPortalMenuSql: `select menu_id, parent_id, menu_name, order_num, path, component, is_cache, menu_type,
			visible, status, ifnull(perms, '') as perms, ifnull(icon, '') as icon, remark,
			create_by, create_time, update_by, update_time `,
		fromPortalMenuSql: ` from customer_workspace_menu`,
		selectPortalRoleSql: `select role_id, role_name, role_key, role_sort, status, del_flag, remark,
			create_by, create_time, update_by, update_time `,
		fromPortalRoleSql: ` from customer_workspace_role`,
	}
}

type customerDao struct {
	selectCustomerSql   string
	fromCustomerSql     string
	selectContactSql    string
	fromContactSql      string
	selectAccountSql    string
	fromAccountSql      string
	selectPortalMenuSql string
	fromPortalMenuSql   string
	selectPortalRoleSql string
	fromPortalRoleSql   string
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
	if customer.SalesUserId != 0 {
		whereSql += " AND sales_user_id = :sales_user_id"
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
	err := datasource.GetMasterDb().Select(&list, `select customer_id, customer_no, customer_name, company_name, sales_user_id, sales_user_name from customer`+whereSql+` order by create_time desc limit 100`, args...)
	if err != nil {
		panic(err)
	}
	return list
}

func (dao *customerDao) InsertCustomer(customer *models.CustomerDML) {
	insertSQL := `insert into customer(
		customer_id, customer_no, customer_name, company_name, contact_name, phone, email,
		sales_user_id, sales_user_name, status, remark, create_by, create_time, update_by, update_time
	) values (
		:customer_id, :customer_no, :customer_name, :company_name, :contact_name, :phone, :email,
		:sales_user_id, :sales_user_name, :status, :remark, :create_by, now(), :update_by, now()
	)`
	if _, err := datasource.GetMasterDb().NamedExec(insertSQL, customer); err != nil {
		panic(err)
	}
}

func (dao *customerDao) UpdateCustomer(customer *models.CustomerDML) {
	updateSQL := `update customer set customer_name = :customer_name, company_name = :company_name,
		contact_name = :contact_name, phone = :phone, email = :email, status = :status,
		sales_user_id = :sales_user_id, sales_user_name = :sales_user_name,
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

func (dao *customerDao) SelectAccountPasswordById(accountId int64) string {
	password := new(string)
	err := datasource.GetMasterDb().Get(password, "select password from customer_account where account_id = ?", accountId)
	if err == sql.ErrNoRows {
		return ""
	} else if err != nil {
		panic(err)
	}
	return *password
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

func (dao *customerDao) UpdatePortalProfile(accountId int64, realName string, phone string, email string, updateBy string) {
	if _, err := datasource.GetMasterDb().Exec(
		`update customer_account set real_name = ?, phone = ?, email = ?, update_by = ?, update_time = now() where account_id = ?`,
		realName,
		phone,
		email,
		updateBy,
		accountId,
	); err != nil {
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

func (dao *customerDao) SelectPortalMenuList(menu *models.CustomerPortalMenuDQL) (list []*models.CustomerPortalMenuVo) {
	whereSql := ``
	if menu.MenuName != "" {
		whereSql += " AND menu_name like concat('%', :menu_name, '%')"
	}
	if menu.Visible != "" {
		whereSql += " AND visible = :visible"
	}
	if menu.Status != "" {
		whereSql += " AND status = :status"
	}
	if whereSql != "" {
		whereSql = " where " + whereSql[4:]
	}
	rows, err := datasource.GetMasterDb().NamedQuery(dao.selectPortalMenuSql+dao.fromPortalMenuSql+whereSql+" order by parent_id, cast(order_num as unsigned), menu_id", menu)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	list = make([]*models.CustomerPortalMenuVo, 0)
	for rows.Next() {
		item := new(models.CustomerPortalMenuVo)
		if err := rows.StructScan(item); err != nil {
			panic(err)
		}
		list = append(list, item)
	}
	return
}

func (dao *customerDao) SelectPortalMenuById(menuId int64) *models.CustomerPortalMenuVo {
	menu := new(models.CustomerPortalMenuVo)
	err := datasource.GetMasterDb().Get(menu, dao.selectPortalMenuSql+dao.fromPortalMenuSql+" where menu_id = ?", menuId)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		panic(err)
	}
	return menu
}

func (dao *customerDao) InsertPortalMenu(menu *models.CustomerPortalMenuDML) {
	insertSQL := `insert into customer_workspace_menu(
		menu_id, parent_id, menu_name, order_num, path, component, is_cache, menu_type, visible, status, perms, icon,
		remark, create_by, create_time, update_by, update_time
	) values (
		:menu_id, :parent_id, :menu_name, :order_num, :path, :component, :is_cache, :menu_type, :visible, :status, :perms, :icon,
		:remark, :create_by, now(), :update_by, now()
	)`
	if _, err := datasource.GetMasterDb().NamedExec(insertSQL, menu); err != nil {
		panic(err)
	}
}

func (dao *customerDao) UpdatePortalMenu(menu *models.CustomerPortalMenuDML) {
	updateSQL := `update customer_workspace_menu set parent_id = :parent_id, menu_name = :menu_name, order_num = :order_num,
		path = :path, component = :component, menu_type = :menu_type, visible = :visible, status = :status,
		is_cache = :is_cache, perms = :perms, icon = :icon, remark = :remark, update_by = :update_by, update_time = now()
		where menu_id = :menu_id`
	if _, err := datasource.GetMasterDb().NamedExec(updateSQL, menu); err != nil {
		panic(err)
	}
}

func (dao *customerDao) DeletePortalMenuById(menuId int64) {
	if _, err := datasource.GetMasterDb().Exec("delete from customer_workspace_menu where menu_id = ?", menuId); err != nil {
		panic(err)
	}
}

func (dao *customerDao) HasPortalMenuChildByMenuId(menuId int64) int {
	var count int
	if err := datasource.GetMasterDb().Get(&count, "select count(1) from customer_workspace_menu where parent_id = ?", menuId); err != nil {
		panic(err)
	}
	return count
}

func (dao *customerDao) CheckPortalMenuNameUnique(menuName string, parentId int64) int64 {
	var menuId int64
	err := datasource.GetMasterDb().Get(&menuId, "select menu_id from customer_workspace_menu where menu_name = ? and parent_id = ? limit 1", menuName, parentId)
	if err == sql.ErrNoRows {
		return 0
	} else if err != nil {
		panic(err)
	}
	return menuId
}

func (dao *customerDao) CheckPortalMenuExistRole(menuId int64) int {
	var count int
	if err := datasource.GetMasterDb().Get(&count, "select count(1) from customer_workspace_role_menu where menu_id = ?", menuId); err != nil {
		panic(err)
	}
	return count
}

func (dao *customerDao) SelectPortalRoleList(role *models.CustomerPortalRoleDQL) (list []*models.CustomerPortalRoleVo, total *int64) {
	whereSql := ` where del_flag = '0'`
	if role.RoleName != "" {
		whereSql += " AND role_name like concat('%', :role_name, '%')"
	}
	if role.RoleKey != "" {
		whereSql += " AND role_key like concat('%', :role_key, '%')"
	}
	if role.Status != "" {
		whereSql += " AND status = :status"
	}
	if role.BeginTime != "" {
		whereSql += " AND date_format(create_time,'%Y-%m-%d') >= :begin_time"
	}
	if role.EndTime != "" {
		whereSql += " AND date_format(create_time,'%Y-%m-%d') <= :end_time"
	}
	countRow, err := datasource.GetMasterDb().NamedQuery(constants.MysqlCount+dao.fromPortalRoleSql+whereSql, role)
	if err != nil {
		panic(err)
	}
	total = new(int64)
	if countRow.Next() {
		countRow.Scan(total)
	}
	countRow.Close()
	list = make([]*models.CustomerPortalRoleVo, 0, role.Size)
	if *total > role.Offset {
		whereSql += " order by role_sort asc, create_time desc"
		if role.Limit != "" {
			whereSql += role.Limit
		}
		rows, err := datasource.GetMasterDb().NamedQuery(dao.selectPortalRoleSql+dao.fromPortalRoleSql+whereSql, role)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		for rows.Next() {
			item := new(models.CustomerPortalRoleVo)
			if err := rows.StructScan(item); err != nil {
				panic(err)
			}
			list = append(list, item)
		}
	}
	return
}

func (dao *customerDao) SelectPortalRoleById(roleId int64) *models.CustomerPortalRoleVo {
	role := new(models.CustomerPortalRoleVo)
	err := datasource.GetMasterDb().Get(role, dao.selectPortalRoleSql+dao.fromPortalRoleSql+" where role_id = ? and del_flag = '0'", roleId)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		panic(err)
	}
	return role
}

func (dao *customerDao) InsertPortalRole(role *models.CustomerPortalRoleDML) {
	insertSQL := `insert into customer_workspace_role(
		role_id, role_name, role_key, role_sort, status, del_flag, remark, create_by, create_time, update_by, update_time
	) values (
		:role_id, :role_name, :role_key, :role_sort, :status, '0', :remark, :create_by, now(), :update_by, now()
	)`
	if _, err := datasource.GetMasterDb().NamedExec(insertSQL, role); err != nil {
		panic(err)
	}
}

func (dao *customerDao) UpdatePortalRole(role *models.CustomerPortalRoleDML) {
	updateSQL := `update customer_workspace_role set role_name = :role_name, role_key = :role_key, role_sort = :role_sort,
		status = :status, remark = :remark, update_by = :update_by, update_time = now() where role_id = :role_id`
	if _, err := datasource.GetMasterDb().NamedExec(updateSQL, role); err != nil {
		panic(err)
	}
}

func (dao *customerDao) UpdatePortalRoleStatus(roleId int64, status string, updateBy string) {
	if _, err := datasource.GetMasterDb().Exec("update customer_workspace_role set status = ?, update_by = ?, update_time = now() where role_id = ?", status, updateBy, roleId); err != nil {
		panic(err)
	}
}

func (dao *customerDao) DeletePortalRoleByIds(roleIds []int64, updateBy string) {
	query, args, err := sqlx.In("update customer_workspace_role set del_flag = '2', update_by = ?, update_time = now() where role_id in (?)", updateBy, roleIds)
	if err != nil {
		panic(err)
	}
	if _, err = datasource.GetMasterDb().Exec(query, args...); err != nil {
		panic(err)
	}
}

func (dao *customerDao) CheckPortalRoleNameUnique(roleName string) int64 {
	var roleId int64
	err := datasource.GetMasterDb().Get(&roleId, "select role_id from customer_workspace_role where role_name = ? and del_flag = '0' limit 1", roleName)
	if err == sql.ErrNoRows {
		return 0
	} else if err != nil {
		panic(err)
	}
	return roleId
}

func (dao *customerDao) CheckPortalRoleKeyUnique(roleKey string) int64 {
	var roleId int64
	err := datasource.GetMasterDb().Get(&roleId, "select role_id from customer_workspace_role where role_key = ? and del_flag = '0' limit 1", roleKey)
	if err == sql.ErrNoRows {
		return 0
	} else if err != nil {
		panic(err)
	}
	return roleId
}

func (dao *customerDao) SelectPortalRoleMenuIds(roleId int64) (menuIds []string) {
	menuIds = make([]string, 0)
	if err := datasource.GetMasterDb().Select(&menuIds, "select cast(menu_id as char) from customer_workspace_role_menu where role_id = ?", roleId); err != nil {
		panic(err)
	}
	return
}

func (dao *customerDao) ReplacePortalRoleMenus(roleId int64, menuIds []int64) {
	tx := datasource.GetMasterDb().MustBegin()
	if _, err := tx.Exec("delete from customer_workspace_role_menu where role_id = ?", roleId); err != nil {
		tx.Rollback()
		panic(err)
	}
	if len(menuIds) > 0 {
		list := make([]*models.CustomerPortalRoleMenu, 0, len(menuIds))
		for _, menuId := range menuIds {
			list = append(list, &models.CustomerPortalRoleMenu{RoleId: roleId, MenuId: menuId})
		}
		if _, err := tx.NamedExec("insert into customer_workspace_role_menu(role_id, menu_id) values (:role_id, :menu_id)", list); err != nil {
			tx.Rollback()
			panic(err)
		}
	}
	if err := tx.Commit(); err != nil {
		panic(err)
	}
}

func (dao *customerDao) CountAccountRoleByRoleIds(roleIds []int64) int {
	query, args, err := sqlx.In("select count(1) from customer_workspace_account_role where role_id in (?)", roleIds)
	if err != nil {
		panic(err)
	}
	var count int
	if err = datasource.GetMasterDb().Get(&count, query, args...); err != nil {
		panic(err)
	}
	return count
}

func (dao *customerDao) SelectPortalRoleOptions() (list []*models.CustomerPortalRoleOptionVo) {
	list = make([]*models.CustomerPortalRoleOptionVo, 0)
	if err := datasource.GetMasterDb().Select(&list, "select role_id, role_name from customer_workspace_role where status = '0' and del_flag = '0' order by role_sort asc, role_id asc"); err != nil {
		panic(err)
	}
	return
}

func (dao *customerDao) SelectAccountRoleIds(accountId int64) (roleIds []string) {
	roleIds = make([]string, 0)
	if err := datasource.GetMasterDb().Select(&roleIds, "select cast(role_id as char) from customer_workspace_account_role where account_id = ?", accountId); err != nil {
		panic(err)
	}
	return
}

func (dao *customerDao) ReplaceAccountRoles(accountId int64, roleIds []int64) {
	tx := datasource.GetMasterDb().MustBegin()
	if _, err := tx.Exec("delete from customer_workspace_account_role where account_id = ?", accountId); err != nil {
		tx.Rollback()
		panic(err)
	}
	if len(roleIds) > 0 {
		list := make([]*models.CustomerPortalAccountRole, 0, len(roleIds))
		for _, roleId := range roleIds {
			list = append(list, &models.CustomerPortalAccountRole{AccountId: accountId, RoleId: roleId})
		}
		if _, err := tx.NamedExec("insert into customer_workspace_account_role(account_id, role_id) values (:account_id, :role_id)", list); err != nil {
			tx.Rollback()
			panic(err)
		}
	}
	if err := tx.Commit(); err != nil {
		panic(err)
	}
}

func (dao *customerDao) SelectPortalRolesByAccountId(accountId int64) (roles []string) {
	roles = make([]string, 0)
	if err := datasource.GetMasterDb().Select(&roles, `select r.role_key
		from customer_workspace_role r
		inner join customer_workspace_account_role ar on r.role_id = ar.role_id
		where ar.account_id = ? and r.status = '0' and r.del_flag = '0'
		order by r.role_sort asc, r.role_id asc`, accountId); err != nil {
		panic(err)
	}
	return
}

func (dao *customerDao) SelectPortalPermissionsByAccountId(accountId int64) (permissions []string) {
	permissions = make([]string, 0)
	if err := datasource.GetMasterDb().Select(&permissions, `select distinct m.perms
		from customer_workspace_menu m
		inner join customer_workspace_role_menu rm on m.menu_id = rm.menu_id
		inner join customer_workspace_account_role ar on rm.role_id = ar.role_id
		inner join customer_workspace_role r on ar.role_id = r.role_id
		where ar.account_id = ? and r.status = '0' and r.del_flag = '0'
			and m.status = '0' and m.menu_type = 'C' and ifnull(m.perms, '') <> ''
		order by m.perms`, accountId); err != nil {
		panic(err)
	}
	return
}

func (dao *customerDao) SelectPortalMenusByAccountId(accountId int64) (list []*models.CustomerPortalMenuVo) {
	list = make([]*models.CustomerPortalMenuVo, 0)
	err := datasource.GetMasterDb().Select(&list, `select distinct m.menu_id, m.parent_id, m.menu_name, m.order_num, m.path, m.component, m.is_cache, m.menu_type,
			m.visible, m.status, ifnull(m.perms, '') as perms, ifnull(m.icon, '') as icon, m.remark,
			m.create_by, m.create_time, m.update_by, m.update_time from customer_workspace_menu m
		inner join customer_workspace_role_menu rm on m.menu_id = rm.menu_id
		inner join customer_workspace_account_role ar on rm.role_id = ar.role_id
		inner join customer_workspace_role r on ar.role_id = r.role_id
		where ar.account_id = ? and r.status = '0' and r.del_flag = '0'
			and m.menu_type in ('M', 'C') and m.status = '0' and m.visible = '0'
		order by m.parent_id, cast(m.order_num as unsigned), m.menu_id`, accountId)
	if err != nil {
		panic(err)
	}
	return
}
