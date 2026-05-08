package domain

import "errors"

type Customer struct {
	ID   uint
	Code string // idTx — national registration code, business key
	Name string
}

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrInvalidCode      = errors.New("customer code cannot be empty")
	ErrInvalidName      = errors.New("customer name cannot be empty")
)

func (c *Customer) UpdateName(name string) error {
	if name == "" {
		return ErrInvalidName
	}
	c.Name = name
	return nil
}
