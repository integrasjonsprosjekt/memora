import { Card } from './card';

export interface Deck {
  id: string;
  cards: Card[];
  owner_id: string;
  shared_emails: string[];
  title: string;
}

export const deckDefaults: Deck = {
  id: '',
  cards: [],
  owner_id: '',
  shared_emails: [],
  title: '',
};
