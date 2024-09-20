'use client'

import { IRepository } from '@/client/repository'
import { StarIcon, EyeOpenIcon, CalendarIcon, GitHubLogoIcon } from '@radix-ui/react-icons'

export default function Info(props: IRepository) {
  const { stargazers_count, open_issues_count, last_modified_at, html_url } = props

  return (
    <div className='grid grid-cols-2 lg:grid-cols-4'>
      <a className='m-1' href={html_url} target='_blank'>
        <GitHubLogoIcon />
      </a>
      <div className='flex gap-x-2 items-center'>
        <EyeOpenIcon className='w-4 h-4' />
        <span>{open_issues_count}</span>
      </div>
      <div className='flex gap-x-2 items-center'>
        <StarIcon className='w-4 h-4' />
        <span>{stargazers_count}</span>
      </div>
      <DayAgo dateString={last_modified_at} />
    </div>
  )
}

function daysAgo(dateString: string) {
  const date = new Date(dateString)
  const currentDate = new Date()
  const timeDiff = currentDate.getTime() - date.getTime()
  const daysDiff = Math.floor(timeDiff / (1000 * 3600 * 24))
  return daysDiff
}

const DayAgo = ({ dateString }: { dateString: string }) => {
  const days = daysAgo(dateString)
  return (
    <div className='flex gap-x-2 items-center'>
      <CalendarIcon />
      <span className='font-thin'>Update {days} days ago</span>
    </div>
  )
}
