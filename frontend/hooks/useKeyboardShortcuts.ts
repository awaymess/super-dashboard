'use client';

import { useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '@/store';
import { toggleCommandPalette } from '@/store/slices/uiSlice';
import { ROUTES, KEYBOARD_SHORTCUTS } from '@/lib/constants';

export function useKeyboardShortcuts() {
  const router = useRouter();
  const dispatch = useDispatch();
  const commandPaletteOpen = useSelector((state: RootState) => state.ui.commandPaletteOpen);

  const handleKeyDown = useCallback(
    (event: KeyboardEvent) => {
      if (event.target instanceof HTMLInputElement || event.target instanceof HTMLTextAreaElement) {
        return;
      }

      if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === 'k') {
        event.preventDefault();
        dispatch(toggleCommandPalette());
        return;
      }

      if (commandPaletteOpen) {
        if (event.key === 'Escape') {
          dispatch(toggleCommandPalette());
        }
        return;
      }

      switch (event.key.toLowerCase()) {
        case 'b':
          router.push(ROUTES.BETTING);
          break;
        case 's':
          router.push(ROUTES.STOCKS);
          break;
        case 'd':
          router.push(ROUTES.DASHBOARD);
          break;
        case 'p':
          router.push(ROUTES.PAPER_TRADING);
          break;
        case 'a':
          router.push(ROUTES.ANALYTICS);
          break;
        case '/':
          event.preventDefault();
          dispatch(toggleCommandPalette());
          break;
        case 'escape':
          break;
      }
    },
    [dispatch, router, commandPaletteOpen]
  );

  useEffect(() => {
    window.addEventListener('keydown', handleKeyDown);
    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [handleKeyDown]);

  return {
    shortcuts: KEYBOARD_SHORTCUTS,
  };
}
