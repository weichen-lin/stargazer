'use client'

import { Dialog, DialogContent, DialogHeader } from '@/components/ui/dialog'
import clsx from 'clsx'
import { useSearchCollection } from '@/app/stars/hook'
import { useStars } from '@/hooks/stars'
import ListCollection from './collection'
import CollectionPagination from './pagination'

export default function SearchCollection() {
  const { open, setOpen } = useSearchCollection()
  const { selectedRepo } = useStars()

  return (
    <Dialog
      open={open}
      onOpenChange={() => {
        setOpen(false)
      }}
    >
      <DialogContent
        className={clsx(
          'flex flex-col justify-start h-[650px] p-0',
          'w-full lg:min-w-[580px] lg:max-w-[580px] overflow-y-auto',
        )}
      >
        <DialogHeader className='text-2xl p-4'>Moving {selectedRepo.length} Repositories to Collection</DialogHeader>
        <CollectionPagination total={100} />
        {/* <Loading /> */}
        <ListCollection />
      </DialogContent>
    </Dialog>
  )
}

const Loading = () => {
  return (
    <div className='flex flex-col items-start justify-start h-full gap-y-4'>
      <div className='w-full h-12 bg-slate-200 animate-pulse'></div>
      <div className='w-full h-12 bg-slate-200 animate-pulse'></div>
      <div className='w-full h-12 bg-slate-200 animate-pulse'></div>
      <div className='w-full h-12 bg-slate-200 animate-pulse'></div>
      <div className='w-full h-12 bg-slate-200 animate-pulse'></div>
    </div>
  )
}
