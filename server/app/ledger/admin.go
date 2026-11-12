package ledger

import (
	"baize/app/common/baize/baizeContext"
	"baize/app/common/datasource"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type adminAccount struct {
	AccountID    int64   `db:"account_id" json:"accountId"`
	UID          string  `db:"uid" json:"uid"`
	LoginID      string  `db:"login_id" json:"loginId"`
	CloudSync    bool    `db:"cloud_sync_enabled" json:"cloudSyncEnabled"`
	CreateTime   string  `db:"create_time" json:"createTime"`
	RecordCount  int64   `db:"record_count" json:"recordCount"`
	IncomeTotal  float64 `db:"income_total" json:"incomeTotal"`
	ExpenseTotal float64 `db:"expense_total" json:"expenseTotal"`
}

type adminRecord struct {
	ID             string  `db:"record_id" json:"id"`
	Date           string  `db:"record_date" json:"date"`
	Merchant       string  `db:"merchant" json:"merchant"`
	Category       string  `db:"category" json:"category"`
	Amount         float64 `db:"amount" json:"amount"`
	Type           string  `db:"record_type" json:"type"`
	Source         string  `db:"source" json:"source"`
	ImportBatchID  string  `db:"import_batch_id" json:"importBatchId"`
	ImportFileName string  `db:"import_file_name" json:"importFileName"`
}

func adminPagination(c *gin.Context) (int, int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	return page, size, (page - 1) * size
}

// AdminAccounts lists cloud-enabled ledger accounts for authorized administrators.
func AdminAccounts(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	_, size, offset := adminPagination(c)
	keyword := strings.TrimSpace(c.Query("loginId"))
	filter := " WHERE a.login_id LIKE ?"
	pattern := "%" + strings.ReplaceAll(strings.ReplaceAll(keyword, "\\", "\\\\"), "%", "\\%") + "%"
	var total int64
	if err := datasource.GetMasterDb().Get(&total, `SELECT COUNT(*) FROM ledger_account a`+filter, pattern); err != nil {
		bzc.BzError()
		return
	}
	rows := make([]adminAccount, 0)
	query := `SELECT a.account_id, CONCAT('QJ-', UPPER(CONV(a.account_id, 10, 36))) AS uid, a.login_id, a.cloud_sync_enabled, DATE_FORMAT(a.create_time, '%Y-%m-%d %H:%i:%s') AS create_time,
		COUNT(r.record_id) AS record_count,
		COALESCE(SUM(CASE WHEN r.record_type='income' THEN r.amount ELSE 0 END), 0) AS income_total,
		COALESCE(SUM(CASE WHEN r.record_type='expense' THEN r.amount ELSE 0 END), 0) AS expense_total
		FROM ledger_account a LEFT JOIN ledger_record r ON r.account_id=a.account_id` + filter +
		` GROUP BY a.account_id, a.login_id, a.cloud_sync_enabled, a.create_time ORDER BY a.create_time DESC LIMIT ? OFFSET ?`
	if err := datasource.GetMasterDb().Select(&rows, query, pattern, size, offset); err != nil {
		bzc.BzError()
		return
	}
	bzc.SuccessListData(rows, &total)
}

// AdminAccountRecords returns a selected account and its cloud-stored ledger records.
func AdminAccountRecords(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	accountID, err := strconv.ParseInt(c.Param("accountId"), 10, 64)
	if err != nil || accountID <= 0 {
		bzc.ParameterError()
		return
	}
	var account adminAccount
	if err := datasource.GetMasterDb().Get(&account, `SELECT account_id, CONCAT('QJ-', UPPER(CONV(account_id, 10, 36))) AS uid, login_id, cloud_sync_enabled, DATE_FORMAT(create_time, '%Y-%m-%d %H:%i:%s') AS create_time FROM ledger_account WHERE account_id=?`, accountID); err != nil {
		bzc.Waring("记账用户不存在")
		return
	}
	page, size, offset := adminPagination(c)
	var total int64
	if err := datasource.GetMasterDb().Get(&total, `SELECT COUNT(*) FROM ledger_record WHERE account_id=?`, accountID); err != nil {
		bzc.BzError()
		return
	}
	rows := make([]adminRecord, 0)
	if err := datasource.GetMasterDb().Select(&rows, `SELECT record_id, record_date, merchant, category, amount, record_type, source, import_batch_id, import_file_name FROM ledger_record WHERE account_id=? ORDER BY record_date DESC LIMIT ? OFFSET ?`, accountID, size, offset); err != nil {
		bzc.BzError()
		return
	}
	bzc.SuccessData(gin.H{"account": account, "rows": rows, "total": total, "pageNum": page, "pageSize": size})
}

// AdminDeleteAccountRecord removes an invalid cloud ledger record.
func AdminDeleteAccountRecord(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	accountID, err := strconv.ParseInt(c.Param("accountId"), 10, 64)
	if err != nil || accountID <= 0 {
		bzc.ParameterError()
		return
	}
	recordID, err := url.PathUnescape(c.Param("recordId"))
	if err != nil || recordID == "" {
		bzc.ParameterError()
		return
	}
	result, err := datasource.GetMasterDb().Exec(`DELETE FROM ledger_record WHERE account_id=? AND record_id=?`, accountID, recordID)
	if err != nil {
		bzc.BzError()
		return
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		bzc.Waring("账单记录不存在")
		return
	}
	bzc.Success()
}
