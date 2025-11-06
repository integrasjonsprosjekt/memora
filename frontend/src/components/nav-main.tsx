'use client';

import { FileBox, ChevronRight, Plus, Trash2, SquarePen, Share2 } from 'lucide-react';
import { Skeleton } from '@/components/ui/skeleton';
import { use, useMemo, Suspense, useState, useEffect } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible';
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuAction,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from '@/components/ui/sidebar';
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuSeparator,
  ContextMenuTrigger,
} from '@/components/ui/context-menu';
import { cn } from '@/lib/utils';
import { getApiEndpoint, USER_ID } from '@/config/api';
import { Deck } from '@/types/deck';
import { fetchApi } from '@/lib/api/config';
import { getCurrentUser } from '@/lib/firebase/auth';

// Create a cache for deck promises
const deckPromiseCache = new Map<string, Promise<Deck[]>>();

function fetchDecks(endpoint: string): Promise<Deck[]> {
  // Return cached promise if it exists
  if (deckPromiseCache.has(endpoint)) {
    return deckPromiseCache.get(endpoint)!;
  }

  // Create new promise and cache it
  const promise = fetch(getApiEndpoint(`/v1/users/${USER_ID}/decks/${endpoint}`), {
    next: { revalidate: 300 }, // Cache for 5 minutes
  })
    .then((res) => {
      if (!res.ok) throw new Error('Failed to fetch decks');
      return res.json();
    })
    .catch((error) => {
      console.error(`Error fetching ${endpoint} decks:`, error);
      // Remove failed promise from cache
      deckPromiseCache.delete(endpoint);
      throw error;
    });

  deckPromiseCache.set(endpoint, promise);
  return promise;
}

function DeckGroup({ title, endpoint, action }: { title: string; endpoint: string; action?: React.ReactNode }) {
  const pathname = usePathname();
  // const decksPromise = useMemo(() => fetchDecks(endpoint), [endpoint]);
  // const decks = use(decksPromise);

  /**
   * Helper to check if main deck item should be active
   * Active only when NOT on a direct subpage (like /dashboard or /today)
   * Still active on nested routes like /cards/{id}
   */
  const isDeckMainActive = (deckId: string) => {
    if (!pathname) return false;

    const deckBasePath = `/decks/${deckId}`;

    // Not on this deck at all
    if (!pathname.startsWith(deckBasePath)) return false;

    const pathAfterDeck = pathname.slice(deckBasePath.length);

    // If we're exactly at /decks/{id} or /decks/{id}/
    if (!pathAfterDeck || pathAfterDeck === '/') return true;

    // Check if there's a second slash (meaning it's a nested route like /cards/{id})
    // Direct subpages like /dashboard or /today don't have a second slash
    const remainingPath = pathAfterDeck.startsWith('/') ? pathAfterDeck.slice(1) : pathAfterDeck;
    return remainingPath.includes('/');
  };

  return (
    <SidebarGroup>
      <SidebarGroupLabel className="flex items-center justify-between pr-1">
        <span>{title}</span>
        {action}
      </SidebarGroupLabel>
      <SidebarMenu>
        {/* {decks.map((deck) => {
          return <DeckItem key={deck.id} deck={deck} pathname={pathname} isDeckMainActive={isDeckMainActive} />;
        })} */}
      </SidebarMenu>
    </SidebarGroup>
  );
}

function DeckItem({
  deck,
  pathname,
  isDeckMainActive,
}: {
  deck: Pick<Deck, 'id' | 'title'>;
  pathname: string | null;
  isDeckMainActive: (deckId: string) => boolean;
}) {
  const shouldBeOpen = pathname?.startsWith(`/decks/${deck.id}`) || false;
  const [isOpen, setIsOpen] = useState(shouldBeOpen);

  useEffect(() => {
    setIsOpen(shouldBeOpen);
  }, [shouldBeOpen]);

  const hoverAnimation = 'transition-all duration-200 hover:translate-x-0.5';

  return (
    <ContextMenu>
      <ContextMenuTrigger>
        <Collapsible asChild open={isOpen} onOpenChange={setIsOpen}>
          <SidebarMenuItem>
            <SidebarMenuButton
              asChild
              tooltip={deck.title}
              className={cn(hoverAnimation, 'transition-transform duration-100 active:scale-98')}
              isActive={isDeckMainActive(deck.id)}
            >
              <Link href={`/decks/${deck.id}`}>
                <FileBox />
                <span>{deck.title}</span>
              </Link>
            </SidebarMenuButton>
            <CollapsibleTrigger asChild>
              <SidebarMenuAction className="transition-transform duration-200 ease-out data-[state=open]:rotate-90">
                <ChevronRight />
                <span className="sr-only">Toggle</span>
              </SidebarMenuAction>
            </CollapsibleTrigger>
            <CollapsibleContent className="data-[state=closed]:animate-collapsible-up data-[state=open]:animate-collapsible-down overflow-hidden transition-all duration-200 ease-out">
              <SidebarMenuSub>
                <SidebarMenuSubItem>
                  <SidebarMenuSubButton
                    asChild
                    className={cn(hoverAnimation, 'transition-transform duration-100 active:scale-95')}
                    isActive={pathname === `/decks/${deck.id}/dashboard`}
                  >
                    <Link href={`/decks/${deck.id}/dashboard`}>
                      <span>Dashboard</span>
                    </Link>
                  </SidebarMenuSubButton>
                </SidebarMenuSubItem>
                <SidebarMenuSubItem>
                  <SidebarMenuSubButton
                    asChild
                    className={cn(hoverAnimation, 'transition-transform duration-100 active:scale-95')}
                    isActive={pathname === `/decks/${deck.id}/today`}
                  >
                    <Link href={`/decks/${deck.id}/today`}>
                      <span>Today</span>
                    </Link>
                  </SidebarMenuSubButton>
                </SidebarMenuSubItem>
              </SidebarMenuSub>
            </CollapsibleContent>
          </SidebarMenuItem>
        </Collapsible>
      </ContextMenuTrigger>

      <ContextMenuContent>
        <ContextMenuItem>
          <SquarePen />
          Edit
        </ContextMenuItem>
        <ContextMenuItem>
          <Share2 />
          Share
        </ContextMenuItem>
        <ContextMenuSeparator />
        <ContextMenuItem variant="destructive">
          <Trash2 />
          Delete
        </ContextMenuItem>
      </ContextMenuContent>
    </ContextMenu>
  );
}

function DeckGroupSuspense({ title, endpoint, action }: { title: string; endpoint: string; action?: React.ReactNode }) {
  return (
    <Suspense
      fallback={
        <SidebarGroup>
          <SidebarGroupLabel className="flex items-center justify-between pr-1">
            <span>{title}</span>
            {action}
          </SidebarGroupLabel>
          <SidebarMenu>
            <Skeleton className="mx-2 h-[20px] rounded-xl" />
          </SidebarMenu>
        </SidebarGroup>
      }
    >
      <DeckGroup title={title} endpoint={endpoint} action={action} />
    </Suspense>
  );
}

export function NavMain() {
  return (
    <>
      <DeckGroupSuspense
        title="Decks"
        endpoint="owned"
        action={
          <button onClick={() => alert('Adding new deck')} className="hover:bg-accent rounded p-0.5">
            <Plus className="h-4 w-4" />
          </button>
        }
      />
      <DeckGroupSuspense title="Shared" endpoint="shared" />
    </>
  );
}
