import { Setting } from './types';

export const DEFAULT_SETTINGS: Setting[] = [
  {
    id: 'use-masonry-layout',
    type: 'boolean',
    label: 'Use masonry layout',
    description: 'Display cards in a masonry layout (reload required)',
    value: true,
  },
];

/**
 * Get a specific setting configuration by ID.
 */
export function getSettingConfig(id: string): Setting | undefined {
  return DEFAULT_SETTINGS.find((setting) => setting.id === id);
}

/**
 * Get all setting IDs.
 */
export function getAllSettingIds(): string[] {
  return DEFAULT_SETTINGS.map((setting) => setting.id);
}

/**
 * Get default value for a specific setting.
 */
export function getDefaultValue(id: string): unknown {
  const setting = getSettingConfig(id);
  return setting?.value;
}
