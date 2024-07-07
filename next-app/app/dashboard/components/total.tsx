'use client'

import clsx from 'clsx'
import { useEffect, useState } from 'react'
import { getRepositoriesCount } from '@/actions/neo4j'
import { useUser } from '@/context'

export default function TotalRepositories() {
  const [isLoading, setIsLoading] = useState(true)
  const [total, setTotal] = useState(0)
  const { email } = useUser()
  useEffect(() => {
    const fetchTotal = async () => {
      try {
        const res = await getRepositoriesCount(email)
        setTotal(res)
      } catch (error) {
        console.error(error)
      } finally {
        setIsLoading(false)
      }
    }
    fetchTotal()
  }, [])

  return (
    <div
      className={clsx(
        'border-[2px] border-slate-300/70 rounded-lg flex flex-col items-center justify-center gap-y-1',
        'bg-white drop-shadow-lg dark:bg-slate-300 dark:border-slate-800 dark:text-white',
        'w-full h-[35%]',
      )}
    >
      {isLoading ? (
        <div className='w-1/3 h-6 bg-slate-300/60 animate-pulse'></div>
      ) : (
        <div className='w-full text-center text-4xl text-slate-500'>{total}</div>
      )}
      <div className='w-full text-center text-2xl text-slate-500'>Total Repositories</div>
    </div>
  )
}
