'use client'

import { Button } from '@/components/ui/button'
import { DialogFooter } from '@/components/ui/dialog'
import { useRepoDetail } from '@/hooks/util'
import { useFetch } from '@/hooks/util'

export default function Footer() {
  const { setOpen, repoID } = useRepoDetail()
  const { isLoading, run: deleteRepo } = useFetch({
    initialRun: false,
    config: {
      url: `/repository/detail/${repoID}`,
      method: 'DELETE',
    },
    onSuccess: () => {
      setOpen(false)
    },
  })

  return (
    <DialogFooter className='flex gap-x-4 flex-row justify-end items-end w-full'>
      <Button
        variant='secondary'
        onClick={() => setOpen(false)}
        disabled={isLoading}
        className='border-slate-300 border-[1px] w-20'
      >
        Close
      </Button>
      <Button
        className='w-20'
        variant='destructive'
        onClick={() => deleteRepo({})}
        disabled={isLoading}
        loading={isLoading}
      >
        Delete
      </Button>
    </DialogFooter>
  )
}
