'use server'

import clsx from 'clsx'
import { GetUser } from '@/actions'
import { getReposByKey, ISearchKey } from '@/actions/neo4j'
import Repo from './repo'

export default async function RepoSearch({ searchKey, title }: { searchKey: ISearchKey; title: string }) {
  const user = await GetUser()
  const repos = await getReposByKey(user.email, searchKey)

  return (
    <div
      className={clsx(
        'border-[2px] border-slate-300/70 rounded-lg flex flex-col items-center justify-start gap-y-2',
        'bg-white drop-shadow-lg dark:bg-slate-300 dark:border-slate-800 dark:text-white',
        'w-[380px] h-[380px] xl:w-[330px] 2xl:w-[380px] xl:h-[330px] 2xl:h-[380px]',
      )}
    >
      <div className='w-full text-left px-4 py-3 text-xl xl:text-lg text-slate-500'>{title}</div>
      <div className='w-full px-3 grid grid-rows-5 flex-1'>
        {repos.map(e => (
          <Repo {...e} key={`${searchKey}-repo-${e.full_name}`} searchKey={searchKey} />
        ))}
      </div>
    </div>
  )
}
