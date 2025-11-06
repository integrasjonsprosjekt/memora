import { RenderCard } from '@/components/card';
import { getApiEndpoint } from '@/config/api';

export default async function CardPage({ params }: { params: Promise<{ deckId: string; cardId: string }> }) {
  const { deckId, cardId } = await params;

  const card = await fetch(getApiEndpoint(`/v1/decks/${deckId}/cards/${cardId}`), {
    cache: 'no-store',
  }).then((res) => res.json());

  return (
    <RenderCard
      key={card.id}
      card={card}
    />
  );
}
