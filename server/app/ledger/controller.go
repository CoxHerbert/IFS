package ledger

import (
	"baize/app/common/baize/baizeContext"
	"baize/app/common/datasource"
	"baize/app/utils/bCryptPasswordEncoder"
	"errors"
	"net/mail"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type credentials struct {
	LoginID  string `json:"loginId" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type accountDTO struct {
	AccountID int64  `json:"accountId"`
	UID       string `json:"uid"`
	LoginID   string `json:"loginId"`
	Token     string `json:"token"`
	CloudSync bool   `json:"cloudSyncEnabled"`
}

type ledgerRecord struct {
	ID             string  `json:"id" db:"record_id" binding:"required"`
	Date           string  `json:"date" db:"record_date" binding:"required"`
	Merchant       string  `json:"merchant" db:"merchant" binding:"required"`
	Category       string  `json:"category" db:"category" binding:"required"`
	Amount         float64 `json:"amount" db:"amount" binding:"required,gt=0"`
	Type           string  `json:"type" db:"record_type" binding:"required,oneof=expense income"`
	Source         string  `json:"source" db:"source" binding:"required,oneof=manual wechat alipay"`
	ImportBatchID  string  `json:"importBatchId" db:"import_batch_id"`
	ImportFileName string  `json:"importFileName" db:"import_file_name"`
}

type recordBatch struct {
	Records []ledgerRecord `json:"records"`
}
type syncConsent struct {
	Enabled *bool `json:"enabled"`
}

var phonePattern = regexp.MustCompile(`^\+?[0-9]{7,20}$`)

func normalizeLoginID(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) < 3 || len(value) > 120 {
		return "", errors.New("账号长度需为 3 到 120 个字符")
	}
	if strings.Contains(value, "@") {
		addr, err := mail.ParseAddress(value)
		if err != nil || addr.Address != value {
			return "", errors.New("邮箱格式不正确")
		}
	} else if !phonePattern.MatchString(value) {
		return "", errors.New("账号需使用邮箱或手机号")
	}
	return strings.ToLower(value), nil
}

func Register(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	var body credentials
	if err := c.ShouldBindJSON(&body); err != nil {
		bzc.ParameterError()
		return
	}
	loginID, err := normalizeLoginID(body.LoginID)
	if err != nil {
		bzc.ErrorMsg(err.Error())
		return
	}
	if len(body.Password) < 8 || len(body.Password) > 72 {
		bzc.ErrorMsg("密码长度需为 8 到 72 个字符")
		return
	}
	hash := bCryptPasswordEncoder.HashPassword(body.Password)
	result, err := datasource.GetMasterDb().Exec(`INSERT INTO ledger_account (login_id, password_hash, cloud_sync_enabled, create_time) VALUES (?, ?, 0, NOW())`, loginID, hash)
	if err != nil {
		bzc.Waring("账号已存在或注册暂不可用")
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		bzc.BzError()
		return
	}
	respondWithToken(bzc, id, loginID)
}

func Login(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	var body credentials
	if err := c.ShouldBindJSON(&body); err != nil {
		bzc.ParameterError()
		return
	}
	loginID, err := normalizeLoginID(body.LoginID)
	if err != nil {
		bzc.Waring("账号或密码错误")
		return
	}
	var row struct {
		ID   int64  `db:"account_id"`
		Hash string `db:"password_hash"`
	}
	if err := datasource.GetMasterDb().Get(&row, `SELECT account_id, password_hash FROM ledger_account WHERE login_id = ?`, loginID); err != nil || !bCryptPasswordEncoder.CheckPasswordHash(body.Password, row.Hash) {
		bzc.Waring("账号或密码错误")
		return
	}
	respondWithToken(bzc, row.ID, loginID)
}

func respondWithToken(bzc *baizeContext.BaiZeContext, id int64, loginID string) {
	token, err := makeToken(id)
	if err != nil {
		bzc.BzError()
		return
	}
	var enabled bool
	if err := datasource.GetMasterDb().Get(&enabled, `SELECT cloud_sync_enabled FROM ledger_account WHERE account_id = ?`, id); err != nil {
		bzc.BzError()
		return
	}
	bzc.SuccessData(accountDTO{AccountID: id, UID: accountUID(id), LoginID: loginID, Token: token, CloudSync: enabled})
}

func accountUID(id int64) string {
	return "QJ-" + strings.ToUpper(strconv.FormatInt(id, 36))
}

func UpdateSyncConsent(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id, ok := accountID(c)
	if !ok {
		bzc.InvalidToken()
		return
	}
	var body syncConsent
	if err := c.ShouldBindJSON(&body); err != nil {
		bzc.ParameterError()
		return
	}
	if body.Enabled == nil {
		bzc.ParameterError()
		return
	}
	enabled := *body.Enabled
	tx, err := datasource.GetMasterDb().BeginTxx(c.Request.Context(), nil)
	if err != nil {
		bzc.BzError()
		return
	}
	if _, err = tx.ExecContext(c.Request.Context(), `UPDATE ledger_account SET cloud_sync_enabled = ?, consent_updated_at = NOW() WHERE account_id = ?`, enabled, id); err != nil {
		_ = tx.Rollback()
		bzc.BzError()
		return
	}
	if !enabled {
		if _, err = tx.ExecContext(c.Request.Context(), `DELETE FROM ledger_record WHERE account_id = ?`, id); err != nil {
			_ = tx.Rollback()
			bzc.BzError()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		bzc.BzError()
		return
	}
	bzc.SuccessData(gin.H{"cloudSyncEnabled": enabled})
}

func GetSyncConsent(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id, ok := accountID(c)
	if !ok {
		bzc.InvalidToken()
		return
	}
	var enabled bool
	if err := datasource.GetMasterDb().Get(&enabled, `SELECT cloud_sync_enabled FROM ledger_account WHERE account_id = ?`, id); err != nil {
		bzc.BzError()
		return
	}
	bzc.SuccessData(gin.H{"cloudSyncEnabled": enabled})
}

func UploadRecords(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id, ok := accountID(c)
	if !ok {
		bzc.InvalidToken()
		return
	}
	var enabled bool
	if err := datasource.GetMasterDb().Get(&enabled, `SELECT cloud_sync_enabled FROM ledger_account WHERE account_id = ?`, id); err != nil || !enabled {
		bzc.Waring("请先明确开启云端同步")
		return
	}
	var body recordBatch
	if err := c.ShouldBindJSON(&body); err != nil || len(body.Records) > 5000 {
		bzc.ParameterError()
		return
	}
	for _, record := range body.Records {
		if strings.TrimSpace(record.ID) == "" || len(record.ID) > 128 || len(record.Date) > 40 || len(record.Merchant) > 255 || len(record.Category) > 100 || len(record.ImportBatchID) > 64 || len(record.ImportFileName) > 255 || record.Amount <= 0 || record.Type != "expense" && record.Type != "income" || record.Source != "manual" && record.Source != "wechat" && record.Source != "alipay" {
			bzc.ErrorMsg("账单数据格式不正确")
			return
		}
	}
	db := datasource.GetMasterDb()
	tx, err := db.BeginTxx(c.Request.Context(), nil)
	if err != nil {
		bzc.BzError()
		return
	}
	for _, record := range body.Records {
		_, err = tx.ExecContext(c.Request.Context(), `INSERT INTO ledger_record (account_id, record_id, record_date, merchant, category, amount, record_type, source, import_batch_id, import_file_name, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW()) ON DUPLICATE KEY UPDATE record_date=VALUES(record_date), merchant=VALUES(merchant), category=VALUES(category), amount=VALUES(amount), record_type=VALUES(record_type), source=VALUES(source), import_batch_id=VALUES(import_batch_id), import_file_name=VALUES(import_file_name), updated_at=NOW()`, id, record.ID, record.Date, record.Merchant, record.Category, record.Amount, record.Type, record.Source, record.ImportBatchID, record.ImportFileName)
		if err != nil {
			_ = tx.Rollback()
			bzc.BzError()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		bzc.BzError()
		return
	}
	bzc.SuccessData(gin.H{"uploaded": len(body.Records)})
}

func ReadRecords(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id, ok := accountID(c)
	if !ok {
		bzc.InvalidToken()
		return
	}
	var enabled bool
	if err := datasource.GetMasterDb().Get(&enabled, `SELECT cloud_sync_enabled FROM ledger_account WHERE account_id = ?`, id); err != nil || !enabled {
		bzc.Waring("云端同步未开启")
		return
	}
	records := make([]ledgerRecord, 0)
	err := datasource.GetMasterDb().Select(&records, `SELECT record_id, record_date, merchant, category, amount, record_type, source, import_batch_id, import_file_name FROM ledger_record WHERE account_id = ? ORDER BY record_date DESC`, id)
	if err != nil {
		bzc.BzError()
		return
	}
	bzc.SuccessData(gin.H{"records": records, "cloudSyncEnabled": true})
}

func DeleteImportBatch(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id, ok := accountID(c)
	if !ok {
		bzc.InvalidToken()
		return
	}
	batchID := strings.TrimSpace(c.Param("batchId"))
	if batchID == "" || len(batchID) > 64 {
		bzc.ParameterError()
		return
	}
	var enabled bool
	if err := datasource.GetMasterDb().Get(&enabled, `SELECT cloud_sync_enabled FROM ledger_account WHERE account_id = ?`, id); err != nil || !enabled {
		bzc.Waring("云端同步未开启")
		return
	}
	result, err := datasource.GetMasterDb().ExecContext(c.Request.Context(), `DELETE FROM ledger_record WHERE account_id = ? AND import_batch_id = ?`, id, batchID)
	if err != nil {
		bzc.BzError()
		return
	}
	deleted, _ := result.RowsAffected()
	bzc.SuccessData(gin.H{"deleted": deleted})
}
