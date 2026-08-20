package user

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
}
