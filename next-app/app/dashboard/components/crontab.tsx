'use client'

import { CalendarDays, CheckCheck, XCircle } from 'lucide-react'
import clsx from 'clsx'
import { useState, useEffect } from 'react'
import { getCrontabInfo } from '@/actions/neo4j'
import { useUser } from '@/context'
import moment from 'moment'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import HourSetting from './crontab-setting'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'

export default function Crontab() {
  const user = useUser()

  const [isLoading, setIsLoading] = useState(true)
  const [status, setStatus] = useState('')
  const [lastTriggerTime, setLastTriggerTime] = useState<Date | null>(null)
  const [hour, setHour] = useState<number | null>(null)

  useEffect(() => {
    const fetchCrontabInfo = async () => {
      const crontabInfo = await getCrontabInfo(user.email)
      setIsLoading(false)
      if (crontabInfo) {
        crontabInfo?.status && setStatus(crontabInfo.status)
        crontabInfo?.lastTriggerTime && setLastTriggerTime(crontabInfo.lastTriggerTime)
        crontabInfo?.hour && setHour(crontabInfo.hour)
      }
    }

    fetchCrontabInfo()
  }, [])

  return isLoading ? (
    <CrontabLoading />
  ) : (
    <Card className='flex flex-col h-[320px] w-full'>
      <CardHeader className='items-start pb-0 gap-y-0'>
        <CardTitle className='text-xl'>Crontab Setting</CardTitle>
        <CardDescription>Sync your github stars</CardDescription>
      </CardHeader>
      <CardContent className='flex-1 pb-0 flex flex-col gap-y-2 py-4'>
        <HourSetting
          hour={hour}
          update={e => {
            setHour(e.getHours())
          }}
        />
        <div className='flex items-start justify-start gap-x-8 w-full'>
          <StatusMapper status={status} />
          <div className='flex flex-col h-12 gap-y-1'>
            <div className='text-slate-300 dark:text-slate-500'>Last Trigger Time</div>
            <span className='text-slate-500 dark:text-slate-700'>
              {lastTriggerTime ? moment(lastTriggerTime).format('YYYY-MM-DD HH:mm:ss') : '--'}
            </span>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

const StatusDict = ['invalid github token']

const StatusMapper = ({ status }: { status: string }) => {
  const isFailed = StatusDict.includes(status) ? 'failed' : 'success'

  if (status === '')
    return (
      <div className='flex flex-col h-12 gap-y-1'>
        <div className='text-slate-300 dark:text-slate-500'>Status</div>
        --
      </div>
    )

  switch (isFailed) {
    case 'success':
      return (
        <div className='flex flex-col h-12 gap-y-1'>
          <div className='text-slate-300 dark:text-slate-500'>Status</div>
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <CheckCheck className='w-6 h-6 text-green-500' />
              </TooltipTrigger>
              <TooltipContent className='bg-white text-black border-[1px] border-green-500'>
                <p>{status}</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </div>
      )
    case 'failed':
      return (
        <div className='flex flex-col h-12 gap-y-1'>
          <div className='text-slate-300 dark:text-slate-500'>Status</div>
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <XCircle className='w-6 h-6 text-red-500' />
              </TooltipTrigger>
              <TooltipContent className='bg-white text-black border-[1px] border-red-300'>
                <p>{status}</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </div>
      )
    default:
      return (
        <div className='flex flex-col h-12 gap-y-1'>
          <div className='text-slate-300 dark:text-slate-500'>Status</div>
          --
        </div>
      )
  }
}

const CrontabLoading = () => {
  return (
    <div
      className={clsx(
        'border-[2px] border-slate-300/70 rounded-lg flex flex-col items-center justify-center gap-y-4',
        'bg-white drop-shadow-lg dark:bg-slate-300 dark:border-slate-800 dark:text-white px-8',
        'w-full h-[60%]',
      )}
    >
      <div className='flex gap-x-4 items-center w-full'>
        <CalendarDays className='w-8 h-8 text-slate-500' />
        <div className='w-full text-left text-2xl text-slate-500'>Cron Job</div>
      </div>
      <div className='w-full flex justify-between items-center'>
        <div className='animate-pulse w-28 h-4 rounded-md bg-slate-300'></div>
        <div className='flex gap-x-2'>
          <div className='animate-pulse w-12 h-4 rounded-md bg-slate-300'></div>
          <div className='animate-pulse w-12 h-4 rounded-md bg-slate-300'></div>
        </div>
      </div>
      <div className='flex items-start justify-start gap-x-8 w-full'>
        <div className='flex flex-col h-12 gap-y-1'>
          <div className='text-slate-300 dark:text-slate-500'>Status</div>
          <div className='animate-pulse w-12 h-4 rounded-md bg-slate-300'></div>
        </div>
        <div className='flex flex-col h-12 gap-y-1'>
          <div className='text-slate-300 dark:text-slate-500'>Last Trigger Time</div>
          <div className='animate-pulse w-12 h-4 rounded-md bg-slate-300'></div>
        </div>
      </div>
    </div>
  )
}
