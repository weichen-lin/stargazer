'use client'

import { Dialog, DialogContent } from '@/components/ui/dialog'
import { useRepoDetail } from '@/hooks/util'
import Title from './title'
import Info from './info'
import Status from './status'
import Footer from './footer'
import { useFetch } from '@/hooks/util'
import { IRepository } from '@/client/repository'
import Tags from './tags'
import clsx from 'clsx'
import { useSearch } from '@/hooks/dashboard'

export default function Detail() {
  const { open, setOpen, repoID, setRepoID } = useRepoDetail()
  const { clear } = useSearch()
  const { data, isLoading } = useFetch<IRepository>({
    initialRun: true,
    config: {
      url: `/repository/detail/${repoID}`,
      method: 'GET',
    },
    onError: () => {
      clear()
      setOpen(false)
    },
  })

  return (
    <Dialog
      open={open}
      onOpenChange={() => {
        setOpen(false)
        setRepoID(0)
      }}
    >
      <DialogContent
        className={clsx(
          'flex flex-col justify-between h-[750px]',
          'w-full lg:min-w-[768px] lg:max-w-[768px] overflow-y-auto',
        )}
      >
        {isLoading && <Loading />}
        {!isLoading && data && (
          <div className='w-full flex flex-col gap-y-4 h-[700px] overflow-y-auto pb-12 px-2'>
            <Title {...data} />
            <Tags repo_id={repoID} />
            <Status {...data} />
            <Info {...data} />
          </div>
        )}
        <Footer />
      </DialogContent>
      <Footer />
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
    </div>
  )
}
