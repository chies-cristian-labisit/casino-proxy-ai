package database

import "github.com/cometagaming/ms-casino-go-v2/internal/usecase"

var _ usecase.CustomerRepository = (*CustomerRepository)(nil)
