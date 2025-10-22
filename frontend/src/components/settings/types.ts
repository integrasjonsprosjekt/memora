// TODO: Refactor like card types
export type SettingType = 'boolean';

interface BaseSetting {
  id: string;
  type: SettingType;
  label: string;
  description?: string;
}

export interface BooleanSetting extends BaseSetting {
  type: 'boolean';
  value: boolean;
}

export type Setting = BooleanSetting;

export interface SettingComponentProps<T extends Setting> {
  setting: T;
  className?: string;
}
