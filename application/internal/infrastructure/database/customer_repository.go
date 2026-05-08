package database

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/cometagaming/casino-proxy-ai/internal/domain"
	"github.com/cometagaming/casino-proxy-ai/internal/usecase"
)

var _ usecase.CustomerRepository = (*CustomerRepository)(nil)

type customerRecord struct {
	ID   uint   `gorm:"primarykey;autoIncrement"`
	Code string `gorm:"uniqueIndex;not null"`
	Name string `gorm:"not null"`
}

func (r customerRecord) toDomain() *domain.Customer {
	return &domain.Customer{ID: r.ID, Code: r.Code, Name: r.Name}
}

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) GetByCode(ctx context.Context, code string) (*domain.Customer, error) {
	var rec customerRecord
	result := r.db.WithContext(ctx).Where("code = ?", code).First(&rec)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, domain.ErrCustomerNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return rec.toDomain(), nil
}

func (r *CustomerRepository) Save(ctx context.Context, customer *domain.Customer) error {
	rec := customerRecord{ID: customer.ID, Code: customer.Code, Name: customer.Name}
	return r.db.WithContext(ctx).Save(&rec).Error
}

// Migrate creates or updates the database schema for all repository models.
// Called from cmd/api/main.go on startup.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&customerRecord{})
}
