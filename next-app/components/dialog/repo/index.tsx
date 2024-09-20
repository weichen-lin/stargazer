'use client'

import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogHeader } from '@/components/ui/dialog'
import { useRepoDetail } from '@/hooks/util'
import Title from './title'
import Info from './info'
import Body from './body'
import Footer from './footer'
import { useFetch } from '@/hooks/util'
import { IRepository } from '@/client/repository'

export default function Detail() {
  const { open, setOpen, repoID, setRepoID } = useRepoDetail()
  const { data, isLoading } = useFetch<IRepository>({
    initialRun: true,
    config: {
      url: `/repository/detail/${repoID}`,
      method: 'GET',
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
      <DialogContent className='w-[380px] md:w-[90%] lg:max-w-[768px] h-[600px] 2xl:h-[800px] flex flex-col justify-between'>
        {isLoading && (
          <div className='flex items-center justify-center h-full'>
            <div className='loader w-20 h-20'></div>
          </div>
        )}
        {!isLoading && data && (
          <>
            <DialogHeader className='flex flex-col gap-y-4 px-3'>
              <Title {...data} />
              <Info {...data} />
              <DialogDescription className='w-full text-start'>{data?.description}</DialogDescription>
            </DialogHeader>
            <Body {...data} />
            <Footer />
          </>
        )}
      </DialogContent>
    </Dialog>
  )
}
