import LanguageMap from './color'
import { GitHubLogoIcon, StarIcon } from '@radix-ui/react-icons'
import { motion } from 'framer-motion'
import clsx from 'clsx'
import Tilt from 'react-parallax-tilt'
import type { Repository } from '@/actions/neo4j'

const config = {
  scale: 1.05,
  tiltMaxAngleX: 8,
  tiltMaxAngleY: 8,
  perspective: 1400,
  glareEnable: true,
  glareMaxOpacity: 0.1,
}

export default function GridRepo(props: Repository & { index: number }) {
  const [repoOwner, repo] = props.full_name.split('/')
  const { stargazers_count, language, description, index, html_url, full_name, owner } = props
  return (
    <Tilt {...config}>
      <motion.div
        initial={{ opacity: 0, y: -100 }}
        animate={{ opacity: 1, y: 0, transition: { delay: index * 0.05 } }}
        className={clsx(
          'shadow-md bg-white dark:bg-slate-700/30 w-[365px] p-3 flex flex-col rounded-md cursor-pointer',
          'border-[1px] dark:border-slate-100/30 border-slate-500/10',
        )}
      >
        <div className='flex flex-col justify-between gap-y-1'>
          <div className='flex gap-y-4 justify-between items-center w-full'>
            <div className='flex flex-col w-full'>
              <div className='text-md pr-2 truncate w-[280px]'>{repoOwner}/</div>
              <div className='font-semibold text-2xl w-[280px] pr-2 truncate'>{repo}</div>
            </div>
            <div className='w-20 h-20'>
              <img src={owner.avatar_url} alt={full_name} width={80} height={80} className='rounded-full' />
            </div>
          </div>
          <div className='text-slate-500/75 w-full h-28 overflow-y-auto truncate text-wrap dark:text-white/70'>
            {description}
          </div>
          <div className='flex gap-x-8 justify-end w-full items-center'>
            <Language language={language} />
            <Stars count={stargazers_count} />
            <a className='p-2 rounded-full hover:bg-slate-300/30' href={html_url} target='_blank'>
              <GitHubLogoIcon />
            </a>
          </div>
        </div>
      </motion.div>
    </Tilt>
  )
}

const Language = ({ language }: { language: string }) => {
  return (
    <div className='flex gap-x-2 items-center justify-center'>
      <div
        className='h-2.5 w-2.5 rounded-full'
        style={{
          backgroundColor: LanguageMap[language],
        }}
      ></div>
      <span className='text-slate-700 dark:text-white/70'>{language}</span>
    </div>
  )
}

const Stars = ({ count }: { count: number }) => {
  return (
    <div className='flex gap-x-2 items-center justify-center'>
      <StarIcon />
      <span className='text-slate-700 dark:text-white/70'>{count}</span>
    </div>
  )
}
