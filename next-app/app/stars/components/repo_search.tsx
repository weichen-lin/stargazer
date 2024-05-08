'use server'

import clsx from 'clsx'
import { GetUser } from '@/actions'
import { getReposByKey } from '@/actions/neo4j'
import Repo from './repo'

export default async function RepoSearch({ searchKey, title }: { searchKey: string; title: string }) {
  const user = await GetUser()
  const repos = await getReposByKey(user.name, searchKey)

  return (
    <div
      className={clsx(
        'border-[2px] border-slate-300/70 rounded-lg flex flex-col items-center justify-start gap-y-2',
        'bg-white drop-shadow-lg dark:bg-slate-300 dark:border-slate-800 dark:text-white',
        'w-[380px] h-[380px]',
      )}
    >
      <div className='w-full text-left px-4 py-3 text-xl text-slate-500'>{title}</div>
      {repos.map(e => (
        <Repo {...e} key={`${searchKey}-repo-${e.full_name}`} searchKey={searchKey} />
      ))}
    </div>
  )
}
