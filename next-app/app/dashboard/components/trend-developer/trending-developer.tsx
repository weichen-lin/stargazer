import { Button } from '@/components/ui/button'
import { ExternalLink } from 'lucide-react'
import { RocketIcon } from '@radix-ui/react-icons'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { ITrendDeveloper } from '@/client/trends/type'

export default function TrendDeveloper(props: ITrendDeveloper) {
  const { avatar_url, name, sub_name, html_url, repo_name, description } = props

  return (
    <div className='flex flex-col p-1 bg-gray-50 rounded-xl border-[1px] border-slate-300'>
      <div className='w-full flex justify-between items-center'>
        <div className='flex items-center space-x-4 p-2 bg-gray-50 rounded-xl w-full'>
          <Avatar className='w-10 h-10'>
            <AvatarImage src={avatar_url} alt={repo_name} />
            <AvatarFallback>{repo_name.substring(0, 2).toUpperCase()}</AvatarFallback>
          </Avatar>
          <div className='flex flex-col justify-center'>
            <span className='text-blue-500 text-lg'>{name}</span>
            <span className='text-sm text-slate-500'>{sub_name}</span>
          </div>
        </div>
        <a href={html_url} target='_blank' className='mr-4'>
          <Button size='sm' variant='outline'>
            <ExternalLink className='w-4 h-4' />
          </Button>
        </a>
      </div>
      <div className='flex flex-col p-2 gap-y-2'>
        <div className='flex gap-x-2 items-center'>
          <RocketIcon className='w-5 h-5' />
          <div className='text-lg'>{repo_name}</div>
        </div>
        <div className='text-gray-500 text-sm'>{description}</div>
      </div>
    </div>
  )
}
