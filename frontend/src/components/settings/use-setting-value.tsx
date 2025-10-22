'use client';

import { useSettings } from './use-settings';

/**
 * Utility hook to get the value of a specific setting.
 * This is a convenience hook for when you only need one setting value.
 *
 * @param settingId - The ID of the setting to retrieve
 * @returns The current value of the setting, or undefined if not found
 *
 * @example
 * ```tsx
 * const isDarkMode = useSettingValue('dark-mode');
 * const useMasonry = useSettingValue('use-masonry-layout');
 *
 * if (useMasonry) {
 *   return <MasonryGrid>{children}</MasonryGrid>;
 * }
 * ```
 */
export function useSettingValue<T = unknown>(settingId: string): T | undefined {
  const { settings } = useSettings();

  const setting = settings.find((s) => s.id === settingId);
  return setting?.value as T | undefined;
}

/**
 * Utility hook to get a boolean setting value with a fallback.
 *
 * @param settingId - The ID of the boolean setting to retrieve
 * @param defaultValue - Fallback value if setting is not found (default: false)
 * @returns The boolean value of the setting or the default
 *
 * @example
 * ```tsx
 * const notifications = useBooleanSetting('notifications', true);
 * const useMasonry = useBooleanSetting('use-masonry-layout');
 * ```
 */
export function useBooleanSetting(settingId: string, defaultValue: boolean = false): boolean {
  const value = useSettingValue<boolean>(settingId);
  return value ?? defaultValue;
}

/**
 * Utility hook to get multiple setting values at once.
 *
 * @param settingIds - Array of setting IDs to retrieve
 * @returns Object with setting IDs as keys and their values
 *
 * @example
 * ```tsx
 * const { 'dark-mode': isDarkMode, 'notifications': notificationsEnabled } =
 *   useSettingValues(['dark-mode', 'notifications']);
 * ```
 */
export function useSettingValues(settingIds: string[]): Record<string, unknown> {
  const { settings } = useSettings();

  return settingIds.reduce(
    (acc, id) => {
      const setting = settings.find((s) => s.id === id);
      acc[id] = setting?.value;
      return acc;
    },
    {} as Record<string, unknown>
  );
}
