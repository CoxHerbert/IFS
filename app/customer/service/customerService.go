package service

import (
	"baize/app/customer/dao"
	"baize/app/customer/models"
	"baize/app/setting"
	"baize/app/utils/bCryptPasswordEncoder"
	"baize/app/utils/snowflake"
	"errors"
	"fmt"
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

func GenCustomerToken(account *models.CustomerAccountVo) (string, error) {
	expireSeconds := setting.Conf.TokenConfig.ExpireTime
	if expireSeconds <= 0 {
		expireSeconds = 7200
	}
	claims := CustomerClaims{
		AccountId:  account.AccountId,
		CustomerId: account.CustomerId,
		Username:   account.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(expireSeconds) * time.Second).Unix(),
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
