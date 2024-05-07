'use server'

import Sidebar, { MobileBar, DesktopBar } from '@/components/sidebar/star'
import Image from 'next/image'
import Link from 'next/link'
import { getServerSession } from 'next-auth'
import { options } from '@/app/api/auth/[...nextauth]/option'
import { UserProvider } from '@/context/user'
import { z } from 'zod'
import { redirect } from 'next/navigation'

const pageSchema = z.object({
  name: z.string().min(1),
  email: z.string().email(),
  image: z.string(),
})

export default async function AuthLayout(props: { children: React.ReactNode; path: string }) {
  const { path, children } = props

  try {
    const session = await getServerSession(options)
    if (!session) {
      redirect('/')
    }

    const user = pageSchema.parse(session.user)
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
          <div className='flex-1 flex flex-col pt-20 lg:pt-0 overflow-y-auto bg-blue-50/70 dark:bg-slate-900'>
            <MobileBar />
            <DesktopBar />
            {children}
          </div>
        </div>
      </UserProvider>
    )
  } catch {
    redirect('/')
  }
}
