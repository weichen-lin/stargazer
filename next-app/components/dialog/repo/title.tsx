'use client'

import { DialogTitle } from '@/components/ui/dialog'
import { IRepoDetail } from '@/actions/neo4j'
import { CheckCheck, XCircle } from 'lucide-react'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { cn } from '@/lib/utils'
import { IRepository } from '@/client/repository'

export default function Title(props: IRepository) {
  const { avatar_url, repo_name } = props

  return (
    <DialogTitle className='flex gap-x-2 items-center'>
      <img width={48} height={48} src={avatar_url} alt={repo_name} />
      <span className='ml-2 lg:text-2xl text-normal max-w-[220px] truncate'>{repo_name}</span>
      <TooltipProvider>
        <Tooltip>
          {/* <TooltipTrigger asChild>
            {is_vectorized ? (
              <CheckCheck className='w-4 h-4 text-green-500 ml-4' />
            ) : (
              <XCircle className='w-4 h-4 text-red-500 ml-4' />
            )}
          </TooltipTrigger> */}
          {/* <TooltipContent
            className={cn('bg-white text-black border-[1px]', is_vectorized ? 'text-green-500' : 'text-red-500')}
          >
            {is_vectorized ? 'This repository is vectorized' : 'This repository is not yet vectorized'}
          </TooltipContent> */}
        </Tooltip>
      </TooltipProvider>
    </DialogTitle>
  )
}
