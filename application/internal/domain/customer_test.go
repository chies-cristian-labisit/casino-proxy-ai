package domain

import "testing"

func TestUpdateName_Empty(t *testing.T) {
	c := &Customer{ID: 1, Code: "BR123", Name: "Old Name"}
	err := c.UpdateName("")
	if err != ErrInvalidName {
		t.Errorf("expected ErrInvalidName, got %v", err)
	}
	if c.Name != "Old Name" {
		t.Error("Name must not change on error")
	}
}

func TestUpdateName_Valid(t *testing.T) {
	c := &Customer{ID: 1, Code: "BR123", Name: "Old Name"}
	err := c.UpdateName("João Silva")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Name != "João Silva" {
		t.Errorf("expected Name=João Silva, got %s", c.Name)
	}
}

func TestSentinelErrors_NonEmpty(t *testing.T) {
	errors := []error{ErrCustomerNotFound, ErrInvalidCode, ErrInvalidName}
	for _, err := range errors {
		if err.Error() == "" {
			t.Errorf("sentinel error has empty message: %T", err)
		}
	}
}
