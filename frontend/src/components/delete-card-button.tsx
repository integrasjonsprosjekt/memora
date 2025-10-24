'use client';

import { Button } from '@/components/ui/button';
import { deleteCard } from '@/app/api';
import { useRouter } from 'next/navigation';

interface DeleteCardButtonProps {
  deckId: string;
  cardId: string;
}

export function DeleteCardButton({ deckId, cardId }: DeleteCardButtonProps) {
  const router = useRouter();
  async function handler({ deckId, cardId }: DeleteCardButtonProps) {
    const res = await deleteCard(deckId, cardId);
    if (res.success) {
      alert('Card deleted successfully');
      router.push(`/decks/${deckId}`);
    } else {
      alert('Failed to delete card');
    }
  }
  return (
    <>
      <Button className="bg-destructive mx-20" onClick={() => handler({ deckId, cardId })}>
        Delete card
      </Button>
    </>
  );
}
