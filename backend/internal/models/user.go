package models

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
