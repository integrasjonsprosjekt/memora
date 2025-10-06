import { CardDisplay } from '@/components/card';
import { Card } from '@/components/card';
import { notFound } from 'next/navigation';
import { AddCardButton } from '@/components/add-card-button';

export default async function DeckPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;

  const decks_cards: Record<string, string[]> = {
    '1': ['gSFxN0hqouvSKyz1jOfZ'],
    '2': [
      '9ybE9uAbhBfoiPa3dblp',
      'YzAn6ycO5vrucqOZVzkR',
      'gSFxN0hqouvSKyz1jOfZ',
      'iHsbWpt6MiH8bmVbSf69',
      'mRamakLsGWKPB8s9jTLX',
    ],
    '3': [],
  };

  if (!(slug in decks_cards)) {
    return notFound();
  }

  const cardIds = decks_cards[slug];

  // TODO: Load on scroll
  // TODO: DO NOT HARDCODE URI
  let cards: Card[] = [];
  try {
    cards = await Promise.all(
      cardIds.map(async (cardId) => {
        const response = await fetch(`http://localhost:8080/api/v1/cards/${cardId}`);
        if (!response.ok) {
          throw new Error(`Failed to fetch card ${cardId}`);
        }
        return response.json();
      }),
    );
  } catch (error) {
    console.error('Error fetching cards:', error);
    return (
      <div className="p-8">
        <h1 className="text-3xl font-bold">Error loading cards</h1>
        <p className="mt-4 text-[var(--muted-foreground)]">
          There was a problem loading the cards for this deck.
        </p>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-4 sm:px-6 lg:px-8">
      <header className="mb-8 lg:mb-12">
        <h1 className="text-2xl font-bold sm:text-3xl">Deck {slug}</h1>
        <p className="text-[var(--muted-foreground) mt-1 text-lg">
          {cards.length} card{cards.length !== 1 ? 's' : ''}
        </p>
      </header>

      <div className="grid auto-rows-auto grid-cols-[repeat(auto-fill,minmax(200px,1fr))] gap-4">
        <AddCardButton />
        {cards.map((card) => (
          <div key={card.id}>
            <CardDisplay card={card} />
          </div>
        ))}
      </div>
    </div>
  );
}
