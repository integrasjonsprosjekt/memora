package models

import "memora/internal/utils"

type Card interface {
	GetType() string
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

type MultipleChoiceCard struct {
	ID      string          `json:"id,omitempty" firestore:"-"`
	Type    string          `json:"type" validate:"required" firestore:"type"`
	Options map[string]bool `json:"options" validate:"required" firestore:"options"`
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
	ID      string   `json:"id,omitempty" firestore:"-"`
	Type    string   `json:"type" validate:"required"  firestore:"type"`
	Options []string `json:"options" validate:"required" firestore:"options"`
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
