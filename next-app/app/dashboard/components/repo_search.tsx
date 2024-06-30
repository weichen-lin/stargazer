'use client'

import clsx from 'clsx'
import { GetUser } from '@/actions'
import { getReposByKey, ISearchKey } from '@/actions/neo4j'
import Repo from './repo'
import { useEffect, useState } from 'react'
import { IRepoAtDashboard } from '@/actions/neo4j'
import { useUser } from '@/context'

export default function RepoSearch({ searchKey, title }: { searchKey: ISearchKey; title: string }) {
  const [repos, setRepos] = useState<IRepoAtDashboard[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const { email } = useUser()

  useEffect(() => {
    const fetchData = async () => {
      try {
        const res = await getReposByKey(email, searchKey)
        setRepos(res)
      } catch (error) {
        console.error(error)
      } finally {
        setIsLoading(false)
      }
    }
    fetchData()
  }, [])

  return (
    <div
      className={clsx(
        'border-[2px] border-slate-300/70 rounded-lg flex flex-col items-center justify-start gap-y-2',
        'bg-white drop-shadow-lg dark:bg-slate-300 dark:border-slate-800 dark:text-white',
        'h-[300px] w-full',
      )}
    >
      <div className='w-full text-left px-4 py-3 text-xl xl:text-lg text-slate-500'>{title}</div>
      {!isLoading && repos.length > 0 && (
        <div className='w-full px-3 grid grid-rows-5 flex-1'>
          {repos.map(e => (
            <Repo {...e} key={`${searchKey}-repo-${e.full_name}`} searchKey={searchKey} />
          ))}
        </div>
      )}
      {!isLoading && repos.length === 0 && (
        <div className='w-full h-full flex flex-col items-center justify-center gap-y-2'>
          <div className='text-slate-300 dark:text-slate-500'>No data available now</div>
        </div>
      )}
      {isLoading && (
        <div className='w-full h-full flex flex-col gap-y-1'>
          <div className='bg-slate-300/60 h-10 w-[95%] rounded-md animate-pulse mx-auto'></div>
          <div className='bg-slate-300/60 h-10 w-[95%] rounded-md animate-pulse mx-auto'></div>
          <div className='bg-slate-300/60 h-10 w-[95%] rounded-md animate-pulse mx-auto'></div>
          <div className='bg-slate-300/60 h-10 w-[95%] rounded-md animate-pulse mx-auto'></div>
          <div className='bg-slate-300/60 h-10 w-[95%] rounded-md animate-pulse mx-auto'></div>
        </div>
      )}
    </div>
  )
}
