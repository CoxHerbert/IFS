package service

import (
	"baize/app/portal/dao"
	"baize/app/portal/models"
	"baize/app/system/service/systemService/systemServiceImpl"
	"baize/app/utils/exceLize"
	"baize/app/utils/snowflake"
	"fmt"
)

var contactServiceImpl *contactService

func init() {
	contactServiceImpl = &contactService{contactDao: dao.GetContactDao()}
}

type contactService struct {
	contactDao interface {
		InsertContact(contact *models.ContactDML)
		SelectContactList(contact *models.ContactDQL) (list []*models.ContactVo, total *int64)
		SelectContactById(contactId int64) (contact *models.ContactVo)
		UpdateContact(contact *models.ContactDML)
		DeleteContactByIds(contactIds []int64)
	}
}

func GetContactService() *contactService {
	return contactServiceImpl
}

func (contactService *contactService) InsertContact(contact *models.ContactDML) string {
	contact.ContactId = snowflake.GenID()
	contact.LeadNo = fmt.Sprintf("CT%d", contact.ContactId)
	if contact.Source == "" {
		contact.Source = "portal-contact"
	}
	contact.Status = "10"
	contact.CreateBy = "portal"
	contact.UpdateBy = "portal"
	contactService.contactDao.InsertContact(contact)
	return contact.LeadNo
}

func (contactService *contactService) SelectContactList(contact *models.ContactDQL) (list []*models.ContactVo, total *int64) {
	return contactService.contactDao.SelectContactList(contact)
}

func (contactService *contactService) SelectContactById(contactId int64) (contact *models.ContactVo) {
	return contactService.contactDao.SelectContactById(contactId)
}

func (contactService *contactService) UpdateContact(contact *models.ContactDML) {
	contactService.contactDao.UpdateContact(contact)
}

func (contactService *contactService) DeleteContactByIds(contactIds []int64) {
	contactService.contactDao.DeleteContactByIds(contactIds)
}

func (contactService *contactService) ExportContact(contact *models.ContactDQL) (data []byte) {
	list, _ := contactService.contactDao.SelectContactList(contact)
	return exceLize.SetRows(models.ContactListToRows(list, contactStatusLabels()))
}

func contactStatusLabels() map[string]string {
	labels := make(map[string]string)
	for _, item := range systemServiceImpl.GetDictDataService().SelectDictDataByType("portal_contact_status") {
		labels[item.DictValue] = item.DictLabel
	}
	return labels
}
