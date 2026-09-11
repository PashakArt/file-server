package repository

import "context"

type DocRepository struct {
}

func NewDocRepository() *DocRepository {
	return &DocRepository{}
}

func (r *DocRepository) Upload(
	ctx context.Context,
) {
}
