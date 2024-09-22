'use client'

import { IRepository } from '@/client/repository'
import { Star, Eye, Clock, AlertCircle } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'

export default function Status(props: IRepository) {
  const { stargazers_count, open_issues_count, updated_at, watchers_count } = props

  return (
    <Card>
      <CardHeader>
        <CardTitle>Repository Stats</CardTitle>
        <CardDescription>Key metrics for this repository</CardDescription>
      </CardHeader>
      <CardContent className='grid grid-cols-1 lg:grid-cols-4 gap-4'>
        <div className='flex flex-col items-center'>
          <Star className='h-8 w-8 text-yellow-400' />
          <span className='mt-2 text-2xl font-semibold'>{stargazers_count}</span>
          <span className='text-sm text-gray-500'>Stars</span>
        </div>
        <div className='flex flex-col items-center'>
          <Eye className='h-8 w-8 text-blue-400' />
          <span className='mt-2 text-2xl font-semibold'>{watchers_count}</span>
          <span className='text-sm text-gray-500'>Watchers</span>
        </div>
        <div className='flex flex-col items-center'>
          <AlertCircle className='h-8 w-8 text-red-400' />
          <span className='mt-2 text-2xl font-semibold'>{open_issues_count}</span>
          <span className='text-sm text-gray-500'>Open Issues</span>
        </div>
        <div className='flex flex-col items-center'>
          <Clock className='h-8 w-8 text-purple-400' />
          <span className='mt-2 text-2xl font-semibold'>{daysAgo(updated_at)} days</span>
          <span className='text-sm text-gray-500'>Last Update</span>
        </div>
      </CardContent>
    </Card>
  )
}

function daysAgo(dateString: string) {
  const date = new Date(dateString)
  const currentDate = new Date()
  const timeDiff = currentDate.getTime() - date.getTime()
  const daysDiff = Math.floor(timeDiff / (1000 * 3600 * 24))
  return daysDiff
}

// const DayAgo = ({ dateString }: { dateString: string }) => {
//   const days = daysAgo(dateString)
//   return (
//     <div className='flex gap-x-2 items-center'>
//       <CalendarIcon />
//       <span className='font-thin'>Update {days} days ago</span>
//     </div>
//   )
// }
