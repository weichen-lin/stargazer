import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { useFetch } from '@/hooks/util'
import { useState } from 'react'
import { useCollection } from '@/app/collections/hooks'

interface AlertDialogProps {
  id: string
  open: boolean
  close: () => void
}

export default function AlertDialog(props: AlertDialogProps) {
  const { id, open, close } = props
  const [value, setValue] = useState('')
  const { data, setData } = useCollection()
  const { run, isLoading } = useFetch<string>({
    initialRun: false,
    config: {
      url: '/collection',
      method: 'DELETE',
    },
    onSuccess: () => {
      const newData = data.filter(d => d.id !== id)
      setData(newData)
    },
  })

  return (
    <Dialog
      open={open}
      onOpenChange={() => {
        if (!isLoading) {
          close()
        }
      }}
    >
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Alert</DialogTitle>
        </DialogHeader>
        <DialogDescription>Are you sure you want to delete this collection?</DialogDescription>
        <Input placeholder='Type "delete" to confirm' value={value} onChange={e => setValue(e.target.value)} />
        <DialogFooter>
          <Button variant='ghost' onClick={close} disabled={isLoading}>
            Cancel
          </Button>
          <Button
            variant='destructive'
            disabled={value !== 'delete'}
            onClick={() => {
              run({ payload: { id } })
            }}
            loading={isLoading}
          >
            Delete
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
