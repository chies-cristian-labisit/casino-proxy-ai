package usecase

import (
	"context"
	"time"

	"github.com/cometagaming/casino-proxy-ai/internal/domain"
)

type UpdateClientNameUseCase struct {
	repo    CustomerRepository
	store   IdempotencyStore
	lockTTL time.Duration
}

func NewUpdateClientNameUseCase(repo CustomerRepository, store IdempotencyStore, lockTTL time.Duration) *UpdateClientNameUseCase {
	return &UpdateClientNameUseCase{repo: repo, store: store, lockTTL: lockTTL}
}

func (uc *UpdateClientNameUseCase) Execute(ctx context.Context, idTx string, newName string) error {
	acquired, err := uc.store.AcquireLock(ctx, idTx, uc.lockTTL)
	if err != nil {
		return err
	}
	if !acquired {
		return nil
	}

	var customer *domain.Customer
	customer, err = uc.repo.GetByCode(ctx, idTx)
	if err != nil {
		_ = uc.store.DeleteKey(ctx, idTx)
		return err
	}

	if err := customer.UpdateName(newName); err != nil {
		_ = uc.store.DeleteKey(ctx, idTx)
		return err
	}

	if err := uc.repo.Save(ctx, customer); err != nil {
		_ = uc.store.DeleteKey(ctx, idTx)
		return err
	}

	return uc.store.SetStatus(ctx, idTx, "COMPLETED")
}
