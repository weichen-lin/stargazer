'use server'

import Menu from '@/components/menu'
import Sidebar from '@/components/sidebar'
import Image from 'next/image'
import Link from 'next/link'
import { UserProvider } from '@/context/user'
import { GetUser } from '@/actions'

export default async function AuthLayout(props: { children: React.ReactNode }) {
  const { children } = props
  const user = await GetUser()

  return (
    <UserProvider user={user}>
      <div className='h-screen w-screen flex flex-col lg:flex-row overflow-x-hidden'>
        <div className='pl-4 lg:flex flex-col w-[260px] hidden items-center'>
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
        <div className='flex-1 flex flex-col pt-[60px] lg:pt-0 overflow-y-auto dark:bg-slate-900'>
          <Menu />
          {children}
        </div>
      </div>
    </UserProvider>
  )
}
