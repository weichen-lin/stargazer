'use client'

import { CalendarDays, CheckCheck, XCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import clsx from 'clsx'
import { useState, useEffect } from 'react'
import { getCrontabInfo } from '@/actions/neo4j'
import { useUser } from '@/context'
import moment from 'moment'
import { useToast } from '@/components/ui/use-toast'
import { syncUserStars } from '@/actions/producer'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'

export default function Crontab() {
  const user = useUser()
  const { toast } = useToast()

  const [isLoading, setIsLoading] = useState(true)
  const [status, setStatus] = useState('')
  const [lastTriggerTime, setLastTriggerTime] = useState<Date | null>(null)
  const [syncing, setSyncing] = useState(false)

  const getStars = async () => {
    setSyncing(true)
    const { status, title, message } = await syncUserStars(user.email)
    setSyncing(false)
    if (status === 200) {
      toast({
        title,
        description: message,
      })
    } else {
      toast({
        title,
        description: message,
        variant: 'destructive',
      })
    }
  }

  useEffect(() => {
    const fetchCrontabInfo = async () => {
      const crontabInfo = await getCrontabInfo(user.email)
      console.log({ crontabInfo })
      setIsLoading(false)
      if (crontabInfo) {
        setStatus(crontabInfo?.status ?? '')
        setLastTriggerTime(crontabInfo?.lastTriggerTime ?? null)
      }
    }

    fetchCrontabInfo()
  }, [])

  return (
    <div
      className={clsx(
        'border-[2px] border-slate-300/70 rounded-lg flex flex-col items-center justify-center gap-y-4',
        'bg-white drop-shadow-lg dark:bg-slate-300 dark:border-slate-800 dark:text-white px-8',
        'w-full h-[170px] md:h-[200px]',
      )}
    >
      <div className='flex gap-x-4 items-center w-full'>
        <CalendarDays className='w-8 h-8 text-slate-500' />
        <div className='w-full text-left text-2xl text-slate-500'>Cron Job</div>
      </div>
      <div className='w-full flex justify-between items-center'>
        <span className='text-slate-500 dark:text-slate-700'>Everyday at 11:00 PM</span>
        <div className='flex gap-x-2'>
          <Button className='' loading={syncing} onClick={getStars}>
            Start
          </Button>
          <Button className=''>Setting</Button>
        </div>
      </div>
      <div className='flex items-start justify-start gap-x-8 w-full'>
        <StatusMapper status={status} />
        <div className='flex flex-col h-12 gap-y-1'>
          <div className='text-slate-300 dark:text-slate-500'>Last Trigger Time</div>
          <span className='text-slate-500 dark:text-slate-700'>
            {lastTriggerTime ? moment(lastTriggerTime).format('YYYY-MM-DD HH:mm:ss') : '--'}
          </span>
        </div>
      </div>
    </div>
  )
}

const StatusDict = ['invalid github token']

const StatusMapper = ({ status }: { status: string }) => {
  const isFailed = StatusDict.includes(status) ? 'failed' : 'success'

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
              <TooltipContent>
                <p>Add to library</p>
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
