package models

import "cloud.google.com/go/firestore"

type CreateDeck struct {
	Title         string   `json:"title"`
	Owner_id      string   `json:"owner_id"`
}

type DeckResponse struct {
	ID            string   `json:"id"`
	Owner_id      string   `json:"owner_id"`
	Shared_emails []string `json:"shared_emails" firestore:"shared_emails"`
	Cards         []Card   `json:"-" firestore:"cards"`
}

type Deck struct {
	ID            string                   `json:"id"`
	Owner_id      string                   `json:"owner_id"`
	Shared_emails []string                 `json:"shared_emails" firestore:"shared_emails"`
	Cards         []*firestore.DocumentRef `json:"-" firestore:"cards"`
}
