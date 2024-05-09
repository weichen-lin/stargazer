'use client'

import { CalendarDays, CheckCheck } from 'lucide-react'
import { Button } from '@/components/ui/button'
import clsx from 'clsx'
import { useState, useEffect } from 'react'
import { getCrontabInfo } from '@/actions/neo4j'
import { useUser } from '@/context'

export default function Crontab() {
  const user = useUser()

  useEffect(() => {
    const fetchCrontabInfo = async () => {
      console.log('fetchCrontabInfo')
      const crontabInfo = await getCrontabInfo(user.name)
      console.log({ crontabInfo })
    }

    fetchCrontabInfo()
  }, [])

  return (
    <div
      className={clsx(
        'border-[2px] border-slate-300/70 rounded-lg flex flex-col items-center justify-center gap-y-4',
        'bg-white drop-shadow-lg dark:bg-slate-300 dark:border-slate-800 dark:text-white px-8',
        'w-full h-[200px] xl:h-[190px] 2xl:h-[230px]',
      )}
    >
      <div className='flex gap-x-4 items-center w-full'>
        <CalendarDays className='w-8 h-8 text-slate-500' />
        <div className='w-full text-left text-2xl text-slate-500'>Cron Job</div>
      </div>
      <div className='w-full flex justify-between items-center'>
        <span className='text-slate-500 dark:text-slate-700'>Everyday at 11:00 PM</span>
        <div className='flex gap-x-2'>
          <Button className=''>Start</Button>
          <Button className=''>Setting</Button>
        </div>
      </div>
      <div className='flex items-start justify-start gap-x-8 w-full'>
        <div className='flex flex-col h-12 gap-y-1'>
          <div className='text-slate-300 dark:text-slate-500'>Status</div>
          <CheckCheck className='w-6 h-6 text-green-500' />
        </div>
        <div className='flex flex-col h-12 gap-y-1'>
          <div className='text-slate-300 dark:text-slate-500'>Last Trigger Time</div>
          <span className='text-slate-500 dark:text-slate-700'>2021-10-10 11:00 PM</span>
        </div>
      </div>
    </div>
  )
}
