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

export default function Operator() {
  return (
    <FloatingPanelRoot>
      <FloatingPanelTrigger
        title='Create a repository collection'
        className='flex items-center space-x-2 px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition-colors'
      >
        <span>Create Collection</span>
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
        <FloatingPanelTextarea id='note-input' className='min-h-[80px]' maxLength={20} />
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
