import { motion } from 'framer-motion'
import clsx from 'clsx'
import { useRepoDetail } from '@/hooks/util'
import { IRepository } from '@/client/repository'
import { getLanguageColor, Language } from '@/app/dashboard/components/pie-chart/config'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Star, Eye, LucideCalendarRange, ExternalLink, DoorOpen } from 'lucide-react'
import { useState } from 'react'
import { Checkbox } from '@/components/ui/checkbox'
import { useStars } from '@/hooks/stars'

export default function GridRepo(props: IRepository & { index: number }) {
  const {
    stargazers_count,
    repo_name,
    language,
    index,
    html_url,
    avatar_url,
    owner_name,
    open_issues_count,
    repo_id,
    updated_at,
  } = props
  const { setOpen, setRepoID } = useRepoDetail()
  const [isHover, setIsHover] = useState<boolean>(false)
  const { selectedRepo, setSelectedRepo } = useStars()

  return (
    <motion.div
      initial={{ opacity: 0, y: -100 }}
      animate={{ opacity: 1, y: 0, transition: { delay: index * 0.05 } }}
      className={clsx(
        'shadow-md bg-white dark:bg-slate-700/30 w-[365px] p-3 flex flex-col rounded-md h-[170px]',
        'border-[1px] dark:border-slate-100/30 border-slate-500/10 mx-auto relative md:overflow-hidden',
      )}
      onHoverStart={() => setIsHover(true)}
      onHoverEnd={() => setIsHover(false)}
    >
      <div className='flex flex-col justify-between gap-y-3 h-full'>
        <div className='flex gap-x-1 w-full'>
          <img src={avatar_url} alt={repo_name} width={60} height={60} className='rounded-full' />
          <div className='flex flex-col px-1 w-[calc(100%-64px)]'>
            <div className='flex justify-between items-center'>
              <div className='text-lg truncate w-[calc(100%-48px)]'>
                {owner_name}/{repo_name}
              </div>
              <Checkbox
                checked={selectedRepo.includes(repo_id)}
                onChange={(checked: boolean) => {
                  if (checked) {
                    setSelectedRepo([...selectedRepo, repo_id])
                  } else {
                    setSelectedRepo(selectedRepo.filter(id => id !== repo_id))
                  }
                }}
              />
            </div>
            <div className='rounded-full text-slate-500/75 dark:text-white/70 h-6'>@{owner_name}</div>
          </div>
        </div>
        <div className='flex justify-between items-center mb-4'>
          <Badge
            className={clsx('px-2 py-1 text-xs font-semibold rounded-full')}
            style={{ backgroundColor: getLanguageColor(language as Language) }}
          >
            {language}
          </Badge>
          <div className='flex gap-x-2 md:hidden'>
            <a href={html_url} target='_blank'>
              <Button size='sm' variant='outline'>
                <ExternalLink className='w-4 h-4' />
              </Button>
            </a>
            <Button
              size='sm'
              variant='outline'
              onClick={() => {
                setRepoID(repo_id)
                setOpen(true)
              }}
            >
              <DoorOpen className='w-4 h-4' />
            </Button>
          </div>
          <DayAgo dateString={updated_at} />
        </div>
        <div className='grid grid-cols-2'>
          <Issues count={open_issues_count} />
          <Stars count={stargazers_count} />
        </div>
      </div>
      {isHover && (
        <motion.div
          initial={{ opacity: 0, y: 170, scale: 0.9 }}
          animate={{ opacity: 1, y: 90, scale: 1, transition: { duration: 0.2 } }}
          className={clsx(
            'absolute inset-0 bg-opacity-20 flex items-center justify-center gap-x-4 w-full h-20 backdrop-blur-sm',
            'from-slate-200/40 to-slate-100/10 dark:from-slate-700 dark:to-slate-800 bg-gradient-to-t',
          )}
        >
          <a href={html_url} target='_blank'>
            <Button size='sm' variant='outline'>
              <ExternalLink className='w-4 h-4' />
            </Button>
          </a>
          <Button
            size='sm'
            variant='outline'
            onClick={() => {
              setRepoID(repo_id)
              setOpen(true)
            }}
          >
            <DoorOpen className='w-4 h-4' />
          </Button>
        </motion.div>
      )}
    </motion.div>
  )
}

const Issues = ({ count }: { count: number }) => {
  return (
    <div className='flex gap-x-2 items-center justify-start w-full'>
      <Eye className='h-5 w-5 text-red-400' />
      <span className='text-slate-700 dark:text-white/70 font-light'>{count} open issues</span>
    </div>
  )
}

const Stars = ({ count }: { count: number }) => {
  return (
    <div className='flex gap-x-2 items-center justify-start w-full'>
      <Star className='h-5 w-5 text-yellow-400' />
      <span className='text-slate-700 dark:text-white/70 font-light'>{count}</span>
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
    <div className='flex gap-x-2 items-center justify-end'>
      <LucideCalendarRange className='h-5 w-5 text-gray-500' />
      <span className='dark:text-white/70 text-sm text-gray-500'>{days} days ago</span>
    </div>
  )
}
