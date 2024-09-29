'use client'

import moment from 'moment'
import HourSetting from './crontab-setting'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { useFetch } from '@/hooks/util'
import { ICrontab } from '@/client/crontab-client'
import { Timer, Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'

const parseTime = (s: string | null): Date | null => {
  const m = moment(s)

  if (!m.isValid()) {
    return null
  }

  return m.toDate()
}

export default function Crontab() {
  const { data, isLoading, run } = useFetch<ICrontab>({
    initialRun: true,
    config: {
      url: '/crontab',
      method: 'GET',
    },
  })

  const triggerAt = parseTime(data?.triggered_at ?? null)

  return (
    <Card className='flex flex-col min-h-[320px] max-h-[320px] w-full max-w-[380px] md:max-w-none'>
      <CardHeader className='items-start pb-0 gap-y-0'>
        <CardTitle className='text-xl'>Crontab Setting</CardTitle>
        <CardDescription>Sync your github stars</CardDescription>
      </CardHeader>
      <CardContent className='flex-1 pb-0 flex flex-col gap-y-2 py-1'>
        {!isLoading && data && (
          <div className='items-center justify-center mt-8'>
            <HourSetting hour={triggerAt?.getHours() ?? null} update={e => {}} />
            <div className='flex flex-col items-start justify-start gap-y-3 w-full mt-2'>
              <div className='flex flex-col h-12 gap-y-1'>
                <div className='text-slate-300 dark:text-slate-500'>Status</div>
                <span className='text-slate-500 dark:text-slate-700'>{data?.status ?? '--'}</span>
              </div>
              <div className='flex flex-col h-12 gap-y-1'>
                <div className='text-slate-300 dark:text-slate-500'>Last Trigger Time</div>
                <span className='text-slate-500 dark:text-slate-700'>
                  {data?.last_triggered_at ? moment(data?.last_triggered_at).format('YYYY-MM-DD HH:mm:ss') : '--'}
                </span>
              </div>
            </div>
          </div>
        )}
        {isLoading && <CrontabLoading />}
        {!isLoading && !data && <EmptyContent onCreate={() => run({})} />}
      </CardContent>
    </Card>
  )
}

const EmptyContent = (props: { onCreate: () => void }) => {
  const { onCreate } = props
  const { isLoading, run } = useFetch<ICrontab>({
    initialRun: false,
    config: {
      url: '/crontab',
      method: 'POST',
    },
  })

  return (
    <div className='w-full flex flex-col items-center justify-center mt-4'>
      <div className='w-32 h-32 relative'>
        <Timer className='w-full h-full text-gray-200' />
        <div className='absolute inset-0 flex items-center justify-center'>
          <Plus className='w-8 h-8 text-gray-400' />
        </div>
      </div>
      <p className='text-center text-gray-500 mb-4'>No avaliable crontab yet.</p>
      <Button
        loading={isLoading}
        onClick={async () => {
          await run({})
          onCreate()
        }}
      >
        Create
        <Plus className='w-4 h-4 ml-2' />
      </Button>
    </div>
  )
}

const CrontabLoading = () => {
  return (
    <div className='flex flex-col gap-y-4 w-full mt-8'>
      <div className='w-full bg-slate-200 h-8 rounded-lg'></div>
      <div className='w-1/3 bg-slate-200 h-8 rounded-lg'></div>
      <div className='w-full bg-slate-200 h-8 rounded-lg'></div>
    </div>
  )
}
