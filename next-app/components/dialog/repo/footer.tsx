'use client'

import { Button } from '@/components/ui/button'
import { DialogFooter } from '@/components/ui/dialog'
import { useRepoDetail } from '@/hooks/util'
import { useState } from 'react'
import { useStars } from '@/hooks/stars'

export default function Footer() {
  const { setOpen, repoID, setRepoID } = useRepoDetail()
  const { search, page } = useStars()
  const [isLoaded, setIsLoaded] = useState(false)

  const deleteRepo = async () => {
    try {
      setIsLoaded(true)
      const response = await fetch(`/api/repos/detail/${repoID}`, {
        method: 'DELETE',
      })
      const data = await response.json()
      if (data) {
        setOpen(false)
        setRepoID(0)
        search(page)
      }
    } catch (error) {
      console.error(error)
    } finally {
      setIsLoaded(false)
    }
  }

  return (
    <DialogFooter className='flex gap-x-4'>
      <Button
        variant='secondary'
        onClick={() => setOpen(false)}
        disabled={isLoaded}
        className='border-slate-300 border-[1px]'
      >
        Close
      </Button>
      <Button variant='destructive' onClick={deleteRepo} disabled={isLoaded} loading={isLoaded}>
        Delete
      </Button>
      <Button onClick={() => setOpen(false)} disabled={isLoaded}>
        Save
      </Button>
    </DialogFooter>
  )
}
