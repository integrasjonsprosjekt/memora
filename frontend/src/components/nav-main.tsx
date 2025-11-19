'use client';

import { FileBox, ChevronRight, Plus, Trash2, SquarePen, Copy } from 'lucide-react';
import { Skeleton } from '@/components/ui/skeleton';
import { useState, useEffect } from 'react';
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

interface DecksResponse {
  owned_decks: Deck[];
  shared_decks: Deck[];
}

function DeckItem({
  deck,
  pathname,
  isDeckMainActive,
  isShared = false,
  onDeckUpdated,
}: {
  deck: Pick<Deck, 'id' | 'title'>;
  pathname: string | null;
  isDeckMainActive: (deckId: string) => boolean;
  isShared?: boolean;
  onDeckUpdated: () => void;
}) {
  const router = useRouter();
  const { user } = useAuth();
  const shouldBeOpen = pathname?.startsWith(`/decks/${deck.id}`) || false;
  const [isOpen, setIsOpen] = useState(shouldBeOpen);
  const [isEditing, setIsEditing] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  useEffect(() => {
    setIsOpen(shouldBeOpen);
  }, [shouldBeOpen]);

  const handleEditClose = (open: boolean) => {
    setIsEditing(open);
    if (!open) {
      onDeckUpdated();
    }
  };

  async function handleDelete() {
    if (!user) return;

    const res = await deleteDeck(user, deck.id);
    if (res.success) {
      if (shouldBeOpen) {
        router.push('/');
      }
      onDeckUpdated();
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
      toast.success('Link copied to clipboard', { icon: <Copy size={16} /> });
    } catch (error) {
      console.error('Failed to copy link:', error);
      toast.error('Failed to copy link');
    }
  }

  const hoverAnimation = 'transition-all duration-200 hover:translate-x-0.5';

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
          {!isShared && (
            <ContextMenuItem onClick={() => setIsEditing(true)}>
              <SquarePen />
              Edit
            </ContextMenuItem>
          )}
          <ContextMenuItem onClick={handleShare}>
            <Copy />
            Copy link
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

function DeckGroup({
  title,
  decks,
  filterType,
  action,
  pathname,
  onDeckUpdated,
}: {
  title: string;
  decks: Deck[];
  filterType: 'owned' | 'shared';
  action?: React.ReactNode;
  pathname: string | null;
  onDeckUpdated: () => void;
}) {
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

  // Don't render the group if there are no decks
  if (decks.length === 0 && !action) {
    return null;
  }

  return (
    <SidebarGroup>
      <SidebarGroupLabel className="flex items-center justify-between pr-1">
        <span>{title}</span>
        {action}
      </SidebarGroupLabel>
      <SidebarMenu>
        {decks.map((deck) => {
          return (
            <DeckItem
              key={deck.id}
              deck={deck}
              pathname={pathname}
              isDeckMainActive={isDeckMainActive}
              isShared={filterType === 'shared'}
              onDeckUpdated={onDeckUpdated}
            />
          );
        })}
      </SidebarMenu>
    </SidebarGroup>
  );
}

export function NavMain() {
  const pathname = usePathname();
  const { user } = useAuth();
  const [decks, setDecks] = useState<DecksResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isAddingDeck, setIsAddingDeck] = useState(false);
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  // Fetch decks whenever user changes or refreshTrigger changes
  useEffect(() => {
    if (!user) {
      setDecks(null);
      setIsLoading(false);
      return;
    }

    let isCancelled = false;

    const fetchDecksData = async () => {
      try {
        setIsLoading(true);
        const response = await fetchApi<DecksResponse>('users/decks', { user });
        if (!isCancelled) {
          setDecks(response);
        }
      } catch (error) {
        console.error('Error fetching decks:', error);
        if (!isCancelled) {
          setDecks({ owned_decks: [], shared_decks: [] });
        }
      } finally {
        if (!isCancelled) {
          setIsLoading(false);
        }
      }
    };

    fetchDecksData();

    return () => {
      isCancelled = true;
    };
  }, [user, refreshTrigger]);

  const handleRefresh = () => {
    setRefreshTrigger((prev) => prev + 1);
  };

  const handleAddDeckClose = (open: boolean) => {
    setIsAddingDeck(open);
    if (!open) {
      handleRefresh();
    }
  };

  if (!user) {
    return null;
  }

  if (isLoading) {
    return (
      <>
        <SidebarGroup>
          <SidebarGroupLabel className="flex items-center justify-between pr-1">
            <span>Decks</span>
            <button
              onClick={() => setIsAddingDeck(true)}
              className="hover:bg-sidebar-accent cursor-pointer rounded p-0.5"
            >
              <Plus className="h-4 w-4" />
            </button>
          </SidebarGroupLabel>
          <SidebarMenu>
            <Skeleton className="mx-2 h-[20px] rounded-xl" />
          </SidebarMenu>
        </SidebarGroup>
        <AddDeckMenu open={isAddingDeck} onOpenChange={handleAddDeckClose} />
      </>
    );
  }

  return (
    <>
      <DeckGroup
        title="Decks"
        decks={decks?.owned_decks || []}
        filterType="owned"
        pathname={pathname}
        onDeckUpdated={handleRefresh}
        action={
          <button
            onClick={() => setIsAddingDeck(true)}
            className="hover:bg-sidebar-accent cursor-pointer rounded p-0.5"
          >
            <Plus className="h-4 w-4" />
          </button>
        }
      />
      <DeckGroup
        title="Shared"
        decks={decks?.shared_decks || []}
        filterType="shared"
        pathname={pathname}
        onDeckUpdated={handleRefresh}
        action={null}
      />
      <AddDeckMenu open={isAddingDeck} onOpenChange={handleAddDeckClose} />
    </>
  );
}
