'use server'

import Sidebar from '@/components/sidebar'
import Image from 'next/image'
import Link from 'next/link'

export default async function HomeLayout(props: { children: React.ReactNode; path: string }) {
  const { children } = props

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
        <Sidebar />
      </div>
      <div className='flex-1 h-full overflow-y-auto'>{children}</div>
    </div>
  )
}
