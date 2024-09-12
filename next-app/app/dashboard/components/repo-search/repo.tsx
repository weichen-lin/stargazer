import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Star, ExternalLink, MessageSquare, EyeIcon } from 'lucide-react'
import { IRepository } from '@/client/repository'

export default function Repo(props: IRepository) {
  const {
    repo_id,
    html_url,
    owner_name,
    avatar_url,
    open_issues_count,
    repo_name,
    language,
    stargazers_count,
    watchers_count,
  } = props

  return (
    <div key={repo_id} className='flex items-start space-x-4 p-4 bg-gray-50 rounded-xl'>
      <Avatar className='w-10 h-10'>
        <AvatarImage src={avatar_url} alt={repo_name} />
        <AvatarFallback>{repo_name.substring(0, 2).toUpperCase()}</AvatarFallback>
      </Avatar>
      <div className='flex-1 space-y-1'>
        <div className='flex items-center justify-between'>
          <h3 className='text-sm font-medium'>{`${owner_name}/${repo_name}`}</h3>
          <Badge variant='outline'>{language}</Badge>
        </div>
        <div className='flex items-center space-x-4 text-sm text-gray-500'>
          <span className='flex items-center'>
            <Star className='w-4 h-4 mr-1' />
            {stargazers_count}
          </span>
          <span className='flex items-center'>
            <MessageSquare className='w-4 h-4 mr-1' />
            {open_issues_count}
          </span>
          <span className='flex items-center'>
            <EyeIcon className='w-4 h-4 mr-1' />
            {watchers_count}
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
  )
}
