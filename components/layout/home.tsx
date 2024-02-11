import Sidebar from '@/components/sidebar/star'
import Image from 'next/image'
import Link from 'next/link'

export default async function HomeLayout(props: { username: string; children: React.ReactNode; path: string }) {
  const { username, path, children } = props

  return (
    <div className='h-screen w-screen flex'>
      <div className='pl-4 lg:flex flex-col w-[240px] hidden items-center'>
        <Link href='/' className='py-8'>
          <Image
            src='/icon.png'
            alt='Logo'
            width={50}
            height={50}
            className='rounded-full cursor-pointer hover:opacity-75'
          />
        </Link>
        <div className='w-full h-[1px] border-b-[1px] border-slate-700/10 mb-8'></div>
        <Sidebar path={path} username={username} />
      </div>
      <div className='flex-1 h-full overflow-y-auto'>{children}</div>
    </div>
  )
}
