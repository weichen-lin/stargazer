'use client'

import { Dialog, DialogContent, DialogHeader } from '@/components/ui/dialog'
import clsx from 'clsx'
import { useSearchCollection } from '@/app/stars/hook'
import { useStars } from '@/hooks/stars'
import ListCollection from './collection'
import CollectionPagination from './pagination'
import { useFetch } from '@/hooks/util'
import { ICollection } from '@/client/collection'

export default function SearchCollection() {
  const { open, setOpen } = useSearchCollection()
  const { selectedRepo } = useStars()
  const { isLoading, data } = useFetch<{
    total: number
    data: ICollection[]
  }>({
    initialRun: true,
    config: {
      url: '/collection',
      method: 'GET',
      params: {
        limit: 20,
        page: 1,
      },
    },
    onSuccess: data => {
      //  setData(data.data)
      //  setTotal(data.total)
    },
  })

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
        {data?.total && <CollectionPagination total={data.total} />}
        {isLoading && <Loading />}
        {!isLoading && data && data?.data.length === 0 && <Empty />}
        {!isLoading && data && data?.data.length > 0 && data?.data.map(e => <ListCollection key={e.id} {...e} />)}
      </DialogContent>
    </Dialog>
  )
}

const Loading = () => {
  return (
    <div className='flex flex-col items-start justify-start h-full gap-y-4 px-4'>
      <div className='w-full h-12 bg-slate-200 animate-pulse'></div>
      <div className='w-full h-12 bg-slate-200 animate-pulse'></div>
      <div className='w-full h-12 bg-slate-200 animate-pulse'></div>
      <div className='w-full h-12 bg-slate-200 animate-pulse'></div>
      <div className='w-full h-12 bg-slate-200 animate-pulse'></div>
    </div>
  )
}

const Empty = () => {
  return (
    <div className='flex flex-col items-center justify-center h-full gap-y-4'>
      <div className='text-2xl font-bold'>No Collection Found</div>
      <div className='text-lg'>Create a new collection to add repositories</div>
    </div>
  )
}
