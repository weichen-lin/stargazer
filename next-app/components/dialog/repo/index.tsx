'use client'

import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader } from '@/components/ui/dialog'
import { useRepoDetail } from '@/hooks/util'
import { useState, useEffect } from 'react'
import { IRepoDetail } from '@/actions/neo4j'
import Title from './title'
import Info from './info'
import Body from './body'
import Footer from './footer'

export default function Detail() {
  const { open, setOpen, repoID, setRepoID } = useRepoDetail()
  const [isLoaded, setIsLoaded] = useState(false)
  const [detail, setDetail] = useState<IRepoDetail | null>(null)

  useEffect(() => {
    const getRepoDetail = async () => {
      try {
        setIsLoaded(true)
        const response = await fetch(`/api/repos/detail/${repoID}`)
        const data = await response.json()
        setDetail(data)
      } catch (error) {
        console.error(error)
      } finally {
        setIsLoaded(false)
      }
    }

    getRepoDetail()
  }, [])

  return (
    <Dialog
      open={open}
      onOpenChange={() => {
        setOpen(false)
        setRepoID(0)
      }}
    >
      <DialogContent className='w-[400px] md:w-[90%] lg:max-w-[768px] h-[800px] flex flex-col justify-between'>
        {isLoaded && (
          <div className='flex items-center justify-center h-full'>
            <div className='loader w-20 h-20'></div>
          </div>
        )}
        {!isLoaded && detail && (
          <>
            <DialogHeader className='flex flex-col gap-y-4'>
              <Title {...detail} />
              <Info {...detail} />
              <DialogDescription>{detail?.description}</DialogDescription>
            </DialogHeader>
            <Body {...detail} />
            <Footer />
          </>
        )}
      </DialogContent>
    </Dialog>
  )
}
