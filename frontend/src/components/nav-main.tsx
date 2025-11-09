'use client';

import { FileBox, ChevronRight, Plus, Trash2, SquarePen, Share2 } from 'lucide-react';
import { Skeleton } from '@/components/ui/skeleton';
import { use, useMemo, Suspense, useState, useEffect, createContext, useContext } from 'react';
import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { useAuth } from '@/context/auth';
import { fetchApi } from '@/lib/api/config';
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
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { cn } from '@/lib/utils';
import { Deck } from '@/types/deck';
import { EditDeckMenu } from './edit-deck-menu';
import { AddDeckMenu } from './add-deck-menu';
import { deleteDeck } from '@/app/api';
import { toast } from 'sonner';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import { User } from 'firebase/auth';

interface DecksResponse {
  owned_decks: Deck[];
  shared_decks: Deck[];
}

// Create a cache for deck promises
let deckPromiseCache: Promise<DecksResponse> | null = null;
let cacheInvalidationCounter = 0;

function invalidateDeckCache() {
  deckPromiseCache = null;
  cacheInvalidationCounter++;
}

async function fetchDecks(user: User): Promise<DecksResponse> {
  // Return cached promise if it exists
  if (deckPromiseCache) {
    return deckPromiseCache;
  }

  // Create new promise and cache it
  const promise = (async () => {
    try {
      return await fetchApi<DecksResponse>('users/decks', { user });
    } catch (error) {
      console.error(`Error fetching decks:`, error);
      // Remove failed promise from cache
      deckPromiseCache = null;
      throw error;
    }
  })();

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
  const { user } = useAuth();

  const decksPromise = useMemo(() => {
    if (!user) return Promise.resolve({ owned_decks: [], shared_decks: [] });
    return fetchDecks(user);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [user, cacheKey]);

  const decksResponse = use(decksPromise);

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

  const filteredDecks = filterType === 'owned' ? decksResponse.owned_decks || [] : decksResponse.shared_decks || [];

  // Don't render the group if there are no decks
  if (filteredDecks.length === 0) {
    return null;
  }

  return (
    <SidebarGroup>
      <SidebarGroupLabel className="flex items-center justify-between pr-1">
        <span>{title}</span>
        {action}
      </SidebarGroupLabel>
      <SidebarMenu>
        {filteredDecks.map((deck) => {
          return (
            <DeckItem
              key={deck.id}
              deck={deck}
              pathname={pathname}
              isDeckMainActive={isDeckMainActive}
              isShared={filterType === 'shared'}
            />
          );
        })}
      </SidebarMenu>
    </SidebarGroup>
  );
}

// TODO: Why isn't this used?
function DeckItem({
  deck,
  pathname,
  isDeckMainActive,
  isShared = false,
}: {
  deck: Pick<Deck, 'id' | 'title'>;
  pathname: string | null;
  isDeckMainActive: (deckId: string) => boolean;
  isShared?: boolean;
}) {
  const router = useRouter();
  const { user } = useAuth();
  const invalidateCache = useContext(DeckCacheContext);
  const shouldBeOpen = pathname?.startsWith(`/decks/${deck.id}`) || false;
  const [isOpen, setIsOpen] = useState(shouldBeOpen);
  const [isEditing, setIsEditing] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
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
    if (!user) return;

    const res = await deleteDeck(user, deck.id);
    if (res.success) {
      if (shouldBeOpen) {
        // Redirect to home
        router.push('/');
      }
      invalidateCache();
      toast.success('Deck deleted', { icon: <Trash2 size={16} /> });
    } else {
      toast.error('Failed to delete deck');
    }
    setDeleteDialogOpen(false);
  }

  async function handleShare() {
    const url = `${window.location.origin}/decks/${deck.id}`;
    try {
      await navigator.clipboard.writeText(url);
      toast.success('Link copied to clipboard', { icon: <Share2 size={16} /> });
    } catch (error) {
      console.error('Failed to copy link:', error);
      toast.error('Failed to copy link');
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
                  <Tooltip delayDuration={800}>
                    <TooltipTrigger asChild>
                      <span>{deck.title}</span>
                    </TooltipTrigger>
                    <TooltipContent>
                      <p>{deck.title}</p>
                    </TooltipContent>
                  </Tooltip>
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
          <ContextMenuItem onClick={handleShare}>
            <Share2 />
            Share
          </ContextMenuItem>
          {!isShared && (
            <>
              <ContextMenuSeparator />
              <ContextMenuItem onClick={() => setDeleteDialogOpen(true)} variant="destructive">
                <Trash2 />
                Delete
              </ContextMenuItem>
            </>
          )}
        </ContextMenuContent>
      </ContextMenu>
      <EditDeckMenu open={isEditing} onOpenChange={handleEditClose} deckId={deck.id} />
      <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone. This will permanently delete the deck &quot;{deck.title}&quot; and all its
              cards.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDelete}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
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
  const { user } = useAuth();

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

  if (!user) {
    return null;
  }

  return (
    <>
      <DeckCacheContext.Provider value={handleInvalidateCache}>
        <DeckGroupSuspense
          title="Decks"
          filterType="owned"
          cacheKey={cacheKey}
          action={
            <button
              onClick={() => setIsAddingDeck(true)}
              className="hover:bg-sidebar-accent cursor-pointer rounded p-0.5"
            >
              <Plus className="h-4 w-4" />
            </button>
          }
        />
        <DeckGroupSuspense title="Shared" filterType="shared" cacheKey={cacheKey} />
      </DeckCacheContext.Provider>
      <AddDeckMenu open={isAddingDeck} onOpenChange={handleAddDeckClose} />
    </>
  );
}
