import { LayersIcon, LockClosedIcon, LockOpen2Icon } from '@radix-ui/react-icons'
import { MagicCard } from '@/components/ui/magic-card'
import { useTheme } from 'next-themes'
import { ICollection } from '@/client/collection'
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { MoreVertical } from 'lucide-react'
import { Button } from '@/components/ui/button'
import moment from 'moment'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'

export default function Collection(props: ICollection) {
  const { id, name, is_public, created_at, updated_at } = props
  const { theme } = useTheme()
  return (
    <MagicCard
      className='cursor-pointer h-[100px] w-full py-2'
      gradientColor={theme === 'dark' ? '#262626' : '#D9D9D955'}
    >
      <div className='flex flex-col justify-between items-center w-full h-full'>
        <div className='flex gap-x-2 px-3 justify-between w-full'>
          <div className='flex gap-x-2 items-center'>
            <LayersIcon className='w-5 h-5' />
            {name}
          </div>
          <div className='flex gap-x-2 items-center'>
            {is_public ? <LockOpen /> : <LockClose />}
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant='ghost' className='h-8 w-8 p-0'>
                  <span className='sr-only'>Open menu</span>
                  <MoreVertical className='h-5 w-5' />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align='end'>
                <DropdownMenuLabel>Actions</DropdownMenuLabel>
                <DropdownMenuItem>Copy payment ID</DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem>View customer</DropdownMenuItem>
                <DropdownMenuItem>View payment details</DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
        <div className='flex px-3 w-full justify-between'>
          <div className='text-slate-500/70'>{moment(updated_at).fromNow()}</div>
        </div>
      </div>
    </MagicCard>
  )
}

const LockOpen = () => {
  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <LockOpen2Icon className='w-5 h-5 text-green-500' />
        </TooltipTrigger>
        <TooltipContent>
          <p>
            This collection is <strong>public</strong>
          </p>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  )
}

const LockClose = () => {
  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <LockClosedIcon className='w-5 h-5 text-red-500' />
        </TooltipTrigger>
        <TooltipContent>
          <p>
            This collection is <strong>private</strong>
          </p>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  )
}
