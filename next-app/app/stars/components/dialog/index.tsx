'use client'

import { Dialog, DialogContent, DialogHeader } from '@/components/ui/dialog'
import clsx from 'clsx'
import ListCollection from './collection'
import CollectionPagination from './pagination'
import { useFetch } from '@/hooks/util'
import { ICollection } from '@/client/collection'
import { useStarsContext } from '@/app/stars/hook'
import { useState } from 'react'

export default function SearchCollection() {
  const { selectRepos, open, setOpen } = useStarsContext()
  const [chosen, setChosen] = useState<string | null>(null)
  const [collections, setCollections] = useState<ICollection[]>([])
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)

  const { isLoading } = useFetch<{
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
      setCollections(data.data)
      setTotal(data.total)
    },
  })

  const select = (id: string) => {
    setChosen(chosen === id ? null : id)
  }

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
          'w-full lg:min-w-[580px] lg:max-w-[580px] overflow-y-auto gap-y-0',
        )}
      >
        <DialogHeader className='text-2xl p-4'>Moving {selectRepos.length} Repositories to Collection</DialogHeader>
        {total > 0 && <CollectionPagination total={total} />}
        {isLoading && <Loading />}
        {!isLoading && collections.length === 0 && <Empty />}
        <div className='py-4'>
          {!isLoading &&
            collections.length > 0 &&
            collections.map(e => <ListCollection key={e.id} {...e} chosen={chosen} select={select} />)}
        </div>
      </DialogContent>
    </Dialog>
  )
}

const Loading = () => {
  return (
    <div className='flex flex-col items-start justify-start h-full gap-y-4 p-4'>
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
