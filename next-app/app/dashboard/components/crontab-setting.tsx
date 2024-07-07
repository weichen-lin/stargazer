'use client'

import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Clock } from 'lucide-react'
import { TimePickerInput } from '@/components/ui/timepicker'
import { useRef, useState } from 'react'
import { Button } from '@/components/ui/button'
import { syncUserStars, updateCrontabHour } from '@/actions/kafka-service'
import { useToast } from '@/components/ui/use-toast'
import { useUser } from '@/context'

interface ICrontabSetting {
  hour: number | null
  update: (date: Date) => void
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
  const { hour, update } = props
  const { toast } = useToast()
  const { email } = useUser()

  const currentDate = new Date()
  if (hour) {
    currentDate.setHours(hour)
  }

  const ref = useRef<HTMLInputElement>(null)
  const [date, setDate] = useState(currentDate)
  const [syncing, setSyncing] = useState(false)
  const [chaning, setChanging] = useState(false)

  const getStars = async () => {
    setSyncing(true)
    const { status, title, message } = await syncUserStars(email)
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

  const updateCrontab = async () => {
    setChanging(true)
    const status = await updateCrontabHour(email, date.getHours())
    if (status) {
      update(date)
      toast({
        title: 'Success',
        description: 'Crontab hour updated successfully.',
      })
    }

    setChanging(false)
  }

  return (
    <div className='w-full flex justify-between items-center'>
      <span className='text-slate-500 dark:text-slate-700 lg:text-sm'>
        {hour !== null ? `Everyday at ${formatHour(hour)}` : '--'}
      </span>
      <div className='flex gap-x-2'>
        <Button loading={syncing} onClick={getStars}>
          Start
        </Button>
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
                <TimePickerInput picker='hours' date={date} setDate={setDate} ref={ref} />
                <Button className='' onClick={updateCrontab} loading={chaning}>
                  Update
                </Button>
              </div>
            </div>
          </PopoverContent>
        </Popover>
      </div>
    </div>
  )
}
