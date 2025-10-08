'use client';

import { useState, useEffect, useCallback } from 'react';
import { Setting } from './types';
import { DEFAULT_SETTINGS } from './settings-config';

export interface UseSettingsReturn {
  settings: Setting[];
  updateSetting: (id: string, value: unknown) => void;
  saveSettings: () => void;
  resetSettings: () => void;
  hasUnsavedChanges: boolean;
  clearAllSettings: () => void;
  clearSetting: (settingId: string) => void;
}

const SETTING_KEY_PREFIX = 'memora_setting_';

function getSettingKey(settingId: string): string {
  return `${SETTING_KEY_PREFIX}${settingId}`;
}

function loadSettingFromStorage(settingId: string, defaultValue: unknown): unknown {
  if (typeof window === 'undefined') {
    return defaultValue;
  }

  try {
    const stored = localStorage.getItem(getSettingKey(settingId));
    if (stored === null) {
      return defaultValue;
    }
    return JSON.parse(stored);
  } catch (error) {
    console.error(`Failed to load setting '${settingId}' from localStorage:`, error);
    return defaultValue;
  }
}

function saveSettingToStorage(settingId: string, value: unknown): void {
  if (typeof window === 'undefined') {
    return;
  }

  try {
    localStorage.setItem(getSettingKey(settingId), JSON.stringify(value));
  } catch (error) {
    console.error(`Failed to save setting '${settingId}' to localStorage:`, error);
  }
}

function loadAllSettingsFromStorage(): Setting[] {
  return DEFAULT_SETTINGS.map(
    (defaultSetting) =>
      ({
        ...defaultSetting,
        value: loadSettingFromStorage(defaultSetting.id, defaultSetting.value),
      }) as Setting,
  );
}

function saveAllSettingsToStorage(settings: Setting[]): void {
  settings.forEach((setting) => {
    saveSettingToStorage(setting.id, setting.value);
  });
}

function clearAllSettings(): void {
  if (typeof window === 'undefined') {
    return;
  }

  DEFAULT_SETTINGS.forEach((setting) => {
    localStorage.removeItem(getSettingKey(setting.id));
  });
}

function clearSetting(settingId: string): void {
  if (typeof window === 'undefined') {
    return;
  }

  localStorage.removeItem(getSettingKey(settingId));
}

export function useSettings(): UseSettingsReturn {
  const [settings, setSettings] = useState<Setting[]>(DEFAULT_SETTINGS);
  const [savedSettings, setSavedSettings] = useState<Setting[]>(DEFAULT_SETTINGS);

  // Load settings on mount
  useEffect(() => {
    const loadedSettings = loadAllSettingsFromStorage();
    setSettings(loadedSettings);
    setSavedSettings(loadedSettings);
  }, []);

  const updateSetting = useCallback((id: string, value: unknown) => {
    setSettings((prev) =>
      prev.map((setting) => (setting.id === id ? ({ ...setting, value } as Setting) : setting)),
    );
  }, []);

  const saveSettings = useCallback(() => {
    saveAllSettingsToStorage(settings);
    setSavedSettings(settings);
  }, [settings]);

  const resetSettings = useCallback(() => {
    setSettings(savedSettings);
  }, [savedSettings]);

  const hasUnsavedChanges = JSON.stringify(settings) !== JSON.stringify(savedSettings);

  return {
    settings,
    updateSetting,
    saveSettings,
    resetSettings,
    hasUnsavedChanges,
    clearAllSettings,
    clearSetting,
  };
}
