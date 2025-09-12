package models

import "context"

type CreateUser struct {
	Name     string `json:"name" firestore:"name" binding:"required"`
	Email    string `json:"email" firestore:"email" binding:"required,email"`
	Password string `json:"password" firestore:"password" binding:"required"`
}

type User struct {
	ID       string `json:"id" firestore:"id" binding:"required"`
	Name     string `json:"name" firestore:"name" binding:"required"`
	Email    string `json:"email" firestore:"email" binding:"required,email"`
	Password string `json:"password" firestore:"password" binding:"required"`
}

type UserRepository interface {
	AddUser(ctx context.Context, u CreateUser) (string, error)
	GetUser(ctx context.Context, id string) (User, error)
	UpdateUser(ctx context.Context, update map[string]interface{}, id string) error
}
