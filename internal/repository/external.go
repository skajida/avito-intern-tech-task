package repository

import "context"

const totalEntries = 1_000

// dummy emulation of the behaviour of an external service containing users database
type ERepository struct {
	database any
}

func NewERepository(users any) *ERepository {
	return &ERepository{database: users}
}

func (*ERepository) Exists(ctx context.Context, id int) bool {
	return id <= totalEntries
}

func (*ERepository) Count(context.Context) int {
	return totalEntries
}
