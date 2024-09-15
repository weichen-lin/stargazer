import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Star, ExternalLink, MessageSquare, StarsIcon } from 'lucide-react'
import { IRepository } from '@/client/repository'

interface ITrendRepository {
  repo_name: string
  owner_name: string
  html_url: string
  description: string
  stargazers_count: number
  language: string
  get_stars: number
  date_range: DateRange
}

export type DateRange = 'daily' | 'weekly' | 'monthly'

export const sinceMap: { [key in DateRange]: string } = {
  daily: 'today',
  weekly: 'this week',
  monthly: 'this month',
}

export default function TrendReposiory(props: ITrendRepository) {
  const { html_url, owner_name, repo_name, description, language, stargazers_count, get_stars, date_range } = props

  return (
    <div className='flex flex-col p-1 bg-gray-50 rounded-xl border-[1px] border-slate-300'>
      <div className='flex items-start space-x-4 p-2 bg-gray-50 rounded-xl w-full'>
        <div className='flex-1 space-y-1'>
          <div className='flex items-center justify-between'>
            <h3 className='text-md font-bold'>{`${owner_name}/${repo_name}`}</h3>
            <Badge variant='outline'>{language}</Badge>
          </div>
          <div className='flex items-center space-x-4 text-sm text-gray-500'>
            <span className='flex items-center'>
              <Star className='w-4 h-4 mr-1' />
              {stargazers_count}
            </span>
          </div>
        </div>
        <div className='flex flex-col space-y-2'>
          <a href={html_url} target='_blank'>
            <Button size='sm' variant='outline'>
              <ExternalLink className='w-4 h-4' />
            </Button>
          </a>
        </div>
      </div>
      <div className='text-gray-500 text-sm p-2'>{description}</div>
      <div className='flex items-center justify-start gap-x-1 p-1 text-slate-700 text-sm'>
        <StarsIcon className='w-4 h-4 mr-1' />
        <span>{get_stars}</span>
        <span>{sinceMap[date_range]}</span>
      </div>
    </div>
  )
}
