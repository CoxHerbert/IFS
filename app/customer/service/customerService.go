package service

import (
	"baize/app/customer/dao"
	"baize/app/customer/models"
	"baize/app/setting"
	"baize/app/utils/bCryptPasswordEncoder"
	"baize/app/utils/snowflake"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var customerServiceImpl *customerService

func init() {
	customerServiceImpl = &customerService{customerDao: dao.GetCustomerDao()}
}

type customerService struct {
	customerDao interface {
		SelectCustomerList(customer *models.CustomerDQL) (list []*models.CustomerVo, total *int64)
		SelectCustomerById(customerId int64) *models.CustomerVo
		SelectCustomerOptions(keyword string) []*models.CustomerOptionVo
		InsertCustomer(customer *models.CustomerDML)
		UpdateCustomer(customer *models.CustomerDML)
		DeleteCustomerByIds(customerIds []int64)
		SelectContactList(contact *models.CustomerContactDQL) (list []*models.CustomerContactVo, total *int64)
		SelectContactById(contactId int64) *models.CustomerContactVo
		InsertContact(contact *models.CustomerContactDML)
		UpdateContact(contact *models.CustomerContactDML)
		DeleteContactByIds(contactIds []int64)
		SelectAccountList(account *models.CustomerAccountDQL) (list []*models.CustomerAccountVo, total *int64)
		SelectAccountById(accountId int64) *models.CustomerAccountVo
		SelectAccountByUsername(username string) (account *models.CustomerAccountVo, password string)
		CheckAccountUsernameUnique(username string) int64
		InsertAccount(account *models.CustomerAccountDML)
		UpdateAccount(account *models.CustomerAccountDML)
		ResetAccountPassword(accountId int64, password string)
		DeleteAccountByIds(accountIds []int64)
		UpdateAccountLoginInfo(accountId int64)
		SelectPortalMenuList(menu *models.CustomerPortalMenuDQL) (list []*models.CustomerPortalMenuVo)
		SelectPortalMenuById(menuId int64) *models.CustomerPortalMenuVo
		InsertPortalMenu(menu *models.CustomerPortalMenuDML)
		UpdatePortalMenu(menu *models.CustomerPortalMenuDML)
		DeletePortalMenuById(menuId int64)
		HasPortalMenuChildByMenuId(menuId int64) int
		CheckPortalMenuNameUnique(menuName string, parentId int64) int64
		CheckPortalMenuExistRole(menuId int64) int
		SelectPortalRoleList(role *models.CustomerPortalRoleDQL) (list []*models.CustomerPortalRoleVo, total *int64)
		SelectPortalRoleById(roleId int64) *models.CustomerPortalRoleVo
		InsertPortalRole(role *models.CustomerPortalRoleDML)
		UpdatePortalRole(role *models.CustomerPortalRoleDML)
		UpdatePortalRoleStatus(roleId int64, status string, updateBy string)
		DeletePortalRoleByIds(roleIds []int64, updateBy string)
		CheckPortalRoleNameUnique(roleName string) int64
		CheckPortalRoleKeyUnique(roleKey string) int64
		SelectPortalRoleMenuIds(roleId int64) (menuIds []string)
		ReplacePortalRoleMenus(roleId int64, menuIds []int64)
		CountAccountRoleByRoleIds(roleIds []int64) int
		SelectPortalRoleOptions() (list []*models.CustomerPortalRoleOptionVo)
		SelectAccountRoleIds(accountId int64) (roleIds []string)
		ReplaceAccountRoles(accountId int64, roleIds []int64)
		SelectPortalRolesByAccountId(accountId int64) (roles []string)
		SelectPortalPermissionsByAccountId(accountId int64) (permissions []string)
		SelectPortalMenusByAccountId(accountId int64) (list []*models.CustomerPortalMenuVo)
	}
}

type CustomerClaims struct {
	AccountId  int64  `json:"accountId"`
	CustomerId int64  `json:"customerId"`
	Username   string `json:"username"`
	jwt.StandardClaims
}

func GetCustomerService() *customerService {
	return customerServiceImpl
}

func (service *customerService) SelectCustomerList(customer *models.CustomerDQL) (list []*models.CustomerVo, total *int64) {
	return service.customerDao.SelectCustomerList(customer)
}

func (service *customerService) SelectCustomerById(customerId int64) *models.CustomerVo {
	return service.customerDao.SelectCustomerById(customerId)
}

func (service *customerService) SelectCustomerOptions(keyword string) []*models.CustomerOptionVo {
	return service.customerDao.SelectCustomerOptions(keyword)
}

func (service *customerService) InsertCustomer(customer *models.CustomerDML) {
	customer.CustomerId = snowflake.GenID()
	customer.CustomerNo = fmt.Sprintf("CU%d", customer.CustomerId)
	if customer.Status == "" {
		customer.Status = "0"
	}
	service.customerDao.InsertCustomer(customer)
}

func (service *customerService) UpdateCustomer(customer *models.CustomerDML) {
	service.customerDao.UpdateCustomer(customer)
}

func (service *customerService) DeleteCustomerByIds(customerIds []int64) {
	service.customerDao.DeleteCustomerByIds(customerIds)
}

func (service *customerService) SelectContactList(contact *models.CustomerContactDQL) (list []*models.CustomerContactVo, total *int64) {
	return service.customerDao.SelectContactList(contact)
}

func (service *customerService) SelectContactById(contactId int64) *models.CustomerContactVo {
	return service.customerDao.SelectContactById(contactId)
}

func (service *customerService) InsertContact(contact *models.CustomerContactDML) error {
	if service.customerDao.SelectCustomerById(contact.CustomerId) == nil {
		return errors.New("客户不存在")
	}
	contact.ContactId = snowflake.GenID()
	if contact.Status == "" {
		contact.Status = "0"
	}
	if contact.IsPrimary == "" {
		contact.IsPrimary = "0"
	}
	service.customerDao.InsertContact(contact)
	return nil
}

func (service *customerService) UpdateContact(contact *models.CustomerContactDML) {
	service.customerDao.UpdateContact(contact)
}

func (service *customerService) DeleteContactByIds(contactIds []int64) {
	service.customerDao.DeleteContactByIds(contactIds)
}

func (service *customerService) SelectAccountList(account *models.CustomerAccountDQL) (list []*models.CustomerAccountVo, total *int64) {
	return service.customerDao.SelectAccountList(account)
}

func (service *customerService) SelectAccountById(accountId int64) *models.CustomerAccountVo {
	return service.customerDao.SelectAccountById(accountId)
}

func (service *customerService) InsertAccount(account *models.CustomerAccountDML) error {
	if service.customerDao.SelectCustomerById(account.CustomerId) == nil {
		return errors.New("客户不存在")
	}
	if service.customerDao.CheckAccountUsernameUnique(account.Username) > 0 {
		return errors.New("账号已存在")
	}
	account.AccountId = snowflake.GenID()
	account.Password = bCryptPasswordEncoder.HashPassword(account.Password)
	if account.Status == "" {
		account.Status = "0"
	}
	if account.IsMain == "" {
		account.IsMain = "0"
	}
	service.customerDao.InsertAccount(account)
	return nil
}

func (service *customerService) UpdateAccount(account *models.CustomerAccountDML) {
	service.customerDao.UpdateAccount(account)
}

func (service *customerService) ResetAccountPassword(accountId int64, password string) {
	service.customerDao.ResetAccountPassword(accountId, bCryptPasswordEncoder.HashPassword(password))
}

func (service *customerService) DeleteAccountByIds(accountIds []int64) {
	service.customerDao.DeleteAccountByIds(accountIds)
}

func (service *customerService) Login(login *models.CustomerLoginBody) (*models.CustomerLoginResult, error) {
	account, password := service.customerDao.SelectAccountByUsername(login.Username)
	if account == nil || !bCryptPasswordEncoder.CheckPasswordHash(login.Password, password) {
		return nil, errors.New("账号不存在或密码错误")
	}
	if account.Status != "0" {
		return nil, errors.New("账号已停用")
	}
	service.customerDao.UpdateAccountLoginInfo(account.AccountId)
	token, err := GenCustomerToken(account)
	if err != nil {
		return nil, err
	}
	return &models.CustomerLoginResult{Token: token, User: account}, nil
}

func (service *customerService) SelectAccountProfile(accountId int64) *models.CustomerAccountVo {
	return service.customerDao.SelectAccountById(accountId)
}

func (service *customerService) SelectPortalMenuList(menu *models.CustomerPortalMenuDQL) []*models.CustomerPortalMenuVo {
	return buildPortalMenuTree(service.customerDao.SelectPortalMenuList(menu), 0)
}

func (service *customerService) SelectPortalMenuById(menuId int64) *models.CustomerPortalMenuVo {
	return service.customerDao.SelectPortalMenuById(menuId)
}

func (service *customerService) InsertPortalMenu(menu *models.CustomerPortalMenuDML) {
	menu.MenuId = snowflake.GenID()
	if menu.Visible == "" {
		menu.Visible = "0"
	}
	if menu.Status == "" {
		menu.Status = "0"
	}
	service.customerDao.InsertPortalMenu(menu)
}

func (service *customerService) UpdatePortalMenu(menu *models.CustomerPortalMenuDML) {
	service.customerDao.UpdatePortalMenu(menu)
}

func (service *customerService) DeletePortalMenuById(menuId int64) {
	service.customerDao.DeletePortalMenuById(menuId)
}

func (service *customerService) HasPortalMenuChildByMenuId(menuId int64) bool {
	return service.customerDao.HasPortalMenuChildByMenuId(menuId) > 0
}

func (service *customerService) CheckPortalMenuExistRole(menuId int64) bool {
	return service.customerDao.CheckPortalMenuExistRole(menuId) > 0
}

func (service *customerService) CheckPortalMenuNameUnique(menu *models.CustomerPortalMenuDML) bool {
	menuId := service.customerDao.CheckPortalMenuNameUnique(menu.MenuName, menu.ParentId)
	return menuId != 0 && menuId != menu.MenuId
}

func (service *customerService) SelectPortalRoleList(role *models.CustomerPortalRoleDQL) (list []*models.CustomerPortalRoleVo, total *int64) {
	return service.customerDao.SelectPortalRoleList(role)
}

func (service *customerService) SelectPortalRoleById(roleId int64) *models.CustomerPortalRoleVo {
	return service.customerDao.SelectPortalRoleById(roleId)
}

func (service *customerService) InsertPortalRole(role *models.CustomerPortalRoleDML) {
	role.RoleId = snowflake.GenID()
	if role.Status == "" {
		role.Status = "0"
	}
	service.customerDao.InsertPortalRole(role)
	service.customerDao.ReplacePortalRoleMenus(role.RoleId, strIdsToInt64(role.MenuIds))
}

func (service *customerService) UpdatePortalRole(role *models.CustomerPortalRoleDML) {
	service.customerDao.UpdatePortalRole(role)
	service.customerDao.ReplacePortalRoleMenus(role.RoleId, strIdsToInt64(role.MenuIds))
}

func (service *customerService) UpdatePortalRoleStatus(roleId int64, status string, updateBy string) {
	service.customerDao.UpdatePortalRoleStatus(roleId, status, updateBy)
}

func (service *customerService) DeletePortalRoleByIds(roleIds []int64, updateBy string) {
	service.customerDao.DeletePortalRoleByIds(roleIds, updateBy)
}

func (service *customerService) CheckPortalRoleNameUnique(role *models.CustomerPortalRoleDML) bool {
	roleId := service.customerDao.CheckPortalRoleNameUnique(role.RoleName)
	return roleId != 0 && roleId != role.RoleId
}

func (service *customerService) CheckPortalRoleKeyUnique(role *models.CustomerPortalRoleDML) bool {
	roleId := service.customerDao.CheckPortalRoleKeyUnique(role.RoleKey)
	return roleId != 0 && roleId != role.RoleId
}

func (service *customerService) CountAccountRoleByRoleIds(roleIds []int64) bool {
	return service.customerDao.CountAccountRoleByRoleIds(roleIds) > 0
}

func (service *customerService) SelectPortalRoleMenuTreeselect(roleId int64) map[string]interface{} {
	data := make(map[string]interface{})
	data["menus"] = service.SelectPortalMenuList(new(models.CustomerPortalMenuDQL))
	data["checkedKeys"] = service.customerDao.SelectPortalRoleMenuIds(roleId)
	return data
}

func (service *customerService) SelectPortalRoleOptions() []*models.CustomerPortalRoleOptionVo {
	return service.customerDao.SelectPortalRoleOptions()
}

func (service *customerService) SelectAccountRoleIds(accountId int64) []string {
	return service.customerDao.SelectAccountRoleIds(accountId)
}

func (service *customerService) UpdateAccountRoles(accountId int64, roleIds []string) {
	service.customerDao.ReplaceAccountRoles(accountId, strIdsToInt64(roleIds))
}

func (service *customerService) SelectPortalProfile(accountId int64) *models.CustomerPortalProfile {
	return &models.CustomerPortalProfile{
		User:        service.customerDao.SelectAccountById(accountId),
		Roles:       service.customerDao.SelectPortalRolesByAccountId(accountId),
		Permissions: service.customerDao.SelectPortalPermissionsByAccountId(accountId),
	}
}

func (service *customerService) SelectPortalRouters(accountId int64) []*models.CustomerPortalRoute {
	menus := buildPortalMenuTree(service.customerDao.SelectPortalMenusByAccountId(accountId), 0)
	return buildPortalRoutes(menus)
}

func GenCustomerToken(account *models.CustomerAccountVo) (string, error) {
	expireMinutes := setting.Conf.TokenConfig.ExpireTime
	if expireMinutes <= 0 {
		expireMinutes = 720
	}
	claims := CustomerClaims{
		AccountId:  account.AccountId,
		CustomerId: account.CustomerId,
		Username:   account.Username,
		StandardClaims: jwt.StandardClaims{
			// The shared token config is defined in minutes across the rest of the system.
			ExpiresAt: time.Now().Add(time.Duration(expireMinutes) * time.Minute).Unix(),
			Issuer:    setting.Conf.TokenConfig.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(setting.Conf.TokenConfig.Secret))
}

func ParseCustomerToken(tokenString string) (*CustomerClaims, error) {
	claims := new(CustomerClaims)
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(setting.Conf.TokenConfig.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func strIdsToInt64(ids []string) []int64 {
	values := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id == "" {
			continue
		}
		value, err := strconv.ParseInt(id, 10, 64)
		if err == nil {
			values = append(values, value)
		}
	}
	return values
}

func buildPortalMenuTree(list []*models.CustomerPortalMenuVo, parentId int64) []*models.CustomerPortalMenuVo {
	children := make([]*models.CustomerPortalMenuVo, 0)
	for _, item := range list {
		if item.ParentId != parentId {
			continue
		}
		item.Children = buildPortalMenuTree(list, item.MenuId)
		children = append(children, item)
	}
	sort.Slice(children, func(i, j int) bool {
		left, _ := strconv.Atoi(children[i].OrderNum)
		right, _ := strconv.Atoi(children[j].OrderNum)
		if left == right {
			return children[i].MenuId < children[j].MenuId
		}
		return left < right
	})
	return children
}

func buildPortalRoutes(menus []*models.CustomerPortalMenuVo) []*models.CustomerPortalRoute {
	routes := make([]*models.CustomerPortalRoute, 0, len(menus))
	for _, menu := range menus {
		route := &models.CustomerPortalRoute{
			Name:   fmt.Sprintf("customer-menu-%d", menu.MenuId),
			Path:   menu.Path,
			Hidden: menu.Visible != "0",
			Meta: models.CustomerPortalRouteMeta{
				Title:  menu.MenuName,
				Icon:   menu.Icon,
				MenuId: menu.MenuId,
			},
		}
		if menu.MenuType == "C" {
			route.Component = menu.Component
		}
		if len(menu.Children) > 0 {
			route.Children = buildPortalRoutes(menu.Children)
		}
		routes = append(routes, route)
	}
	return routes
}
