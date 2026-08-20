package user

import "context"

type UserService interface {
	CreateUser(ctx context.Context, user *User) error
}
