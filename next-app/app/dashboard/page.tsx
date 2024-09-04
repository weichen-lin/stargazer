'use client'

import clsx from 'clsx'
import { TotalRepositories, Crontab, RepoSearch, LanguageDistribution } from './components'

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
      <LanguageDistribution />
      <RepoSearch />
      <Crontab />
    </div>
  )
}
