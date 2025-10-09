export type CardType = 'front_back' | 'blanks' | 'multiple_choice' | 'ordered';

interface BaseCard {
  id: string;
  type: CardType;
}

export interface FrontBackCard extends BaseCard {
  type: 'front_back';
  front: string;
  back: string;
}

export interface FillBlanksCard extends BaseCard {
  type: 'blanks';
  question: string;
  answers: string[];
}

export interface MultipleChoiceCard extends BaseCard {
  type: 'multiple_choice';
  question: string;
  options: { [option: string]: boolean };
}

export interface OrderedCard extends BaseCard {
  type: 'ordered';
  question: string;
  options: string[];
}

export type CardComponentProps<T> = {
  card: T;
  className?: string;
};

// Union type for all cards
export type Card = FrontBackCard | FillBlanksCard | MultipleChoiceCard | OrderedCard;
