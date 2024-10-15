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
import { FolderPen } from 'lucide-react'
import { cn } from '@/lib/utils'
import { useCollectionContext } from '@/app/collections/hooks/useCollectionContext'
import { useFetch } from '@/hooks/util'
import { ICollection } from '@/client/collection'

export default function Rename() {
  const { collection } = useCollectionContext()

  return (
    <FloatingPanelRoot defaultText={collection.name}>
      <FloatingPanelTrigger
        title='Rename collection'
        className={cn(
          'flex items-center space-x-2 hover:bg-slate-400/20 transition-colors',
          'bg-white rounded-md border-slate-300 shadow-md p-2',
        )}
      >
        <div className='flex gap-x-2 items-center'>
          <FolderPen className='w-4 h-4' />
          <span>Rename</span>
        </div>
      </FloatingPanelTrigger>
      <FloatingPanelContent className='w-80'>
        <RenamePanel />
      </FloatingPanelContent>
    </FloatingPanelRoot>
  )
}

const RenamePanel = () => {
  const { note, error, setError, closeFloatingPanel } = useFloatingPanel()
  const { collection, isSearch, update } = useCollectionContext()
  const { run, isLoading } = useFetch<ICollection>({
    config: {
      url: `/collection/${collection.id}`,
      method: 'PATCH',
    },
    initialRun: false,
    onSuccess: data => {
      update(data)
      closeFloatingPanel()
    },
    onError: ({ error }) => {
      setError(error)
    },
  })

  const loading = isSearch || isLoading

  return (
    <FloatingPanelForm>
      <FloatingPanelBody>
        <FloatingPanelLabel htmlFor='note-input'>
          <span className='bg-slate-300 px-2 py-1'>Name</span>
        </FloatingPanelLabel>
        <FloatingPanelTextarea id='note-input' className='min-h-[80px]' maxLength={20} />
        {error && <p className='text-red-500 text-sm'>{error}</p>}
      </FloatingPanelBody>
      <FloatingPanelFooter>
        <FloatingPanelCloseButton />
        <FloatingPanelSubmitButton
          isLoading={loading}
          text='Rename'
          onClick={() => {
            run({ payload: { name: note } })
          }}
        />
      </FloatingPanelFooter>
    </FloatingPanelForm>
  )
}
