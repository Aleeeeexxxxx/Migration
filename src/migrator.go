package src

import (
	"context"
)

const (
	dbOrigin = "migration_origin"

	dbMigrated1 = "migration_migrated1"
	dbMigrated2 = "migration_migrated2"
	dbMigrated3 = "migration_migrated3"
)

type Migrator struct {
	origin *Service
}

func NewMigrator(cfg DBConfig) (*Migrator, error) {
	cfg1 := cfg.Copy()
	cfg1.Schema = dbMigrated3

	db, err := cfg1.Dial()
	if err != nil {
		return nil, err
	}

	origin, err := NewService(db)
	if err != nil {
		return nil, err
	}
	return &Migrator{
		origin: origin,
	}, nil
}

func (m *Migrator) Create(ctx context.Context, item Model) (*Model, error) {
	return m.origin.Create(ctx, item)
}

func (m *Migrator) Update(ctx context.Context, item Model) error {
	return m.origin.Update(ctx, item)
}

func (m *Migrator) Read(ctx context.Context, id string) (*Model, error) {
	return m.origin.Read(ctx, id)
}

func (m *Migrator) Delete(ctx context.Context, id string) error {
	return m.origin.Delete(ctx, id)
}
