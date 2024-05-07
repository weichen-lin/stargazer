'use server'

import { LanguagePie } from '@/components/charts'
import { TotalRepositories, Crontab, Recently } from './components'
import clsx from 'clsx'

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
        <TotalRepositories />
        <Crontab />
      </div>
      <LanguagePie />
      <Recently />
      <Recently />
      <Recently />
      <Recently />
    </div>
  )
}
