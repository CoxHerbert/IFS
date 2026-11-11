package service

import (
	"baize/app/freight/dao"
	"baize/app/freight/models"
	"baize/app/utils/snowflake"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

type receiptService struct {
	receiptDao interface {
		SelectList(*models.ReceiptDQL) ([]*models.ReceiptVo, *int64)
		SelectById(int64) *models.ReceiptVo
		Insert(*models.ReceiptDML, []*models.ReceiptAllocationDML)
		Delete(int64)
		Allocate(int64, *models.ReceiptAllocationDML) error
		InsertDeclaration(*models.PaymentDeclarationDML)
		SelectDeclarationList(*models.PaymentDeclarationDQL) ([]*models.PaymentDeclarationVo, *int64)
		SelectDeclarationById(int64) *models.PaymentDeclarationVo
		ApproveDeclaration(int64, *models.ReceiptDML, *models.ReceiptAllocationDML, string, string) error
		RejectDeclaration(int64, string, string) error
		SelectPaymentLedger(*models.PaymentLedgerDQL) ([]*models.PaymentLedgerVo, *int64)
		UpdatePayableAmount(int64, float64) error
		SelectShipmentCharges(int64) []*models.ShipmentChargeItem
		ReplaceShipmentCharges(int64, []*models.ShipmentChargeItem, string) error
	}
	shipmentDao interface {
		SelectShipmentById(int64) *models.ShipmentPlanVo
	}
}

var supportedReceiptCurrencies = map[string]struct{}{
	"CNY": {}, "USD": {}, "EUR": {}, "GBP": {}, "JPY": {}, "HKD": {}, "SGD": {}, "AUD": {}, "CAD": {},
}

func normalizeReceiptCurrency(value string) (string, error) {
	currency := strings.ToUpper(strings.TrimSpace(value))
	if currency == "" {
		currency = "CNY"
	}
	if _, ok := supportedReceiptCurrencies[currency]; !ok {
		return "", errors.New("不支持的收款币种")
	}
	return currency, nil
}

func (s *receiptService) SelectDeclarationList(q *models.PaymentDeclarationDQL) ([]*models.PaymentDeclarationVo, *int64) {
	return s.receiptDao.SelectDeclarationList(q)
}
func (s *receiptService) SelectDeclarationById(id int64) *models.PaymentDeclarationVo {
	return s.receiptDao.SelectDeclarationById(id)
}
func (s *receiptService) ApproveDeclaration(id int64, reviewBy, remark string) (*models.ReceiptVo, error) {
	item := s.receiptDao.SelectDeclarationById(id)
	if item == nil {
		return nil, errors.New("付款申报不存在")
	}
	if item.Status != "PENDING" {
		return nil, errors.New("申报已处理，请勿重复审核")
	}
	receiptId := snowflake.GenID()
	receipt := &models.ReceiptDML{ReceiptId: receiptId, ReceiptNo: fmt.Sprintf("RC%s%06d", time.Now().Format("20060102"), receiptId%1000000), CustomerId: item.CustomerId, CustomerName: item.CustomerName, Amount: item.Amount, Currency: item.Currency, ReceiptTime: item.PaymentTime, PaymentMethod: "BANK_TRANSFER", Status: "ALLOCATED", VoucherUrl: item.VoucherUrl, VoucherName: item.VoucherName, Remark: "由付款申报 " + item.DeclarationNo + " 审核生成", CreateBy: reviewBy}
	allocation := &models.ReceiptAllocationDML{AllocationId: snowflake.GenID(), ReceiptId: receiptId, ShipmentId: item.ShipmentId, AllocatedAmount: item.Amount}
	if err := s.receiptDao.ApproveDeclaration(id, receipt, allocation, reviewBy, remark); err != nil {
		return nil, err
	}
	return s.receiptDao.SelectById(receiptId), nil
}
func (s *receiptService) RejectDeclaration(id int64, reviewBy, remark string) error {
	if strings.TrimSpace(remark) == "" {
		return errors.New("驳回时必须填写原因")
	}
	return s.receiptDao.RejectDeclaration(id, reviewBy, remark)
}

func (s *receiptService) CreatePaymentDeclaration(item *models.PaymentDeclarationDML) error {
	if item == nil || item.CustomerId == 0 || item.Amount <= 0 || strings.TrimSpace(item.VoucherUrl) == "" {
		return errors.New("付款金额和付款凭证不能为空")
	}
	plan := s.shipmentDao.SelectShipmentById(item.ShipmentId)
	if plan == nil || plan.CustomerId != item.CustomerId {
		return errors.New("关联出货计划不存在")
	}
	if item.Currency != "" && strings.ToUpper(strings.TrimSpace(item.Currency)) != plan.Currency {
		return errors.New("付款币种必须与出货计划结算币种一致")
	}
	item.Currency = plan.Currency
	item.DeclarationId = snowflake.GenID()
	item.DeclarationNo = fmt.Sprintf("PD%s%06d", time.Now().Format("20060102"), item.DeclarationId%1000000)
	item.Status = "PENDING"
	var err error
	item.Currency, err = normalizeReceiptCurrency(item.Currency)
	if err != nil {
		return err
	}
	item.ShipmentNo = plan.ShipmentNo
	s.receiptDao.InsertDeclaration(item)
	return nil
}

var receiptServiceImpl = &receiptService{receiptDao: dao.GetReceiptDao(), shipmentDao: dao.GetShipmentDao()}

func GetReceiptService() *receiptService { return receiptServiceImpl }
func (s *receiptService) SelectPaymentLedger(q *models.PaymentLedgerDQL) ([]*models.PaymentLedgerVo, *int64) {
	return s.receiptDao.SelectPaymentLedger(q)
}
func (s *receiptService) UpdatePayableAmount(shipmentId int64, amount float64) error {
	if shipmentId == 0 || amount < 0 {
		return errors.New("应付金额不能小于0")
	}
	return s.receiptDao.UpdatePayableAmount(shipmentId, round2(amount))
}
func (s *receiptService) SelectShipmentCharges(shipmentId int64) []*models.ShipmentChargeItem {
	return s.receiptDao.SelectShipmentCharges(shipmentId)
}
func (s *receiptService) ReplaceShipmentCharges(shipmentId int64, items []*models.ShipmentChargeItem, username string) error {
	if shipmentId == 0 {
		return errors.New("出货计划不能为空")
	}
	plan := s.shipmentDao.SelectShipmentById(shipmentId)
	if plan == nil {
		return errors.New("出货计划不存在")
	}
	for _, item := range items {
		if item == nil || strings.TrimSpace(item.FeeName) == "" || item.Amount <= 0 {
			return errors.New("费用名称和金额不能为空")
		}
		item.FeeName = strings.TrimSpace(item.FeeName)
		item.Amount = round2(item.Amount)
		item.Currency = plan.Currency
	}
	return s.receiptDao.ReplaceShipmentCharges(shipmentId, items, username)
}

func (s *receiptService) SelectList(q *models.ReceiptDQL) ([]*models.ReceiptVo, *int64) {
	return s.receiptDao.SelectList(q)
}
func (s *receiptService) SelectById(id int64) *models.ReceiptVo { return s.receiptDao.SelectById(id) }

func (s *receiptService) Create(req *models.ReceiptCreateReq, voucherUrl, voucherName, username string) (*models.ReceiptVo, error) {
	if req == nil || req.CustomerId == 0 || req.Amount <= 0 {
		return nil, errors.New("客户和收款金额不能为空")
	}
	currency, err := normalizeReceiptCurrency(req.Currency)
	if err != nil {
		return nil, err
	}
	allocated := 0.0
	seen := map[int64]bool{}
	allocations := make([]*models.ReceiptAllocationDML, 0, len(req.Allocations))
	for _, item := range req.Allocations {
		if item == nil || item.ShipmentId == 0 || item.AllocatedAmount <= 0 {
			continue
		}
		if seen[item.ShipmentId] {
			return nil, errors.New("同一出货计划不能重复核销")
		}
		plan := s.shipmentDao.SelectShipmentById(item.ShipmentId)
		if plan == nil || plan.CustomerId != req.CustomerId {
			return nil, errors.New("核销的出货计划不属于所选客户")
		}
		if plan.Currency != currency {
			return nil, errors.New("收款币种必须与出货计划结算币种一致")
		}
		seen[item.ShipmentId] = true
		allocated += item.AllocatedAmount
		allocations = append(allocations, &models.ReceiptAllocationDML{AllocationId: snowflake.GenID(), ShipmentId: item.ShipmentId, AllocatedAmount: round2(item.AllocatedAmount)})
	}
	if allocated-req.Amount > 0.001 {
		return nil, errors.New("核销金额合计不能超过收款金额")
	}
	id := snowflake.GenID()
	for _, item := range allocations {
		item.ReceiptId = id
	}
	receiptTime := strings.TrimSpace(req.ReceiptTime)
	if receiptTime == "" {
		receiptTime = time.Now().Format("2006-01-02 15:04:05")
	}
	status := "UNALLOCATED"
	if allocated > 0 {
		status = "PARTIAL"
	}
	if math.Abs(allocated-req.Amount) < 0.001 {
		status = "ALLOCATED"
	}
	receipt := &models.ReceiptDML{ReceiptId: id, ReceiptNo: fmt.Sprintf("RC%s%06d", time.Now().Format("20060102"), id%1000000), CustomerId: req.CustomerId, CustomerName: req.CustomerName,
		Amount: round2(req.Amount), Currency: currency, ReceiptTime: receiptTime, PaymentMethod: req.PaymentMethod, Status: status, VoucherUrl: voucherUrl, VoucherName: voucherName, Remark: req.Remark, CreateBy: username}
	s.receiptDao.Insert(receipt, allocations)
	return s.receiptDao.SelectById(id), nil
}

func (s *receiptService) Delete(id int64) error {
	if s.receiptDao.SelectById(id) == nil {
		return errors.New("收款单不存在")
	}
	s.receiptDao.Delete(id)
	return nil
}

func (s *receiptService) Allocate(receiptId, shipmentId int64, amount float64) (*models.ReceiptVo, error) {
	receipt := s.receiptDao.SelectById(receiptId)
	if receipt == nil {
		return nil, errors.New("收款单不存在")
	}
	plan := s.shipmentDao.SelectShipmentById(shipmentId)
	if plan == nil || plan.CustomerId != receipt.CustomerId {
		return nil, errors.New("出货计划不存在或不属于收款客户")
	}
	if plan.Currency != receipt.Currency {
		return nil, errors.New("收款币种必须与出货计划结算币种一致")
	}
	allocation := &models.ReceiptAllocationDML{AllocationId: snowflake.GenID(), ReceiptId: receiptId, ShipmentId: shipmentId, AllocatedAmount: round2(amount)}
	if amount <= 0 {
		return nil, errors.New("核销金额必须大于0")
	}
	if err := s.receiptDao.Allocate(receiptId, allocation); err != nil {
		return nil, err
	}
	return s.receiptDao.SelectById(receiptId), nil
}
