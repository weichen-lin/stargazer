import { LanguagePie } from '@/components/charts'
import { TotalRepositories, Crontab, RepoSearch } from './components'
import clsx from 'clsx'
import { Suspense } from 'react'
import Loading from './loading'

export const dynamic = 'force-dynamic'

export default async function Dashboard() {
  return (
    <div
      className={clsx(
        'bg-blue-50/70 dark:bg-slate-900 flex flex-col gap-3 items-center',
        'w-full lg:h-full md:flex-row md:flex-wrap md:items-center md:justify-center py-4 lg:py-0 lg:p-2 lg:gap-6 xl:gap-4 2xl:gap-12',
        'from-blue-100/40 to-blue-200/60 bg-gradient-to-br dark:from-slate-500 dark:to-slate-100',
      )}
    >
      <div
        className={clsx(
          'flex flex-col justify-between items-center backdrop-blur-md backdrop-opacity-60',
          'w-[380px] h-[380px] xl:w-[330px] 2xl:w-[380px] xl:h-[330px] 2xl:h-[380px]',
        )}
      >
        <Suspense fallback={<Loading type='total' title='' />}>
          <TotalRepositories />
        </Suspense>
        <Crontab />
      </div>
      <LanguagePie />
      <Suspense fallback={<Loading type='' title='Recently Added' />}>
        <RepoSearch searchKey='created_at' title='Recently Added' />
      </Suspense>
      <Suspense fallback={<Loading type='' title='Repository with the most issues open' />}>
        <RepoSearch searchKey='open_issues_count' title='Repository with the most issues open' />
      </Suspense>
      <Suspense fallback={<Loading type='' title='Most recently active' />}>
        <RepoSearch searchKey='last_updated_at' title='Most recently active' />
      </Suspense>
      <Suspense fallback={<Loading type='' title='Recently manually updated' />}>
        <RepoSearch searchKey='last_updated_at' title='Recently manually updated' />
      </Suspense>
    </div>
  )
}
