import Image from 'next/image'
import LanguageMap from './color'
import { Star, GithubLogo } from '@phosphor-icons/react'
import { motion } from 'framer-motion'
import clsx from 'clsx'

export default function ListRepo(props: typeof githubStar & { index: number }) {
  const [owner, repo] = props.full_name.split('/')
  const { stargazers_count, language, description, html_url, index } = props
  return (
    <motion.div
      initial={{ opacity: 0, y: -100 }}
      animate={{ opacity: 1, y: 0, transition: { delay: index * 0.05 } }}
      whileHover={{ scale: 1.01 }}
      className={clsx('shadow-md w-full p-2 flex rounded-md cursor-pointer justify-between')}
    >
      <Image
        src={githubStar.owner.avatar_url}
        alt={githubStar.full_name}
        width={60}
        height={60}
        className='rounded-full'
      />
      <div className='flex flex-col justify-center items-start w-1/4 px-4'>
        <span className='text-md truncate w-full'>{owner}</span>
        <span className='font-semibold text-2xl truncate w-full'>{repo}</span>
      </div>
      <div className='text-slate-500/75 w-2/3 truncate p-3'>{description}</div>
      <div className='flex gap-x-8 justify-end items-center pr-6'>
        <Language language={language} />
        <Stars count={stargazers_count} />
        <a className='p-2 rounded-full hover:bg-slate-300/30' href={html_url} target='_blank'>
          <GithubLogo size={18} />
        </a>
      </div>
    </motion.div>
  )
}

const Language = ({ language }: { language: string }) => {
  const langColor = Object.keys(LanguageMap).includes(language) ? LanguageMap[language] : 'black'
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
      <Star size={18} weight='thin' />
      <span className='text-slate-700'>{count}</span>
    </div>
  )
}

const githubStar = {
  id: 107505869,
  full_name: 'firecracker-microvm/firecracker',
  owner: {
    avatar_url: 'https://avatars.githubusercontent.com/u/44477506?v=4',
  },
  html_url: 'https://github.com/firecracker-microvm/firecracker',
  description: 'Secure and fast microVMs for serverless computing.',
  updated_at: '2024-02-02T04:50:52Z',
  homepage: 'http://firecracker-microvm.io',
  stargazers_count: 23457,
  language: 'Rust',
}
