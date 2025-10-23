import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/**
 * Merge an object with defaults, replacing null/undefined values.
 *
 * @param obj - Object to merge with defaults.
 * @param defaults - Default values for the object.
 * @returns Merged object with defaults.
 */
export function withDefaults<T extends object>(obj: Partial<T>, defaults: T): T {
  if (!obj) return defaults;

  const result = { ...defaults };

  for (const key in defaults) {
    if (obj[key] != null) {
      result[key] = obj[key] as T[typeof key];
    }
  }

  return result;
}
