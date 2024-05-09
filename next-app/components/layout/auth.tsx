'use server'

import Sidebar, { MobileBar, DesktopBar } from '@/components/sidebar/star'
import Image from 'next/image'
import Link from 'next/link'
import { UserProvider } from '@/context/user'
import { z } from 'zod'
import { GetUser } from '@/actions'

const pageSchema = z.object({
  name: z.string().min(1),
  email: z.string().email(),
  image: z.string(),
})

export default async function AuthLayout(props: { children: React.ReactNode; path: string }) {
  const { path, children } = props

  const user = await GetUser()

  return (
    <UserProvider user={user}>
      <div className='h-screen w-screen flex overflow-x-hidden'>
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
          <Sidebar path={path} />
        </div>
        <div className='flex-1 flex flex-col h-full pt-[60px] lg:pt-0 overflow-y-auto dark:bg-slate-900'>
          <MobileBar />
          <DesktopBar />
          {children}
        </div>
      </div>
    </UserProvider>
  )
}
