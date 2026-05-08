package database

import "github.com/cometagaming/casino-proxy-ai/internal/usecase"

var _ usecase.CustomerRepository = (*CustomerRepository)(nil)
