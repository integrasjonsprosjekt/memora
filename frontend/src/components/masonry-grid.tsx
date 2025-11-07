'use client';

import dynamic from 'next/dynamic';
import { ReactNode } from 'react';
import { Skeleton } from '@/components/ui/skeleton';

const spinner = (
  <div className="flex flex-row space-x-3 space-y-3">
    <Skeleton className="h-[145px] w-[250px] rounded-xl" />
    <Skeleton className="h-[125px] w-[250px] rounded-xl" />
    <Skeleton className="h-[130px] w-[250px] rounded-xl" />
  </div>
);

// We have to dynamically import react-responsive-masonry to avoid hydration errors,
// This will emit a loading indicator while the component is being loaded.
const ResponsiveMasonry = dynamic(() => import('react-responsive-masonry').then((mod) => mod.ResponsiveMasonry), {
  loading: () => spinner,
  ssr: false,
});
const Masonry = dynamic(() => import('react-responsive-masonry').then((mod) => mod.default), {
  loading: () => spinner,
  ssr: false,
});

interface MasonryGridProps {
  children: ReactNode;
  columnsCountBreakPoints?: Record<number, number>;
  gutter?: string;
}

/**
 * MasonryGrid component - Client-side wrapper for react-responsive-masonry
 * Provides a responsive masonry grid layout for card components.
 *
 * @param {Object} props - Component props
 * @param {ReactNode} props.children - Child elements to be rendered in the masonry grid
 * @param {Record<number, number>} props.columnsCountBreakPoints - Breakpoints for responsive column counts
 * @param {string} props.gutter - Gap between items in the grid
 */
export function MasonryGrid({
  children,
  columnsCountBreakPoints = { 350: 1, 750: 2, 900: 3, 1200: 4 },
  gutter = '1rem',
}: MasonryGridProps) {
  return (
    <ResponsiveMasonry columnsCountBreakPoints={columnsCountBreakPoints}>
      <Masonry gutter={gutter}>{children}</Masonry>
    </ResponsiveMasonry>
  );
}
