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
        'w-full md:flex-row md:flex-wrap md:items-center md:justify-center lg:p-12 lg:gap-6',
        'from-blue-100/40 to-blue-200/60 bg-gradient-to-br dark:from-slate-500 dark:to-slate-100',
      )}
    >
      <div className='flex flex-col w-[380px] h-[380px] justify-between items-center backdrop-blur-md backdrop-opacity-60 opacity-75'>
        <Suspense key='total' fallback={<Loading type='total' title='' />}>
          <TotalRepositories />
        </Suspense>
        <Crontab />
      </div>
      <LanguagePie />
      <Suspense key='total' fallback={<Loading type='' title='Recently Added' />}>
        <RepoSearch searchKey='created_at' title='Recently Added' />
      </Suspense>
      <Suspense key='total' fallback={<Loading type='' title='Repository with the most issues open' />}>
        <RepoSearch searchKey='open_issues_count' title='Repository with the most issues open' />
      </Suspense>
      <Suspense key='total' fallback={<Loading type='' title='Most recently active' />}>
        <RepoSearch searchKey='last_updated_at' title='Most recently active' />
      </Suspense>
      <Suspense key='total' fallback={<Loading type='' title='Recently manually updated' />}>
        <RepoSearch searchKey='last_updated_at' title='Recently manually updated' />
      </Suspense>
    </div>
  )
}
