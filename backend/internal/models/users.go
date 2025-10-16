package models

type CreateUser struct {
	Name  string `json:"name" firestore:"name" validate:"required"`
	Email string `json:"email" firestore:"email" validate:"required,email"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name" firestore:"name"`
	Email string `json:"email" firestore:"email"`
}
