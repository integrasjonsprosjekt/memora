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

type PatchUser struct {
	Name  string `json:"name" validate:"omitempty"`
	Email string `json:"email" validate:"omitempty,email"`
}

type UserDecks struct {
	OwnedDecks  []DisplayDeck `json:"owned_decks"`
	SharedDecks []DisplayDeck `json:"shared_decks"`
}
