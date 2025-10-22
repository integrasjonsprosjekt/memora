import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger, DialogFooter } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import { RenderSetting } from './settings-renderer';
import { useSettings } from './use-settings';

interface SettingsDialogProps {
  children?: React.ReactNode;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
}

export function SettingsDialog({ children, open, onOpenChange }: SettingsDialogProps) {
  const { settings, updateSetting, saveSettings, resetSettings, hasUnsavedChanges } = useSettings();

  const handleSave = () => {
    saveSettings();
    onOpenChange?.(false);
  };

  const handleCancel = () => {
    resetSettings();
    onOpenChange?.(false);
  };

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen && hasUnsavedChanges) {
      // If user tries to close with unsaved changes, reset to saved state
      resetSettings();
    }
    onOpenChange?.(newOpen);
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      {children && <DialogTrigger asChild>{children}</DialogTrigger>}
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle>Settings</DialogTitle>
        </DialogHeader>

        <div className="space-y-6">
          {settings.map((setting, index) => (
            <div key={setting.id}>
              <RenderSetting setting={setting} onValueChange={updateSetting} />
              {index < settings.length - 1 && <Separator className="mt-6" />}
            </div>
          ))}
        </div>

        <DialogFooter className="gap-2">
          <Button variant="outline" onClick={handleCancel}>
            Cancel
          </Button>
          <Button onClick={handleSave} disabled={!hasUnsavedChanges}>
            Save Changes
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
