package subnet

import "context"

type Repository interface {
	Create(ctx context.Context, subnet *Subnet) (string, error)
	FindAll(ctx context.Context) ([]Subnet, error)
	FindOne(ctx context.Context, id string) (Subnet, error)
	Update(ctx context.Context, subnet *Subnet) error
	Delete(ctx context.Context, id string) error
}
