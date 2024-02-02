import Image from 'next/image'
import LanguageMap from './color'
import { Star, GithubLogo } from '@phosphor-icons/react'

export default function GridRepo(props: typeof githubStar) {
  const [owner, repo] = props.full_name.split('/')
  const { stargazers_count, language, description } = props
  return (
    <div className='shadow-md bg-white w-[420px] p-6 flex flex-col rounded-md'>
      <div className='flex flex-col justify-between'>
        <div className='flex gap-y-4 justify-between items-center w-full'>
          <div className='flex flex-col'>
            <span className='text-md'>{owner}/</span>
            <span className='font-semibold text-2xl'>{repo}</span>
          </div>
          <div>
            <Image
              src={githubStar.owner.avatar_url}
              alt={githubStar.full_name}
              width={80}
              height={80}
              className='rounded-full'
            />
          </div>
        </div>
        <div className='text-slate-500/75 w-[calc(100%-80px)] h-28 overflow-y-auto'>{description}</div>
        <div className='flex gap-x-8 justify-start w-full items-center'>
          <Language language={language} />
          <Stars count={stargazers_count} />
          <GithubLogo size={18} />
        </div>
      </div>
    </div>
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
