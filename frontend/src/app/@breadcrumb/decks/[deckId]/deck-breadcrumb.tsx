'use client';

import { BreadcrumbItem, BreadcrumbLink, BreadcrumbPage } from '@/components/ui/breadcrumb';
import { usePathname } from 'next/navigation';

export function DeckBreadcrumb({ deckId, deckTitle }: { deckId: string; deckTitle: string }) {
  const pathname = usePathname();
  const isOnDeckPage = pathname === `/decks/${deckId}`;

  return (
    <BreadcrumbItem>
      {isOnDeckPage ? (
        <BreadcrumbPage>{deckTitle}</BreadcrumbPage>
      ) : (
        <BreadcrumbLink href={`/decks/${deckId}`}>{deckTitle}</BreadcrumbLink>
      )}
    </BreadcrumbItem>
  );
}
