package integration_test

import (
	"context"
	"testing"

	"github.com/cometagaming/ms-casino-go-v2/internal/domain"
)

func truncateCustomers(t *testing.T) {
	t.Helper()
	if err := sharedDB.Exec("TRUNCATE TABLE customer_records RESTART IDENTITY CASCADE").Error; err != nil {
		t.Fatalf("truncate customer_records: %v", err)
	}
}

func TestSave_PersistsNewCustomer(t *testing.T) {
	truncateCustomers(t)
	ctx := context.Background()

	if err := sharedRepo.Save(ctx, &domain.Customer{Code: "BR100000001", Name: "Alice"}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := sharedRepo.GetByCode(ctx, "BR100000001")
	if err != nil {
		t.Fatalf("GetByCode: %v", err)
	}
	if got.Name != "Alice" {
		t.Errorf("Name: want %q, got %q", "Alice", got.Name)
	}
}

func TestGetByCode_ReturnsErrCustomerNotFound(t *testing.T) {
	truncateCustomers(t)

	_, err := sharedRepo.GetByCode(context.Background(), "DOES_NOT_EXIST")
	if err != domain.ErrCustomerNotFound {
		t.Errorf("error: want domain.ErrCustomerNotFound, got %v", err)
	}
}

func TestSave_UpdatesExistingCustomer(t *testing.T) {
	truncateCustomers(t)
	ctx := context.Background()

	if err := sharedRepo.Save(ctx, &domain.Customer{Code: "BR200000002", Name: "First Name"}); err != nil {
		t.Fatalf("first Save: %v", err)
	}

	// Fetch the auto-assigned ID so the second Save is an UPDATE, not a conflicting INSERT.
	existing, err := sharedRepo.GetByCode(ctx, "BR200000002")
	if err != nil {
		t.Fatalf("GetByCode: %v", err)
	}

	if err := sharedRepo.Save(ctx, &domain.Customer{ID: existing.ID, Code: "BR200000002", Name: "Second Name"}); err != nil {
		t.Fatalf("second Save: %v", err)
	}

	got, err := sharedRepo.GetByCode(ctx, "BR200000002")
	if err != nil {
		t.Fatalf("GetByCode after update: %v", err)
	}
	if got.Name != "Second Name" {
		t.Errorf("Name: want %q, got %q", "Second Name", got.Name)
	}
}

func TestGetByCode_ReturnsAllFields(t *testing.T) {
	truncateCustomers(t)
	ctx := context.Background()

	if err := sharedRepo.Save(ctx, &domain.Customer{Code: "BR300000003", Name: "Bob"}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := sharedRepo.GetByCode(ctx, "BR300000003")
	if err != nil {
		t.Fatalf("GetByCode: %v", err)
	}
	if got.Code != "BR300000003" {
		t.Errorf("Code: want %q, got %q", "BR300000003", got.Code)
	}
	if got.Name != "Bob" {
		t.Errorf("Name: want %q, got %q", "Bob", got.Name)
	}
	if got.ID == 0 {
		t.Error("ID: want > 0, got 0 — auto-increment did not fire")
	}
}
