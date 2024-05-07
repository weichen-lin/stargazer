import { CalendarDays, CheckCheck } from 'lucide-react'
import { Button } from '@/components/ui/button'

export default function Crontab() {
  return (
    <div className='border-[2px] border-slate-300/70 rounded-lg bg-white drop-shadow w-full h-[200px] flex flex-col items-center justify-center gap-y-4 px-8'>
      <div className='flex gap-x-4 items-center w-full'>
        <CalendarDays className='w-8 h-8 text-slate-500' />
        <div className='w-full text-left text-2xl text-slate-500'>Cron Job</div>
      </div>
      <div className='w-full flex justify-between items-center'>
        11.00 PM
        <div className='flex gap-x-2'>
          <Button className=''>Start</Button>
          <Button className=''>Setting</Button>
        </div>
      </div>
      <div className='flex items-start justify-start gap-x-8 w-full'>
        <div className='flex flex-col h-12 gap-y-1'>
          <div className='text-slate-300'>Status</div>
          <CheckCheck className='w-6 h-6 text-green-500' />
        </div>
        <div className='flex flex-col h-12 gap-y-1'>
          <div className='text-slate-300'>Last Trigger Time</div>
          2021-09-01 11:00 PM
        </div>
      </div>
    </div>
  )
}
