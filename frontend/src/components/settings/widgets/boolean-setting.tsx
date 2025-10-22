import { Checkbox } from '@/components/ui/checkbox';
import { Label } from '@/components/ui/label';
import { BooleanSetting, SettingComponentProps } from '../types';

export interface BooleanSettingWidgetProps extends SettingComponentProps<BooleanSetting> {
  onValueChange?: (value: boolean) => void;
}

export function BooleanSettingWidget({ setting, className, onValueChange }: BooleanSettingWidgetProps) {
  return (
    <div className={`flex items-start space-x-3 ${className || ''}`}>
      <Checkbox
        id={setting.id}
        checked={setting.value}
        onCheckedChange={(checked) => onValueChange?.(checked === true)}
        className="mt-0.5"
      />
      <div className="space-y-1 leading-none">
        <Label htmlFor={setting.id} className="cursor-pointer">
          <strong>{setting.label}</strong>
        </Label>
        {setting.description && <p className="text-muted-foreground text-sm">{setting.description}</p>}
      </div>
    </div>
  );
}
