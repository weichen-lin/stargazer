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

export default function Rename() {
  return (
    <FloatingPanelRoot>
      <FloatingPanelTrigger
        title='Rename collection'
        className='flex items-center space-x-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition-colors'
      >
        <div className='flex gap-x-2 items-center'>
          <FolderPen className='w-4 h-4' />
          <span>Rename</span>
        </div>
      </FloatingPanelTrigger>
      <FloatingPanelContent className='w-80'>
        <FloatingPanel />
      </FloatingPanelContent>
    </FloatingPanelRoot>
  )
}

const FloatingPanel = () => {
  const { note } = useFloatingPanel()
  const { data, setData, setIsFetching, loading } = useCollection()
  const { run } = useFetch<ICollection>({
    initialRun: false,
    config: {
      url: '/collection',
      method: 'POST',
    },
    onSuccess: e => {
      setData([e, ...data])
    },
  })

  return (
    <FloatingPanelForm>
      <FloatingPanelBody>
        <FloatingPanelLabel htmlFor='note-input'>
          <span className='bg-slate-300 px-2 py-1'>Name</span>
        </FloatingPanelLabel>
        <FloatingPanelTextarea id='note-input' className='min-h-[80px]' />
      </FloatingPanelBody>
      <FloatingPanelFooter>
        <FloatingPanelCloseButton />
        <FloatingPanelSubmitButton
          isLoading={loading}
          text='Create'
          onClick={() => {
            try {
              setIsFetching(true)
              run({ payload: { name: note } })
            } finally {
              setIsFetching(false)
            }
          }}
        />
      </FloatingPanelFooter>
    </FloatingPanelForm>
  )
}
