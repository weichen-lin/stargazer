'use client'

import { DialogTitle } from '@/components/ui/dialog'
import { IRepository } from '@/client/repository'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { ExternalLink } from 'lucide-react'

export default function Title(props: IRepository) {
  const { avatar_url, repo_name, owner_name, html_url } = props

  return (
    <DialogTitle className='flex gap-x-4 items-center w-full justify-between'>
      <div className='flex gap-x-3 items-center'>
        <Avatar className='w-10 h-10'>
          <AvatarImage src={avatar_url} alt={repo_name} />
          <AvatarFallback>{repo_name.substring(0, 2).toUpperCase()}</AvatarFallback>
        </Avatar>
        <div className='text-xl'>
          {owner_name}/{repo_name}
        </div>
      </div>
      <a href={html_url} target='_blank'>
        <Button size='sm' variant='outline'>
          <ExternalLink className='w-4 h-4' />
        </Button>
      </a>
    </DialogTitle>
  )
}
