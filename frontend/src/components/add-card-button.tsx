'use client';

import { Button } from '@/components/ui/button';

export function AddCardButton() {
  return (
    <Button
      className="h-[100px] w-full rounded-2xl border border-dashed border-[var(--border)] bg-transparent"
      onClick={() => alert('Add card clicked!')}
    >
      Add card
    </Button>
  );
}
