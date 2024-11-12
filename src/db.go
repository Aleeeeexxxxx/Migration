package src

import (
	"encoding/json"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Schema   string `json:"schema"`
}

func LoadDBConfigFromFile(path string) (DBConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return DBConfig{}, err
	}
	defer func() { _ = f.Close() }()

	var config DBConfig
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return DBConfig{}, err
	}
	return config, nil
}

func (cfg DBConfig) dsn() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		cfg.Username,
		cfg.Password,
		cfg.IP, cfg.Port,
		cfg.Schema,
	)
}

func (cfg DBConfig) Dial() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.dsn()), &gorm.Config{
		Logger: NewCustomLogger(),
	})
	if err != nil {
		return nil, fmt.Errorf("fail to open db: %w", err)
	}
	return db, nil
}

func (cfg DBConfig) Copy() DBConfig {
	return DBConfig{
		Username: cfg.Username,
		Password: cfg.Password,
		IP:       cfg.IP,
		Port:     cfg.Port,
	}
}
