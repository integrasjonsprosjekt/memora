package models

type CreateUser struct {
	Name     string `json:"name" firestore:"name" binding:"required"`
	Email    string `json:"email" firestore:"email" binding:"required,email"`
	Password string `json:"password" firestore:"password" binding:"required"`
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name" firestore:"name"`
	Email    string `json:"email" firestore:"email"`
	Password string `json:"password" firestore:"password"`
}

type PatchUser struct {
	Name     string `json:"name" validate:"omitempty"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=12"`
}
