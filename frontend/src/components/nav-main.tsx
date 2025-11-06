'use client';

import { FileBox, ChevronRight, Plus, Trash2, SquarePen, Share2 } from 'lucide-react';
import { Skeleton } from '@/components/ui/skeleton';
import { use, useMemo, Suspense, useState, useEffect, createContext, useContext } from 'react';
import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';

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
import { EditDeckMenu } from './edit-deck-menu';
import { AddDeckMenu } from './add-deck-menu';
import { deleteDeck } from '@/app/api';

// Create a cache for deck promises
let deckPromiseCache: Promise<Deck[]> | null = null;
let cacheInvalidationCounter = 0;

function invalidateDeckCache() {
  deckPromiseCache = null;
  cacheInvalidationCounter++;
}

function fetchDecks(): Promise<Deck[]> {
  // Return cached promise if it exists
  if (deckPromiseCache) {
    return deckPromiseCache;
  }

  // Create new promise and cache it
  const promise = fetch(getApiEndpoint(`/v1/users/${USER_ID}/decks`), {
    next: { revalidate: 300 }, // Cache for 5 minutes
  })
    .then((res) => {
      if (!res.ok) throw new Error('Failed to fetch decks');
      return res.json();
    })
    .catch((error) => {
      console.error(`Error fetching decks:`, error);
      // Remove failed promise from cache
      deckPromiseCache = null;
      throw error;
    });

  deckPromiseCache = promise;
  return promise;
}

const DeckCacheContext = createContext<() => void>(() => {});

function DeckGroup({
  title,
  filterType,
  action,
  cacheKey,
}: {
  title: string;
  filterType: 'owned' | 'shared';
  action?: React.ReactNode;
  cacheKey: number;
}) {
  const pathname = usePathname();
  // Use cacheKey to force re-fetching when cache is invalidated
  const decksPromise = useMemo(() => {
    void cacheKey;
    return fetchDecks();
  }, [cacheKey]);
  const allDecks = use(decksPromise);

  const decks = useMemo(() => {
    if (filterType === 'owned') {
      return allDecks.filter((deck) => deck.owner_id === USER_ID);
    } else {
      return allDecks.filter((deck) => deck.owner_id !== USER_ID);
    }
  }, [allDecks, filterType]);

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
        {decks.map((deck) => {
          return <DeckItem key={deck.id} deck={deck} pathname={pathname} isDeckMainActive={isDeckMainActive} />;
        })}
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
  const router = useRouter();
  const invalidateCache = useContext(DeckCacheContext);
  const shouldBeOpen = pathname?.startsWith(`/decks/${deck.id}`) || false;
  const [isOpen, setIsOpen] = useState(shouldBeOpen);
  const [isEditing, setIsEditing] = useState(false);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  useEffect(() => {
    setIsOpen(shouldBeOpen);
  }, [shouldBeOpen]);

  const handleEditClose = (open: boolean) => {
    setIsEditing(open);
    if (!open) {
      invalidateCache();
    }
  };

  async function handleDelete() {
    if (!confirm(`Are you sure you want to delete ${deck.title}?`)) {
      return;
    } else {
      const res = await deleteDeck(deck.id);
      if (res.success) {
        if (shouldBeOpen) {
          // Redirect to home
          router.push('/');
        }
        invalidateCache();
      } else {
        alert('Failed to delete deck');
      }
    }
  }

  const hoverAnimation = 'transition-all duration-200 hover:translate-x-0.5';

  if (!mounted) {
    return (
      <SidebarMenuItem>
        <div className="px-2 py-1.5">
          <Skeleton className="h-5 w-full" />
        </div>
      </SidebarMenuItem>
    );
  }

  return (
    <>
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
          <ContextMenuItem onClick={() => setIsEditing(true)}>
            <SquarePen />
            Edit
          </ContextMenuItem>
          <ContextMenuItem>
            <Share2 />
            Share
          </ContextMenuItem>
          <ContextMenuSeparator />
          <ContextMenuItem onClick={() => handleDelete()} variant="destructive">
            <Trash2 />
            Delete
          </ContextMenuItem>
        </ContextMenuContent>
      </ContextMenu>
      <EditDeckMenu open={isEditing} onOpenChange={handleEditClose} deckId={deck.id} />
    </>
  );
}

function DeckGroupSuspense({
  title,
  filterType,
  action,
  cacheKey,
}: {
  title: string;
  filterType: 'owned' | 'shared';
  action?: React.ReactNode;
  cacheKey: number;
}) {
  return (
    <Suspense
      key={cacheKey}
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
      <DeckGroup title={title} filterType={filterType} action={action} cacheKey={cacheKey} />
    </Suspense>
  );
}

export function NavMain() {
  const [cacheKey, setCacheKey] = useState(cacheInvalidationCounter);
  const [isAddingDeck, setIsAddingDeck] = useState(false);

  const handleInvalidateCache = () => {
    invalidateDeckCache();
    setCacheKey(cacheInvalidationCounter);
  };

  const handleAddDeckClose = (open: boolean) => {
    setIsAddingDeck(open);
    if (!open) {
      handleInvalidateCache();
    }
  };

  return (
    <>
      <DeckCacheContext.Provider value={handleInvalidateCache}>
        <DeckGroupSuspense
          title="Decks"
          filterType="owned"
          cacheKey={cacheKey}
          action={
            <button onClick={() => setIsAddingDeck(true)} className="hover:bg-accent rounded p-0.5">
              <Plus className="h-4 w-4" />
            </button>
          }
        />
        <DeckGroupSuspense title="Shared" filterType="shared" cacheKey={cacheKey} />
      </DeckCacheContext.Provider>
      <AddDeckMenu userId={USER_ID} open={isAddingDeck} onOpenChange={handleAddDeckClose} />
    </>
  );
}
