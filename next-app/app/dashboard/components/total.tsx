'use server'

import clsx from 'clsx'
import { GetUser } from '@/actions'
import { getRepositoriesCount } from '@/actions/neo4j'

export default async function TotalRepositories() {
  const user = await GetUser()
  const count = await getRepositoriesCount(user.name)

  return (
    <div
      className={clsx(
        'border-[2px] border-slate-300/70 rounded-lg flex flex-col items-center justify-center gap-y-4',
        'bg-white drop-shadow-lg dark:bg-slate-300 dark:border-slate-800 dark:text-white',
        'w-full h-[148px] xl:h-[128px] 2xl:h-[168px]',
      )}
    >
      <div className='w-full text-center text-6xl text-slate-500'>{count}</div>
      <div className='w-full text-center text-2xl text-slate-500'>Total Repositories</div>
    </div>
  )
}
