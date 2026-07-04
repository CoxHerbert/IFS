package dao

import (
	"baize/app/agent/model"
	"baize/app/common/datasource"
	"database/sql"
)

type ChatDao struct{}

var chatDao = &ChatDao{}

func GetChatDao() *ChatDao {
	return chatDao
}

func (dao *ChatDao) InsertSession(userID int64, title, modelName string) int64 {
	result, err := datasource.GetMasterDb().Exec(
		"insert into chat_session(user_id, title, model_name, created_at, updated_at) values(?, ?, ?, now(), now())",
		userID, title, modelName,
	)
	if err != nil {
		panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	return id
}

func (dao *ChatDao) SelectSession(sessionID int64, userID int64) *model.ChatSession {
	session := new(model.ChatSession)
	err := datasource.GetMasterDb().Get(session,
		"select id, user_id, title, model_name, ifnull(summary, '') as summary, created_at, updated_at from chat_session where id = ? and user_id = ?",
		sessionID, userID,
	)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		panic(err)
	}
	return session
}

func (dao *ChatDao) SelectSessions(userID int64) []*model.ChatSession {
	list := make([]*model.ChatSession, 0)
	err := datasource.GetMasterDb().Select(&list,
		"select id, user_id, title, model_name, ifnull(summary, '') as summary, created_at, updated_at from chat_session where user_id = ? order by updated_at desc limit 100",
		userID,
	)
	if err != nil {
		panic(err)
	}
	return list
}

func (dao *ChatDao) UpdateSessionTitle(sessionID int64, userID int64, title string) bool {
	result, err := datasource.GetMasterDb().Exec(
		"update chat_session set title = ?, updated_at = now() where id = ? and user_id = ?",
		title, sessionID, userID,
	)
	if err != nil {
		panic(err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	return affected > 0
}

func (dao *ChatDao) DeleteSession(sessionID int64, userID int64) bool {
	tx := datasource.GetMasterDb().MustBegin()
	result, err := tx.Exec("delete from chat_session where id = ? and user_id = ?", sessionID, userID)
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	if affected == 0 {
		tx.Rollback()
		return false
	}
	if _, err = tx.Exec("delete from chat_message where session_id = ?", sessionID); err != nil {
		tx.Rollback()
		panic(err)
	}
	if err = tx.Commit(); err != nil {
		panic(err)
	}
	return true
}

func (dao *ChatDao) InsertMessage(sessionID int64, role, content, blockResult, modelName string) int64 {
	result, err := datasource.GetMasterDb().Exec(
		"insert into chat_message(session_id, role, content, block_result, model_name, created_at) values(?, ?, ?, ?, ?, now())",
		sessionID, role, content, nullableJSON(blockResult), modelName,
	)
	if err != nil {
		panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	_, err = datasource.GetMasterDb().Exec("update chat_session set updated_at = now() where id = ?", sessionID)
	if err != nil {
		panic(err)
	}
	return id
}

func (dao *ChatDao) SelectMessages(sessionID int64) []*model.ChatMessage {
	list := make([]*model.ChatMessage, 0)
	err := datasource.GetMasterDb().Select(&list,
		"select id, session_id, role, content, ifnull(cast(block_result as char), '') as block_result, ifnull(model_name, '') as model_name, created_at from chat_message where session_id = ? order by id asc",
		sessionID,
	)
	if err != nil {
		panic(err)
	}
	return list
}

func (dao *ChatDao) SelectRecentMessages(sessionID int64, limit int) []*model.ChatMessage {
	list := make([]*model.ChatMessage, 0)
	err := datasource.GetMasterDb().Select(&list,
		"select * from (select id, session_id, role, content, ifnull(cast(block_result as char), '') as block_result, ifnull(model_name, '') as model_name, created_at from chat_message where session_id = ? order by id desc limit ?) t order by id asc",
		sessionID, limit,
	)
	if err != nil {
		panic(err)
	}
	return list
}

func (dao *ChatDao) SelectMemories(userID int64) []*model.ChatMemory {
	list := make([]*model.ChatMemory, 0)
	err := datasource.GetMasterDb().Select(&list,
		"select id, user_id, memory_key, memory_value, memory_type, created_at, updated_at from chat_memory where user_id = ? order by updated_at desc limit 20",
		userID,
	)
	if err != nil {
		panic(err)
	}
	return list
}

func (dao *ChatDao) InsertFormSubmission(sessionID int64, formCode string, formValues string) int64 {
	result, err := datasource.GetMasterDb().Exec(
		"insert into agent_form_submission(session_id, form_code, form_values, created_at) values(?, ?, ?, now())",
		sessionID, formCode, formValues,
	)
	if err != nil {
		panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	return id
}

func nullableJSON(value string) any {
	if value == "" {
		return nil
	}
	return value
}
