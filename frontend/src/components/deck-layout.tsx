'use client';

import { ReactNode } from 'react';
import { useBooleanSetting } from '@/components/settings/use-setting-value';
import { MasonryGrid } from '@/components/masonry-grid';

interface DeckLayoutProps {
  children: ReactNode;
}

export function DeckLayout({ children }: DeckLayoutProps) {
  const useMasonryLayout = useBooleanSetting('use-masonry-layout', true);

  if (useMasonryLayout) {
    return <MasonryGrid>{children}</MasonryGrid>;
  }

  return (
    <div className="grid grid-cols-[repeat(auto-fill,minmax(200px,1fr))] gap-4">{children}</div>
  );
}
