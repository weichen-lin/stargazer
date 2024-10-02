'use client'

import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Clock } from 'lucide-react'
import { TimePickerInput } from '@/components/ui/timepicker'
import { useRef, useState } from 'react'
import { Button } from '@/components/ui/button'
import { useFetch } from '@/hooks/util'
import { useToast } from '@/components/ui/use-toast'
import moment from 'moment'

interface ICrontabSetting {
  triggered_at: Date | null
  update: () => void
}

function formatHour(hour: number) {
  if (hour === 0) {
    return '12:00 AM'
  } else if (hour < 12) {
    return hour + ':00 AM'
  } else if (hour === 12) {
    return '12:00 PM'
  } else {
    return hour - 12 + ':00 PM'
  }
}

export default function HourSetting(props: ICrontabSetting) {
  const { triggered_at, update } = props
  const { isLoading, run: syncRepository } = useFetch<string>({
    initialRun: false,
    config: {
      url: '/repository/sync-repository',
      method: 'GET',
    },
    onSuccess: () => {
      toast({
        title: 'Scheduled: Catch up',
        description: 'Start Sync Repository, You will receive a notification when it is done',
      })
    },
    onError: error => {
      console.log(error)
      toast({
        title: 'Action Restricted',
        description: `You can only use this service after ${error?.expires ?? ''}. Please try again later.`,
        variant: 'destructive',
      })
    },
  })
  const { isLoading: updateLoading, run: updateCrontab } = useFetch<any>({
    initialRun: false,
    config: {
      url: '/crontab',
      method: 'PATCH',
    },
    onSuccess: data => {
      const newDate = moment(data.triggered_at).toDate()
      setDate(newDate)
      toast({
        title: 'Scheduled: Catch up',
        description: 'Update Crontab Hour Successfully',
      })
    },
  })
  const { toast } = useToast()

  const ref = useRef<HTMLInputElement>(null)
  const [date, setDate] = useState(triggered_at ?? new Date())
  const [settingDate, setSettingDate] = useState(triggered_at ?? new Date())

  return (
    <div className='w-full flex justify-between items-center'>
      <span className='text-slate-500 dark:text-slate-700'>
        {triggered_at !== null ? `Everyday at ${formatHour(date.getHours())}` : '--'}
      </span>
      <div className='flex gap-x-2'>
        <Popover>
          <PopoverTrigger>
            <div className='px-2 py-[6px] border-slate-700 border-[1px] rounded-md'>Setting</div>
          </PopoverTrigger>
          <PopoverContent className='w-[250px]'>
            <div className='flex flex-col gap-y-4'>
              <div className='flex items-center gap-x-2'>
                <Clock size={18} />
                <span>Set Crontab Hour</span>
              </div>
              <div className='flex gap-x-4'>
                <TimePickerInput picker='hours' date={settingDate} setDate={setSettingDate} ref={ref} />
                <Button
                  className=''
                  onClick={() => {
                    const dateString = moment(settingDate).utc().format('YYYY-MM-DDTHH:mm:ss[Z]')
                    updateCrontab({ params: { triggered_at: dateString } })
                  }}
                  loading={updateLoading}
                >
                  Update
                </Button>
              </div>
            </div>
          </PopoverContent>
        </Popover>
        <Button loading={isLoading} onClick={() => syncRepository({})}>
          Start
        </Button>
      </div>
    </div>
  )
}
