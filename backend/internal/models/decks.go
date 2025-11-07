package models

import "encoding/json"

type CreateDeck struct {
	Title        string   `json:"title" validate:"required" firestore:"title"`
	OwnerID      string   `json:"owner_id" validate:"required" firestore:"owner_id"`
	SharedEmails []string `json:"shared_emails" validate:"omitempty,dive,email" firestore:"shared_emails"`
}

type DeckResponse struct {
	ID           string   `json:"id" firestore:"-"`
	Title        string   `json:"title,omitempty" firestore:"title"`
	OwnerID      string   `json:"owner_id,omitempty" firestore:"owner_id"`
	SharedEmails []string `json:"shared_emails,omitempty" firestore:"shared_emails"`
	Cards        []Card   `json:"cards" firestore:"cards"`
}

type Deck struct {
	OwnerID      string   `json:"owner_id" firestore:"owner_id"`
	Title        string   `json:"title" firestore:"title"`
	SharedEmails []string `json:"shared_emails" firestore:"shared_emails"`
}

type UpdateDeck struct {
	Title string `json:"title,omitempty" firestore:"title"`
}

type UpdateDeckEmails struct {
	Opp    string   `json:"opp" validate:"required,oneof=add remove"`
	Emails []string `json:"shared_emails" firestore:"shared_emails" validate:"required"`
}

type DisplayDeck struct {
	ID      string `json:"id" firestore:"-"`
	Title   string `json:"title,omitempty" firestore:"title"`
	OwnerID string `json:"owner_id,omitempty" firestore:"owner_id"`
}

func (d DeckResponse) MarshalJSON() ([]byte, error) {
	type Alias DeckResponse
	if d.SharedEmails == nil {
		d.SharedEmails = []string{}
	}
	if d.Cards == nil {
		d.Cards = []Card{}
	}
	return json.Marshal(Alias(d))
}
