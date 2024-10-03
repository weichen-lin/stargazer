'use client'

import { motion } from 'framer-motion'
import clsx from 'clsx'
import Link from 'next/link'
import { StarIcon, ExitIcon, DashboardIcon, MixIcon } from '@radix-ui/react-icons'
import { Button } from '@/components/ui/button'
import { signOut } from 'next-auth/react'
import { useUser } from '@/context/user'
import { usePathname } from 'next/navigation'

const Bars = [
  { name: 'Dashboard', icon: <DashboardIcon />, path: 'dashboard' },
  { name: 'My Stars', icon: <StarIcon />, path: 'stars' },
  { name: 'Collections', icon: <MixIcon />, path: 'collections' },
]

export default function Sidebar() {
  const { name, image } = useUser()
  const pathname = usePathname()

  return (
    <div className='flex flex-col h-full justify-between w-[260px]'>
      <div className='flex flex-col gap-y-4 justify-between h-full pb-6'>
        <div className='flex flex-col gap-y-4'>
          {Bars.map((e, i) => {
            const isCurrent = pathname.includes(e.path)
            return (
              <motion.div whileTap={{ scale: 0.9 }} key={i}>
                <Link
                  href={`/${e.path}`}
                  className={clsx(
                    'py-2 w-[240px] cursor-pointer flex items-center justify-between px-4',
                    `${isCurrent ? 'bg-slate-300/40' : 'hover:bg-slate-300/40'}`,
                  )}
                >
                  <div className='flex items-center justify-between gap-x-4'>
                    <div>{e.icon}</div>
                    <div>{e.name}</div>
                  </div>
                  {isCurrent && (
                    <span className='animate-[ping_1.5s_ease-in-out_infinite] inline-flex h-1.5 w-1.5 rounded-full bg-sky-400 opacity-75'></span>
                  )}
                </Link>
              </motion.div>
            )
          })}
        </div>
        <Button
          variant='secondary'
          className='flex border-[1px] border-slate-300 gap-x-2 max-w-[150px] mx-auto'
          onClick={() => signOut()}
        >
          <ExitIcon className='rotate-180' />
          Log Out
        </Button>
      </div>
      <div className='flex flex-col py-2 border-t-[1px] border-slate-500/10'>
        <div className='flex items-center justify-start my-4 gap-x-4 pl-4 hover:bg-slate-300/40 cursor-pointer py-2'>
          <img src={image} alt='avatar' width={40} height={40} className='rounded-full' />
          <div>{name}</div>
        </div>
        <div
          className={clsx('py-2 w-[240px] cursor-pointer flex items-center justify-between px-4 hover:bg-slate-300/40')}
        ></div>
        <p className='text-slate-500 text-sm text-center pt-4 border-t-[1px] border-slate-500/10'>WeiChen Lin © 2024</p>
      </div>
    </div>
  )
}
