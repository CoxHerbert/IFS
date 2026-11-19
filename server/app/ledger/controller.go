package ledger

import (
	"baize/app/common/baize/baizeContext"
	"baize/app/common/datasource"
	"baize/app/system/service/loginService/loginServiceImpl"
	"baize/app/utils/bCryptPasswordEncoder"
	"errors"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

type credentials struct {
	LoginID     string `json:"loginId" binding:"required"`
	Password    string `json:"password" binding:"required"`
	InviterUID  string `json:"inviterUid"`
	CaptchaID   string `json:"captchaId"`
	CaptchaCode string `json:"captchaCode"`
}

type deleteAccountRequest struct {
	Password string `json:"password" binding:"required"`
}

type updatePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

type accountDTO struct {
	AccountID  int64  `json:"accountId"`
	UID        string `json:"uid"`
	LoginID    string `json:"loginId"`
	Token      string `json:"token"`
	CloudSync  bool   `json:"cloudSyncEnabled"`
	InviterUID string `json:"inviterUid"`
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

type importBatch struct {
	BatchID        string `json:"batchId" db:"batch_id"`
	Source         string `json:"source" db:"source"`
	FileName       string `json:"fileName" db:"file_name"`
	ImportedAt     string `json:"importedAt" db:"imported_at"`
	ImportedCount  int    `json:"importedCount" db:"imported_count"`
	DuplicateCount int    `json:"duplicateCount" db:"duplicate_count"`
}

type importBatchList struct {
	Batches []importBatch `json:"batches"`
}

type syncConsent struct {
	Enabled *bool `json:"enabled"`
}

type accountProfile struct {
	Nickname       string `json:"nickname" db:"nickname"`
	AvatarKey      string `json:"avatarKey" db:"avatar_key"`
	AvatarFrameKey string `json:"avatarFrameKey" db:"avatar_frame_key"`
}

var phonePattern = regexp.MustCompile(`^\+?[0-9]{7,20}$`)

func GetRegisterCaptcha(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	bzc.SuccessData(loginServiceImpl.GenerateCode())
}

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
	if strings.TrimSpace(body.CaptchaID) == "" || !loginServiceImpl.VerityCaptcha(body.CaptchaID, strings.TrimSpace(body.CaptchaCode)) {
		bzc.ErrorMsg("验证码错误或已过期")
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
	inviterUID := strings.ToUpper(strings.TrimSpace(body.InviterUID))
	if inviterUID != "" {
		if len(inviterUID) < 4 || len(inviterUID) > 32 || !strings.HasPrefix(inviterUID, "QJ-") {
			bzc.ErrorMsg("邀请人信息无效")
			return
		}
		var inviterExists int
		if err := datasource.GetMasterDb().Get(&inviterExists, `SELECT COUNT(*) FROM ledger_account WHERE CONCAT('QJ-', UPPER(CONV(account_id, 10, 36))) = ?`, inviterUID); err != nil || inviterExists == 0 {
			bzc.ErrorMsg("邀请人账号不存在或已注销")
			return
		}
	}
	result, err := datasource.GetMasterDb().Exec(`INSERT INTO ledger_account (login_id, password_hash, cloud_sync_enabled, inviter_uid, create_time) VALUES (?, ?, 0, ?, NOW())`, loginID, hash, inviterUID)
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
	var account struct {
		Enabled    bool   `db:"cloud_sync_enabled"`
		InviterUID string `db:"inviter_uid"`
	}
	if err := datasource.GetMasterDb().Get(&account, `SELECT cloud_sync_enabled, inviter_uid FROM ledger_account WHERE account_id = ?`, id); err != nil {
		bzc.BzError()
		return
	}
	bzc.SuccessData(accountDTO{AccountID: id, UID: accountUID(id), LoginID: loginID, Token: token, CloudSync: account.Enabled, InviterUID: account.InviterUID})
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
		if _, err = tx.ExecContext(c.Request.Context(), `DELETE FROM ledger_import_batch WHERE account_id = ?`, id); err != nil {
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

func GetProfile(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id, ok := accountID(c)
	if !ok {
		bzc.InvalidToken()
		return
	}
	var profile accountProfile
	if err := datasource.GetMasterDb().Get(&profile, `SELECT nickname, avatar_key, avatar_frame_key FROM ledger_account WHERE account_id = ?`, id); err != nil {
		bzc.BzError()
		return
	}
	bzc.SuccessData(gin.H{"profile": profile})
}

func UpdateProfile(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id, ok := accountID(c)
	if !ok {
		bzc.InvalidToken()
		return
	}
	var profile accountProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		bzc.ParameterError()
		return
	}
	profile.AvatarKey = strings.TrimSpace(profile.AvatarKey)
	profile.AvatarFrameKey = strings.TrimSpace(profile.AvatarFrameKey)
	profile.Nickname = strings.TrimSpace(profile.Nickname)
	if utf8.RuneCountInString(profile.Nickname) > 30 || len(profile.AvatarKey) > 32 || !validProfileAvatarKey(profile.AvatarKey) || len(profile.AvatarFrameKey) > 32 || !validAvatarFrameKey(profile.AvatarFrameKey) {
		bzc.ErrorMsg("个人资料格式不正确")
		return
	}
	if profile.AvatarFrameKey != "" {
		var invitedCount int
		if err := datasource.GetMasterDb().Get(&invitedCount, `SELECT COUNT(*) FROM ledger_account WHERE inviter_uid = ?`, accountUID(id)); err != nil || invitedCount < 3 {
			bzc.ErrorMsg("邀请 3 位好友注册后即可使用星光头像框")
			return
		}
	}
	if _, err := datasource.GetMasterDb().ExecContext(c.Request.Context(), `UPDATE ledger_account SET nickname = ?, avatar_key = ?, avatar_frame_key = ? WHERE account_id = ?`, profile.Nickname, profile.AvatarKey, profile.AvatarFrameKey, id); err != nil {
		bzc.BzError()
		return
	}
	bzc.SuccessData(gin.H{"profile": profile})
}

func GetReferrals(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id, ok := accountID(c)
	if !ok {
		bzc.InvalidToken()
		return
	}
	var invitedCount int
	if err := datasource.GetMasterDb().Get(&invitedCount, `SELECT COUNT(*) FROM ledger_account WHERE inviter_uid = ?`, accountUID(id)); err != nil {
		bzc.BzError()
		return
	}
	bzc.SuccessData(gin.H{"invitedCount": invitedCount, "inviteStars": invitedCount, "avatarFrameUnlocked": invitedCount >= 3, "nextFrameAt": 3})
}

func UpdatePassword(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id, ok := accountID(c)
	if !ok {
		bzc.InvalidToken()
		return
	}
	var body updatePasswordRequest
	if err := c.ShouldBindJSON(&body); err != nil || len(body.OldPassword) < 8 || len(body.OldPassword) > 72 || len(body.NewPassword) < 8 || len(body.NewPassword) > 72 {
		bzc.ErrorMsg("密码长度需为 8 到 72 个字符")
		return
	}
	if body.OldPassword == body.NewPassword {
		bzc.ErrorMsg("新密码不能与当前密码相同")
		return
	}
	var passwordHash string
	if err := datasource.GetMasterDb().Get(&passwordHash, `SELECT password_hash FROM ledger_account WHERE account_id = ?`, id); err != nil || !bCryptPasswordEncoder.CheckPasswordHash(body.OldPassword, passwordHash) {
		bzc.ErrorMsg("当前密码错误")
		return
	}
	newHash := bCryptPasswordEncoder.HashPassword(body.NewPassword)
	if _, err := datasource.GetMasterDb().ExecContext(c.Request.Context(), `UPDATE ledger_account SET password_hash = ? WHERE account_id = ?`, newHash, id); err != nil {
		bzc.BzError()
		return
	}
	bzc.SuccessData(gin.H{"updated": true})
}

func DeleteAccount(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id, ok := accountID(c)
	if !ok {
		bzc.InvalidToken()
		return
	}
	var body deleteAccountRequest
	if err := c.ShouldBindJSON(&body); err != nil || len(body.Password) < 8 || len(body.Password) > 72 {
		bzc.ErrorMsg("请输入正确的账号密码")
		return
	}
	var passwordHash string
	if err := datasource.GetMasterDb().Get(&passwordHash, `SELECT password_hash FROM ledger_account WHERE account_id = ?`, id); err != nil || !bCryptPasswordEncoder.CheckPasswordHash(body.Password, passwordHash) {
		bzc.ErrorMsg("账号密码错误，未注销账号")
		return
	}
	tx, err := datasource.GetMasterDb().BeginTxx(c.Request.Context(), nil)
	if err != nil {
		bzc.BzError()
		return
	}
	if _, err = tx.ExecContext(c.Request.Context(), `DELETE FROM ledger_record WHERE account_id = ?`, id); err == nil {
		_, err = tx.ExecContext(c.Request.Context(), `DELETE FROM ledger_import_batch WHERE account_id = ?`, id)
	}
	if err == nil {
		_, err = tx.ExecContext(c.Request.Context(), `DELETE FROM ledger_account WHERE account_id = ?`, id)
	}
	if err != nil {
		_ = tx.Rollback()
		bzc.BzError()
		return
	}
	if err = tx.Commit(); err != nil {
		bzc.BzError()
		return
	}
	bzc.SuccessData(gin.H{"deleted": true})
}

func validProfileAvatarKey(key string) bool {
	switch key {
	case "", "cat", "bear", "bunny", "fox", "panda", "frog", "koala", "tiger", "lion", "penguin", "whale", "unicorn", "hamster", "deer", "owl", "dog", "chick":
		return true
	default:
		return false
	}
}

func validAvatarFrameKey(key string) bool {
	return key == "" || key == "sparkle"
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

func canWriteCloudRecords(bzc *baizeContext.BaiZeContext, accountID int64) bool {
	var enabled bool
	if err := datasource.GetMasterDb().Get(&enabled, `SELECT cloud_sync_enabled FROM ledger_account WHERE account_id = ?`, accountID); err != nil || !enabled {
		bzc.Waring("请先明确开启云端同步")
		return false
	}
	return true
}

func validateLedgerRecord(record ledgerRecord) error {
	record.ID = strings.TrimSpace(record.ID)
	if record.ID == "" || len(record.ID) > 128 || strings.TrimSpace(record.Date) == "" || len(record.Date) > 40 || strings.TrimSpace(record.Merchant) == "" || len(record.Merchant) > 255 || strings.TrimSpace(record.Category) == "" || len(record.Category) > 100 || len(record.ImportBatchID) > 64 || len(record.ImportFileName) > 255 || record.Amount <= 0 || record.Type != "expense" && record.Type != "income" || record.Source != "manual" && record.Source != "wechat" && record.Source != "alipay" {
		return errors.New("账单数据格式不正确")
	}
	return nil
}

func CreateRecord(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id, ok := accountID(c)
	if !ok {
		bzc.InvalidToken()
		return
	}
	if !canWriteCloudRecords(bzc, id) {
		return
	}
	var record ledgerRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		bzc.ParameterError()
		return
	}
	if err := validateLedgerRecord(record); err != nil {
		bzc.ErrorMsg(err.Error())
		return
	}
	_, err := datasource.GetMasterDb().ExecContext(c.Request.Context(), `INSERT INTO ledger_record (account_id, record_id, record_date, merchant, category, amount, record_type, source, import_batch_id, import_file_name, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`, id, record.ID, record.Date, strings.TrimSpace(record.Merchant), strings.TrimSpace(record.Category), record.Amount, record.Type, record.Source, record.ImportBatchID, record.ImportFileName)
	if err != nil {
		bzc.Waring("账单已存在或暂时无法保存")
		return
	}
	bzc.SuccessData(gin.H{"record": record})
}

func UpdateRecord(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id, ok := accountID(c)
	if !ok {
		bzc.InvalidToken()
		return
	}
	if !canWriteCloudRecords(bzc, id) {
		return
	}
	recordID := strings.TrimSpace(c.Param("recordId"))
	if recordID == "" || len(recordID) > 128 {
		bzc.ParameterError()
		return
	}
	var record ledgerRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		bzc.ParameterError()
		return
	}
	if record.ID != recordID {
		bzc.ErrorMsg("账单编号不匹配")
		return
	}
	if err := validateLedgerRecord(record); err != nil {
		bzc.ErrorMsg(err.Error())
		return
	}
	var exists int
	if err := datasource.GetMasterDb().Get(&exists, `SELECT COUNT(*) FROM ledger_record WHERE account_id = ? AND record_id = ?`, id, recordID); err != nil {
		bzc.BzError()
		return
	}
	if exists == 0 {
		bzc.ErrorMsg("账单不存在")
		return
	}
	_, err := datasource.GetMasterDb().ExecContext(c.Request.Context(), `UPDATE ledger_record SET record_date = ?, merchant = ?, category = ?, amount = ?, record_type = ?, source = ?, import_batch_id = ?, import_file_name = ?, updated_at = NOW() WHERE account_id = ? AND record_id = ?`, record.Date, strings.TrimSpace(record.Merchant), strings.TrimSpace(record.Category), record.Amount, record.Type, record.Source, record.ImportBatchID, record.ImportFileName, id, recordID)
	if err != nil {
		bzc.BzError()
		return
	}
	bzc.SuccessData(gin.H{"record": record})
}

func DeleteRecord(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id, ok := accountID(c)
	if !ok {
		bzc.InvalidToken()
		return
	}
	if !canWriteCloudRecords(bzc, id) {
		return
	}
	recordID := strings.TrimSpace(c.Param("recordId"))
	if recordID == "" || len(recordID) > 128 {
		bzc.ParameterError()
		return
	}
	result, err := datasource.GetMasterDb().ExecContext(c.Request.Context(), `DELETE FROM ledger_record WHERE account_id = ? AND record_id = ?`, id, recordID)
	if err != nil {
		bzc.BzError()
		return
	}
	deleted, _ := result.RowsAffected()
	bzc.SuccessData(gin.H{"deleted": deleted})
}

func ReadImportBatches(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id, ok := accountID(c)
	if !ok {
		bzc.InvalidToken()
		return
	}
	if !canWriteCloudRecords(bzc, id) {
		return
	}
	batches := make([]importBatch, 0)
	err := datasource.GetMasterDb().Select(&batches, `SELECT batch_id, source, file_name, imported_at, imported_count, duplicate_count FROM ledger_import_batch WHERE account_id = ? ORDER BY imported_at DESC`, id)
	if err != nil {
		bzc.BzError()
		return
	}
	bzc.SuccessData(gin.H{"batches": batches})
}

func validateImportBatch(batch importBatch) error {
	if strings.TrimSpace(batch.BatchID) == "" || len(batch.BatchID) > 64 || batch.Source != "wechat" && batch.Source != "alipay" || strings.TrimSpace(batch.FileName) == "" || len(batch.FileName) > 255 || strings.TrimSpace(batch.ImportedAt) == "" || len(batch.ImportedAt) > 40 || batch.ImportedCount < 0 || batch.DuplicateCount < 0 {
		return errors.New("导入批次数据格式不正确")
	}
	return nil
}

func upsertImportBatches(c *gin.Context, bzc *baizeContext.BaiZeContext, accountID int64, batches []importBatch) {
	if len(batches) > 1000 {
		bzc.ParameterError()
		return
	}
	for _, batch := range batches {
		if err := validateImportBatch(batch); err != nil {
			bzc.ErrorMsg(err.Error())
			return
		}
	}
	tx, err := datasource.GetMasterDb().BeginTxx(c.Request.Context(), nil)
	if err != nil {
		bzc.BzError()
		return
	}
	for _, batch := range batches {
		_, err = tx.ExecContext(c.Request.Context(), `INSERT INTO ledger_import_batch (account_id, batch_id, source, file_name, imported_at, imported_count, duplicate_count, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, NOW()) ON DUPLICATE KEY UPDATE source=VALUES(source), file_name=VALUES(file_name), imported_at=VALUES(imported_at), imported_count=VALUES(imported_count), duplicate_count=VALUES(duplicate_count), updated_at=NOW()`, accountID, strings.TrimSpace(batch.BatchID), batch.Source, strings.TrimSpace(batch.FileName), batch.ImportedAt, batch.ImportedCount, batch.DuplicateCount)
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
	bzc.SuccessData(gin.H{"saved": len(batches)})
}

func SaveImportBatch(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id, ok := accountID(c)
	if !ok {
		bzc.InvalidToken()
		return
	}
	if !canWriteCloudRecords(bzc, id) {
		return
	}
	var batch importBatch
	if err := c.ShouldBindJSON(&batch); err != nil {
		bzc.ParameterError()
		return
	}
	upsertImportBatches(c, bzc, id, []importBatch{batch})
}

func UploadImportBatches(c *gin.Context) {
	bzc := baizeContext.NewBaiZeContext(c)
	id, ok := accountID(c)
	if !ok {
		bzc.InvalidToken()
		return
	}
	if !canWriteCloudRecords(bzc, id) {
		return
	}
	var body importBatchList
	if err := c.ShouldBindJSON(&body); err != nil {
		bzc.ParameterError()
		return
	}
	upsertImportBatches(c, bzc, id, body.Batches)
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
	if !canWriteCloudRecords(bzc, id) {
		return
	}
	tx, err := datasource.GetMasterDb().BeginTxx(c.Request.Context(), nil)
	if err != nil {
		bzc.BzError()
		return
	}
	result, err := tx.ExecContext(c.Request.Context(), `DELETE FROM ledger_record WHERE account_id = ? AND import_batch_id = ?`, id, batchID)
	if err != nil {
		_ = tx.Rollback()
		bzc.BzError()
		return
	}
	deleted, _ := result.RowsAffected()
	if _, err = tx.ExecContext(c.Request.Context(), `DELETE FROM ledger_import_batch WHERE account_id = ? AND batch_id = ?`, id, batchID); err != nil {
		_ = tx.Rollback()
		bzc.BzError()
		return
	}
	if err = tx.Commit(); err != nil {
		bzc.BzError()
		return
	}
	bzc.SuccessData(gin.H{"deleted": deleted})
}
