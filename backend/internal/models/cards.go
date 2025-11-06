package models

import (
	"memora/internal/utils"
	"time"
)

// Card is an interface that all card types implement.
type Card interface {
	// GetType returns the type of the card as a string.
	GetType() string

	// SetID sets the ID of the card.
	SetID(id string)
}

// This tells Swagger that the response can be one of these types
type AnyCard struct {
	// @swagger:oneOf
	FrontBackCard      FrontBackCard
	MultipleChoiceCard MultipleChoiceCard
	OrderedCard        OrderedCard
	BlanksCard         BlanksCard
}

type AnyCardWithPaging struct {
	Cards   []AnyCard `json:"cards"`
	HasMore bool      `json:"has_more"`
	Cursor  string    `json:"cursor,omitempty"`
}

type MultipleChoiceCard struct {
	ID       string          `json:"id,omitempty" validate:"omitempty" firestore:"-"`
	Type     string          `json:"type" validate:"required" firestore:"type"`
	Question string          `json:"question" validate:"required" firestore:"question"`
	Options  map[string]bool `json:"options" validate:"required" firestore:"options"`
}

func (m MultipleChoiceCard) GetType() string  { return utils.MULTIPLE_CHOICE_CARD }
func (m *MultipleChoiceCard) SetID(id string) { m.ID = id }

type FrontBackCard struct {
	ID    string `json:"id,omitempty" firestore:"-"`
	Type  string `json:"type" validate:"required" firestore:"type"`
	Front string `json:"front" validate:"required" firestore:"front"`
	Back  string `json:"back" validate:"required" firestore:"back"`
}

func (f FrontBackCard) GetType() string  { return utils.FRONT_BACK_CARD }
func (f *FrontBackCard) SetID(id string) { f.ID = id }

type OrderedCard struct {
	ID       string   `json:"id,omitempty" firestore:"-"`
	Type     string   `json:"type" validate:"required"  firestore:"type"`
	Question string   `json:"question" validate:"required" firestore:"question"`
	Options  []string `json:"options" validate:"required" firestore:"options"`
}

func (o OrderedCard) GetType() string  { return utils.ORDERED_CARD }
func (o *OrderedCard) SetID(id string) { o.ID = id }

type BlanksCard struct {
	ID       string   `json:"id,omitempty" firestore:"-"`
	Type     string   `json:"type" validate:"required" firestore:"type"`
	Question string   `json:"question" validate:"required" firestore:"question"`
	Answers  []string `json:"answers" validate:"required" firestore:"answers"`
}

func (b BlanksCard) GetType() string  { return utils.BLANKS_CARD }
func (b *BlanksCard) SetID(id string) { b.ID = id }

type CardType struct {
	Type string `json:"type"`
}

type CardsResponse struct {
	Cards   []Card `json:"cards"`
	HasMore bool   `json:"has_more"`
}

type CardRating struct {
	Rating string `json:"rating" validate:"oneof=again hard good easy"`
}

type CardProgress struct {
	EaseFactor   int       `firestore:"ease_factor" json:"ease_factor"`
	Interval     float64   `firestore:"interval" json:"interval"`
	Due          time.Time `firestore:"due" json:"due"`
	Reps         int       `firestore:"reps" json:"reps"`
	Lapses       int       `firestore:"lapses" json:"lapses"`
	LastReviewed time.Time `firestore:"last_reviewed_at" json:"last_reviewed_at"`
}
