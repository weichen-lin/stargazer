'use client'

import { getUserStarsRelation } from '@/actions/neo4j'
import { useSession } from 'next-auth/react'
import { useEffect, useState } from 'react'
import { AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { useChatAlert } from '@/hooks/chat'

export default function VectorizeInfo() {
  const { data: session } = useSession()
  const [isLoading, setIsLoading] = useState(true)
  const [isVectorized, setIsVectorized] = useState(false)
  const [total, setTotal] = useState(0)
  const [vectorized, setVectorized] = useState(0)
  const { startEvent } = useChatAlert()

  useEffect(() => {
    const getCounts = async () => {
      if (session?.user?.email) {
        try {
          setIsLoading(true)
          const { total, vectorized } = await getUserStarsRelation(session?.user?.email)
          setTotal(total)
          setVectorized(vectorized)
        } catch (error) {
          console.error(error)
        } finally {
          setIsLoading(false)
        }
      }
    }
    getCounts()
  }, [session?.user?.email])

  const handleVectorize = () => {
    setIsVectorized(true)
    startEvent(true)
  }

  return (
    <AlertDescription className='flex gap-x-2 items-center my-3'>
      {isLoading ? (
        <div className='flex gap-x-2 mx-10 items-center'>
          <span className='w-6 h-4 animate-pulse bg-slate-300/60 rounded-lg'></span>
          <span>/</span>
          <span className='w-6 h-4 animate-pulse bg-slate-300/60 rounded-lg'></span>
        </div>
      ) : (
        <div className='flex gap-x-2 mx-12'>
          <span>{vectorized}</span>
          <span>/</span>
          <span>{total}</span>
        </div>
      )}
      <Button disabled={isLoading} loading={isVectorized} onClick={handleVectorize}>
        Vectorize
      </Button>
    </AlertDescription>
  )
}
