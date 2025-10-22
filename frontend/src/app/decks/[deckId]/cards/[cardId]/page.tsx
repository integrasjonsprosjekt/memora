import { RenderCard } from '@/components/card';
import { Card } from '@/components/ui/card';
import { getApiEndpoint } from '@/config/api';

export default async function CardPage({
  params,
}: {
  params: Promise<{ deckId: string; cardId: string }>;
}) {
  const { deckId, cardId } = await params;

  const card = await fetch(getApiEndpoint(`/v1/decks/${deckId}/cards/${cardId}`), {
    cache: 'no-store',
  }).then((res) => res.json());

  return (
    <div className="flex flex-1 flex-col items-center justify-center px-4 sm:px-6 lg:px-8">
      <Card className="w-full max-w-sm -translate-y-30 transform px-10 py-5 text-xl sm:max-w-md md:max-w-lg lg:max-w-xl xl:max-w-2xl">
        <RenderCard
          key={card.id}
          card={card}
          // Counteract padding for rulers
          className="[&>hr]:-mx-10 [&>hr]:w-auto"
        />
      </Card>
    </div>
  );
}
