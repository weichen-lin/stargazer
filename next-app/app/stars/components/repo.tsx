import clsx from 'clsx'
import Image from 'next/image'
import { GitHubLogoIcon } from '@radix-ui/react-icons'

export default function Repo() {
  const data = {
    stargazers_count: 123,
    language: 'GO',
    description: 'go',
    html_url: 'https://github.com/weichen-lin/stargazer',
    full_name: 'stargazer',
    owner: 'asdasd',
    index: 1,
    id: '1213',
    homepage: 'https://github.com/weichen-lin/stargazer',
    avatar_url: 'https://avatars.githubusercontent.com/u/123123',
  }
  return (
    <div className='w-full flex justify-between px-3 bg-slate-300/40 py-2'>
      <div className='flex gap-x-2 items-center'>
        <Image width={40} height={40} src='https://avatars.githubusercontent.com/u/123123' alt=''></Image>
        <div className='text-slate-700 max-w-[200px] truncate'>weichen-lin/stargazer</div>
      </div>
      <div className='flex gap-x-2 items-center justify-center'>
        <div className='text-sm pt-1'>2024/05/11</div>
        <a className='rounded-full hover:bg-slate-300/30' href={data.html_url} target='_blank'>
          <GitHubLogoIcon />
        </a>
      </div>
    </div>
  )
}
