'use client';

import { Button } from '@/components/ui/button';
import { deleteDeck } from '@/app/api';
import { useRouter } from 'next/navigation';

export function DeleteDeckButton({ deckId }: { deckId: string }) {
  const router = useRouter();
  async function handler({ deckId }: { deckId: string }) {
    const res = await deleteDeck(deckId);
    if (res.success) {
      alert('Deck deleted successfully');
      router.push(`/`);
    } else {
      alert('Failed to delete deck');
    }
  }
  return (
    <>
      <Button className="bg-destructive mx-20" onClick={() => handler({ deckId })}>
        Delete card
      </Button>
    </>
  );
}
