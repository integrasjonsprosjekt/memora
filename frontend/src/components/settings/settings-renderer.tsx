import { JSX } from 'react';
import { match } from 'ts-pattern';
import { Setting, BooleanSetting, SettingComponentProps } from './types';
import { BooleanSettingWidget } from './widgets/boolean-setting';

export interface RenderSettingProps extends SettingComponentProps<Setting> {
  onValueChange?: (id: string, value: unknown) => void;
}

/**
 * Renders a setting item based on its type.
 */
export function RenderSetting({
  setting,
  className,
  onValueChange,
}: RenderSettingProps): JSX.Element {
  const settingComponent = match(setting)
    .with({ type: 'boolean' }, () => (
      <BooleanSettingWidget
        setting={setting as BooleanSetting}
        className={className}
        onValueChange={(value) => onValueChange?.(setting.id, value)}
      />
    ))
    .exhaustive();

  return settingComponent;
}
