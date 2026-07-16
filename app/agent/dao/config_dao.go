package dao

import (
	"baize/app/agent/model"
	"baize/app/common/datasource"
	"database/sql"
)

type ConfigDao struct{}

var configDao = &ConfigDao{}

func GetConfigDao() *ConfigDao {
	return configDao
}

func (dao *ConfigDao) SelectAll() map[string]string {
	rows := make([]*model.RuntimeConfig, 0)
	err := datasource.GetMasterDb().Select(&rows, "select config_key, config_value from agent_runtime_config")
	if err == sql.ErrNoRows {
		return map[string]string{}
	}
	if err != nil {
		panic(err)
	}
	result := make(map[string]string, len(rows))
	for _, row := range rows {
		result[row.ConfigKey] = row.ConfigValue
	}
	return result
}

func (dao *ConfigDao) Upsert(values map[string]string) {
	tx := datasource.GetMasterDb().MustBegin()
	for key, value := range values {
		_, err := tx.Exec(
			"insert into agent_runtime_config(config_key, config_value, update_time) values(?, ?, now()) on duplicate key update config_value = values(config_value), update_time = now()",
			key, value,
		)
		if err != nil {
			_ = tx.Rollback()
			panic(err)
		}
	}
	if err := tx.Commit(); err != nil {
		panic(err)
	}
}
