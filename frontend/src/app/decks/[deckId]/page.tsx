import { RenderCardThumbnail } from '@/components/card/card-renderer';
import { Card } from '@/types/card';
import { Deck, deckDefaults } from '@/types/deck';
import { AddCardButton } from '@/components/add-card-button';
import { DeckLayout } from '@/components/deck-layout';
import { getApiEndpoint } from '@/config/api';
import { withDefaults } from '@/lib/utils';

export default async function DeckPage({ params }: { params: Promise<{ deckId: string }> }) {
  const { deckId } = await params;

  const deck: Deck | null = await fetch(getApiEndpoint(`/v1/decks/${deckId}/cards/`), {
    cache: 'no-store',
  })
    .then((res) => res.json())
    .then((data) => withDefaults(data, deckDefaults))
    .catch((error) => {
      console.error(`Error fetching deck:`, error);
      return null;
    });

  if (!deck) {
    return (
      <div className="p-8">
        <h1 className="text-3xl font-bold">Error loading deck</h1>
        <p className="mt-4 text-[var(--muted-foreground)]">There was a problem loading this deck.</p>
      </div>
    );
  }

  const cards: Card[] = deck.cards || [];

  return (
    <div className="container mx-auto px-4 py-4 sm:px-6 lg:px-8">
      <header className="mb-8 lg:mb-12">
        <h1 className="text-2xl font-bold sm:text-3xl">{deck.title}</h1>
        <p className="text-[var(--muted-foreground) mt-1 text-lg">
          {cards.length} card{cards.length !== 1 ? 's' : ''}
        </p>
      </header>

      <DeckLayout>
        <AddCardButton deckId={deckId} />
        {cards.map((card) => (
          <RenderCardThumbnail key={card.id} card={card} deckId={deckId} />
        ))}
      </DeckLayout>
    </div>
  );
}
