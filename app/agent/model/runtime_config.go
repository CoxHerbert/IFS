package model

type RuntimeConfig struct {
	ConfigKey   string `db:"config_key"`
	ConfigValue string `db:"config_value"`
}
