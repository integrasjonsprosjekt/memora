'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { AddCardMenu } from '@/components/add-card-menu';

interface AddCardButtonProps {
  deckId: string;
}

export function AddCardButton({ deckId }: AddCardButtonProps) {
  const [open, setOpen] = useState(false);

  return (
    <>
      <Button
        className="text-foreground h-[125px] w-full rounded-2xl border border-dashed border-[var(--border)] bg-transparent"
        onClick={() => setOpen(true)}
      >
        Add card
      </Button>
      <AddCardMenu open={open} onOpenChange={setOpen} deckId={deckId} />
    </>
  );
}
