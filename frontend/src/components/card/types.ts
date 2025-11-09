export type CardRendererProps<T> = {
  card: T;
  deckId: string;
  className?: string;
};

export type CardComponentProps<T> = Pick<CardRendererProps<T>, 'card' | 'className'>;
