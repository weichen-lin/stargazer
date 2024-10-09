'use client'

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
  useFloatingPanel,
} from '@/components/ui/floating-panel'
import { useFetch } from '@/hooks/util'
import { ICollection } from '@/client/collection'
import { useCollection } from '@/app/collections/hooks'
import { FolderPen } from 'lucide-react'
import { cn } from '@/lib/utils'

export default function EditDescription() {
  return (
    <FloatingPanelRoot>
      <FloatingPanelTrigger
        title='Edit Description'
        className={cn(
          'flex items-center space-x-2 hover:bg-slate-400/20 transition-colors',
          'bg-white rounded-md border-slate-300 shadow-md p-2',
        )}
      >
        <div className='flex gap-x-2 items-center'>
          <FolderPen className='w-4 h-4' />
          <span>Edit Description</span>
        </div>
      </FloatingPanelTrigger>
      <FloatingPanelContent className='w-80'>
        <EditDescriptionPanel />
      </FloatingPanelContent>
    </FloatingPanelRoot>
  )
}

const EditDescriptionPanel = () => {
  const { note } = useFloatingPanel()

  return (
    <FloatingPanelForm>
      <FloatingPanelBody>
        <FloatingPanelLabel htmlFor='note-input'>
          <span className='bg-slate-300 px-2 py-1'>Description</span>
        </FloatingPanelLabel>
        <FloatingPanelTextarea id='note-input' className='min-h-[80px]' disabled={true} maxLength={100} />
      </FloatingPanelBody>
      <FloatingPanelFooter>
        <FloatingPanelCloseButton />
        <FloatingPanelSubmitButton isLoading={true} text='Rename' onClick={() => {}} />
      </FloatingPanelFooter>
    </FloatingPanelForm>
  )
}
