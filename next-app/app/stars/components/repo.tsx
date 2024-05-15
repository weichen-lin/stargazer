import LanguageMap from './color'
import { GitHubLogoIcon, StarIcon, EyeOpenIcon, CalendarIcon, OpenInNewWindowIcon } from '@radix-ui/react-icons'
import { motion } from 'framer-motion'
import clsx from 'clsx'
import type { Repository } from '@/actions/neo4j'
import { useRepoDetail } from '@/hooks/util'

// const config = {
//   scale: 1.05,
//   tiltMaxAngleX: 8,
//   tiltMaxAngleY: 8,
//   perspective: 1400,
//   glareEnable: true,
//   glareMaxOpacity: 0.1,
// }

export default function GridRepo(props: Repository & { index: number }) {
  const {
    stargazers_count,
    language,
    index,
    html_url,
    full_name,
    avatar_url,
    owner_url,
    owner_name,
    open_issues_count,
    last_updated_at,
    repo_id,
  } = props
  const { setOpen, setRepoID } = useRepoDetail()

  return (
    <motion.div
      initial={{ opacity: 0, y: -100 }}
      animate={{ opacity: 1, y: 0, transition: { delay: index * 0.05 } }}
      className={clsx(
        'shadow-md bg-white dark:bg-slate-700/30 w-[365px] p-3 flex flex-col rounded-md h-[170px]',
        'border-[1px] dark:border-slate-100/30 border-slate-500/10 mx-auto',
      )}
    >
      <div className='flex flex-col justify-between gap-y-1'>
        <div className='flex gap-x-1 w-full'>
          <img src={avatar_url} alt={full_name} width={60} height={60} className='rounded-full' />
          <div className='flex flex-col'>
            <div className='text-xl truncate w-[280px]'>{full_name}</div>
            <div className='flex gap-x-[2px] items-center'>
              <a className='rounded-full hover:bg-slate-300/30 w-6 h-6 p-1' href={html_url} target='_blank'>
                <GitHubLogoIcon className='w-4 h-4' />
              </a>
              <div className='rounded-full text-slate-500/75 dark:text-white/70 h-6 px-1'>@{owner_name}</div>
              <div className='flex-1 flex justify-end'>
                <OpenInNewWindowIcon
                  className='w-4 h-4 cursor-pointer hover:text-slate-500 dark:hover:text-white/70'
                  onClick={() => {
                    setOpen(true)
                    setRepoID(repo_id)
                  }}
                />
              </div>
            </div>
          </div>
        </div>
        <div className='grid grid-cols-2 mt-8'>
          <Issues count={open_issues_count} />
          <Language language={language} />
          <Stars count={stargazers_count} />
          <DayAgo dateString={last_updated_at} />
        </div>
      </div>
    </motion.div>
  )
}

const Language = ({ language }: { language: string }) => {
  return (
    <div className='flex gap-x-2 items-center justify-end'>
      <div
        className='h-2.5 w-2.5 rounded-full'
        style={{
          backgroundColor: LanguageMap[language],
        }}
      ></div>
      <span className='text-slate-500 dark:text-white/70'>{language}</span>
    </div>
  )
}

const Issues = ({ count }: { count: number }) => {
  return (
    <div className='flex gap-x-2 items-center justify-start w-full'>
      <EyeOpenIcon />
      <span className='text-slate-700 dark:text-white/70 font-light'>{count} open issues</span>
    </div>
  )
}

const Stars = ({ count }: { count: number }) => {
  return (
    <div className='flex gap-x-2 items-center justify-start w-full'>
      <StarIcon />
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
      <CalendarIcon />
      <span className='text-slate-500 dark:text-white/70 font-thin'>Update {days} days ago</span>
    </div>
  )
}
