'use server'

import { LanguagePie } from '@/components/charts'
import { TotalRepositories, Crontab } from './components'

export default async function Dashboard() {
  return (
    <div className='w-full bg-blue-50/70 flex flex-col lg:flex-row lg:flex-wrap gap-3 items-center'>
      <div className='flex flex-col w-[380px] h-[380px] justify-between items-center'>
        <TotalRepositories />
        <Crontab />
      </div>
      <LanguagePie />
      <LanguagePie />
      <LanguagePie />
      <LanguagePie />
      <LanguagePie />
    </div>
  )
}
