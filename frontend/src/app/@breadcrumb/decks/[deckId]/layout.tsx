import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';
import { getApiEndpoint } from '@/config/api';
import { Deck } from '@/types/deck';
import { DeckBreadcrumb } from './deck-breadcrumb';

export default async function Layout({
  children,
  params,
}: Readonly<{
  children: React.ReactNode;
  params: Promise<{ deckId: string }>;
}>) {
  const { deckId } = await params;

  const deck: Deck = await fetch(getApiEndpoint(`/v1/decks/${deckId}`), {
    cache: 'no-store',
  })
    .then((res) => res.json())
    .catch(() => null);

  const deckTitle = deck?.title || deckId;

  return (
    <Breadcrumb>
      <BreadcrumbList>
        <BreadcrumbItem>
          <BreadcrumbPage>Decks</BreadcrumbPage>
        </BreadcrumbItem>
        <BreadcrumbSeparator />
        <DeckBreadcrumb deckId={deckId} deckTitle={deckTitle} />
        {children}
      </BreadcrumbList>
    </Breadcrumb>
  );
}
