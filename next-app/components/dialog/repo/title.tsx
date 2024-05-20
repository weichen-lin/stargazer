'use client'

import { DialogTitle } from '@/components/ui/dialog'
import { IRepoDetail } from '@/actions/neo4j'
import { CheckCheck, XCircle } from 'lucide-react'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { cn } from '@/lib/utils'

export default function Title(props: IRepoDetail) {
  const { avatar_url, full_name, is_vectorized } = props

  return (
    <DialogTitle className='flex gap-x-2 items-center'>
      <img width={48} height={48} src={avatar_url} alt={full_name} />
      <span className='ml-2 lg:text-2xl text-normal max-w-[220px] truncate'>{full_name}</span>
      <TooltipProvider>
        <Tooltip>
          <TooltipTrigger asChild>
            {is_vectorized ? (
              <CheckCheck className='w-4 h-4 text-green-500 ml-4' />
            ) : (
              <XCircle className='w-4 h-4 text-red-500 ml-4' />
            )}
          </TooltipTrigger>
          <TooltipContent
            className={cn('bg-white text-black border-[1px]', is_vectorized ? 'text-green-500' : 'text-red-500')}
          >
            {is_vectorized ? 'This repository is vectorized' : 'This repository is not yet vectorized'}
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>
    </DialogTitle>
  )
}
