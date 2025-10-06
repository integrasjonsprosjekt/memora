import { FillBlanksCardType as FillBlanksCardType } from '../types';
import CardWrapper from '../card-wrapper';

export function FillBlanksCard({ card }: { card: FillBlanksCardType }) {
  return (
    <CardWrapper>
      <p>{card.question}</p>
    </CardWrapper>
  );
}
