import LanguageMap from './color'
import { GitHubLogoIcon, StarIcon } from '@radix-ui/react-icons'
import { motion } from 'framer-motion'
import clsx from 'clsx'
import type { Repository } from '@/actions/neo4j'

export default function ListRepo(props: Repository & { index: number }) {
  const [repoOwner, repo] = props.full_name.split('/')
  const { stargazers_count, language, description, html_url, full_name, owner, index } = props
  return (
    <motion.div
      initial={{ opacity: 0, y: -100 }}
      animate={{ opacity: 1, y: 0, transition: { delay: index * 0.05 } }}
      whileHover={{ scale: 1.01 }}
      className={clsx(
        'shadow-md w-full p-2 flex rounded-md cursor-pointer justify-start md:justify-between',
        'dark:border-slate-100/30 dark:border-b-[1px]',
      )}
    >
      <div className='flex'>
        <img src={owner.avatar_url} alt={full_name} width={60} height={60} className='rounded-full' />
        <div className='flex flex-col justify-center items-start w-[280px] px-4'>
          <span className='text-md truncate w-full'>{repoOwner}</span>
          <span className='font-semibold text-2xl truncate w-full'>{repo}</span>
        </div>
      </div>
      <div className='text-slate-500/75 flex-1 truncate px-3 py-1 hidden lg:block text-wrap max-h-[64px]'>
        {description}
      </div>
      <div className='md:flex gap-x-8 justify-end items-center pr-6 hidden'>
        <Language language={language} />
        <Stars count={stargazers_count} />
        <a className='p-2 rounded-full hover:bg-slate-300/30' href={html_url} target='_blank'>
          <GitHubLogoIcon />
        </a>
      </div>
    </motion.div>
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
      <span className='text-slate-700'>{language}</span>
    </div>
  )
}

const Stars = ({ count }: { count: number }) => {
  return (
    <div className='flex gap-x-2 items-center justify-center'>
      <StarIcon />
      <span className='text-slate-700'>{count}</span>
    </div>
  )
}
