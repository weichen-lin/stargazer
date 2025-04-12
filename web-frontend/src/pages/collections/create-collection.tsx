import {
  FloatingPanelBody,
  FloatingPanelCloseButton,
  FloatingPanelContent,
  FloatingPanelFooter,
  FloatingPanelForm,
  FloatingPanelLabel,
  FloatingPanelRoot,
  FloatingPanelSubmitButton,
  FloatingPanelTextarea,
  FloatingPanelTrigger,
} from '@/components/ui/floating-panel';
import { Plus } from 'lucide-react';
import { cn } from '@/lib/utils';

export default function CreateCollection() {
  return (
    <FloatingPanelRoot>
      <FloatingPanelTrigger
        title='Create a repository collection'
        className={cn(
          'bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition-colors'
        )}
      >
        <div className='flex gap-x-3 items-center w-[160px] justify-center'>
          <Plus className='h-4 w-4' />
          <span>Create Collection</span>
        </div>
      </FloatingPanelTrigger>
      <FloatingPanelContent className='w-80'>
        <FloatingPanelForm>
          <FloatingPanelBody>
            <FloatingPanelLabel htmlFor='note-input'>
              <span className='bg-slate-300 px-2 py-1'>Name</span>
            </FloatingPanelLabel>
            <FloatingPanelTextarea
              id='note-input'
              className='min-h-[80px]'
              maxLength={20}
            />
            {/* {error && <p className='text-red-500 text-sm'>{error}</p>} */}
          </FloatingPanelBody>
          <FloatingPanelFooter>
            <FloatingPanelCloseButton />
            <FloatingPanelSubmitButton
              isLoading={false}
              text='Create'
              onClick={() => {}}
            />
          </FloatingPanelFooter>
        </FloatingPanelForm>{' '}
      </FloatingPanelContent>
    </FloatingPanelRoot>
  );
}
