export type CardRendererProps<T> = {
  card: T;
  deckId: string;
  className?: string;
};

// User answer types for different card types
export type FillBlanksAnswer = string[];
export type MultipleChoiceAnswer = string | Record<string, boolean>;
export type OrderedAnswer = string[];
export type FrontBackAnswer = boolean;

export type UserAnswer = FillBlanksAnswer | MultipleChoiceAnswer | OrderedAnswer | FrontBackAnswer;

export type CardComponentProps<T> = Pick<CardRendererProps<T>, 'card' | 'className'> & {
  onAnswerChange?: (answer: UserAnswer) => void;
};

export type CardVerificationResult = {
  isCorrect: boolean;
  userAnswer?: UserAnswer;
};
