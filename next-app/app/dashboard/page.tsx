'use client'

import clsx from 'clsx'
import { TotalRepositories, Crontab, LanguagePie, RepoSearch } from './components'

export default function Dashboard() {
  return (
    <div
      className={clsx(
        'flex flex-col items-center dark:bg-black',
        'w-full flex-1 overflow-y-auto px-[2.5%] py-8 gap-y-3',
        'lg:grid lg:grid-cols-2 lg:gap-x-3 lg:gap-y-3 xl:grid-cols-3 xl:gap-x-3 xl:gap-y-3',
        'lg:content-around',
      )}
    >
      <div className='w-full rounded-md h-[300px] flex flex-col justify-between'>
        <TotalRepositories />
        <Crontab />
      </div>
      <LanguagePie />
      <RepoSearch searchKey='created_at' title='Recently Added' />
      <RepoSearch searchKey='open_issues_count' title='Repository with the most issues open' />
      <RepoSearch searchKey='last_updated_at' title='Most recently active' />
      <RepoSearch searchKey='last_modified_at' title='Recently manually updated' />
    </div>
  )
}
