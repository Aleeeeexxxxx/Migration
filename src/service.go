package src

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Model struct {
	ID        string `json:"id,omitempty" gorm:"primaryKey"`
	Msg       string `json:"msg" gorm:"size:255"`
	UpdatedAt int64  `json:"updated,omitempty" gorm:"autoUpdateTime:false"`
}

type Service struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewService(db *gorm.DB) (*Service, error) {
	if err := db.AutoMigrate(&Model{}); err != nil {
		return nil, err
	}
	return &Service{
		db:     db,
		logger: GetDefaultLogger().With(zap.String(ZModule, "service")),
	}, nil
}

func (svc *Service) Create(ctx context.Context, item Model) (*Model, error) {
	item.ID = uuid.New().String()
	item.UpdatedAt = time.Now().Unix()

	logger := WithContext(ctx, svc.logger)
	logger.Debug(
		"creating new model",
		zap.String("id", item.ID),
		zap.String("msg", item.Msg),
	)

	if result := svc.db.Create(item); result.Error != nil {
		logger.Error("failed to create model", zap.Error(result.Error))
		return nil, result.Error
	}

	logger.Info(
		"new model created",
		zap.String("id", item.ID),
	)
	return &item, nil
}

func (svc *Service) Update(ctx context.Context, item Model) error {
	item.UpdatedAt = time.Now().Unix()

	logger := WithContext(ctx, svc.logger)
	logger.Debug(
		"updating model",
		zap.String("id", item.ID),
		zap.String("msg", item.Msg),
	)

	if result := svc.db.Model(&Model{}).
		Where("id = ?", item.ID).
		Update("msg", item.Msg); result.Error != nil {
		logger.Error("failed to update model", zap.Error(result.Error))
		return result.Error
	}

	logger.Info(
		"updated model",
		zap.String("id", item.ID),
	)
	return nil
}

func (svc *Service) Read(ctx context.Context, id string) (*Model, error) {
	logger := WithContext(ctx, svc.logger)
	logger.Debug(
		"read model",
		zap.String("id", id),
	)

	var ret Model
	if result := svc.db.Where("id = ?", id).First(&ret); result.Error != nil {
		logger.Error("failed to read model", zap.Error(result.Error))
		return nil, result.Error
	}
	return nil, nil
}

func (svc *Service) Delete(ctx context.Context, id string) error {
	logger := WithContext(ctx, svc.logger)
	logger.Debug(
		"delete model",
		zap.String("id", id),
	)

	if result := svc.db.Delete(&Model{ID: id}); result.Error != nil {
		logger.Error("failed to delete model", zap.Error(result.Error))
		return result.Error
	}
	logger.Info(
		"model deleted",
		zap.String("id", id),
	)
	return nil
}
