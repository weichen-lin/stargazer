'use client'

import { IRepository } from '@/client/repository'
import { Tag } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'

export default function Info(props: IRepository) {
  const { description, topics } = props

  return (
    <Card>
      <CardHeader>
        <CardTitle>About</CardTitle>
        <CardDescription>Repository description and details</CardDescription>
      </CardHeader>
      <CardContent>
        <p className='text-gray-700 mb-4'>{description}</p>
        <div className='flex items-center space-x-2 mb-4'>
          <Tag className='h-5 w-5 text-gray-400' />
          <span className='text-sm font-semibold text-gray-700'>Topics:</span>
        </div>
        <div className='flex flex-wrap gap-2'>
          {topics?.map((topic, index) => (
            <Badge key={index} variant='secondary'>
              {topic}
            </Badge>
          ))}
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
