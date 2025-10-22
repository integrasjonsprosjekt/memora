import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';
import React from 'react';

// List of route patterns that should not be linked
const NON_LINKABLE_ROUTES = [/^\/decks$/, /^\/cards$/, /^\/decks\/[^/]+\/cards$/];

function isValidRoute(path: string): boolean {
  return !NON_LINKABLE_ROUTES.some((pattern) => pattern.test(path));
}

interface BreadcrumbPathProps {
  segments: string[];
  labelOverrides?: Map<number, string>;
  capitalizeFrom?: number; // Index from which to start capitalizing (default: 0 = all)
}

/**
 * Generates a breadcrumb navigation component from path segments
 * @param segments - Array of path segments (e.g., ['decks', 'deck-id', 'cards'])
 * @param labelOverrides - Optional map to override labels at specific indices
 * @param capitalizeFrom - Index from which to start capitalizing labels (default: 0)
 */
export function BreadcrumbPath({ segments, labelOverrides = new Map(), capitalizeFrom = 0 }: BreadcrumbPathProps) {
  return (
    <Breadcrumb>
      <BreadcrumbList>
        {segments.map((segment, index) => {
          const href = `/${segments.slice(0, index + 1).join('/')}`;
          const isLast = index === segments.length - 1;
          const shouldLink = !isLast && isValidRoute(href);
          const label = labelOverrides.get(index) || segment;
          const shouldCapitalize = index >= capitalizeFrom;

          return (
            <React.Fragment key={href}>
              <BreadcrumbItem>
                {isLast ? (
                  <BreadcrumbPage className={shouldCapitalize ? 'capitalize' : ''}>{label}</BreadcrumbPage>
                ) : shouldLink ? (
                  <BreadcrumbLink href={href} className={shouldCapitalize ? 'capitalize' : ''}>
                    {label}
                  </BreadcrumbLink>
                ) : (
                  <BreadcrumbPage className={shouldCapitalize ? 'capitalize' : ''}>{label}</BreadcrumbPage>
                )}
              </BreadcrumbItem>
              {!isLast && <BreadcrumbSeparator />}
            </React.Fragment>
          );
        })}
      </BreadcrumbList>
    </Breadcrumb>
  );
}
