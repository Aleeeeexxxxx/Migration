//go:build mysql

package src

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func createMigrators(n int, rq *require.Assertions, cfg DBConfig) []*Migrator {
	var ret []*Migrator

	for i := 0; i < n; i++ {
		m, err := NewMigrator(cfg)
		rq.NoError(err)
		ret = append(ret, m)
	}
	return ret
}

func TestMigration_local(t *testing.T) {
	rq := require.New(t)
	SetDefaultLoggerLevel(zapcore.ErrorLevel)

	cfg, loadErr := LoadDBConfigFromFile("../mysql.json")
	rq.NoError(loadErr)

	db, err := cfg.Dial()
	rq.NoError(err)

	m := createMigrators(3, rq, cfg)

	RunCurdConcurrency(rq, m, 10, 1, C, U, R)

	ValidateConsistence(rq, db)
}

func TestMigration_server(t *testing.T) {

}
