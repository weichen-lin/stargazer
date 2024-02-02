'use client'

import { motion } from 'framer-motion'
import { Star, Chat, Wrench } from '@phosphor-icons/react'
import clsx from 'clsx'
import Link from 'next/link'
import Image from 'next/image'

const Bars = [
  { name: 'My Stars', icon: <Star size={24} weight='light' />, path: 'stars' },
  { name: 'Start Chat', icon: <Chat size={28} weight='light' />, path: 'chats' },
  { name: 'Settings', icon: <Wrench size={28} weight='light' />, path: 'settings' },
]

interface SidebarProps {
  path: string
  username: string
}

export default function Sidebar(props: SidebarProps) {
  const { path, username } = props

  return (
    <div className='flex flex-col h-full justify-between w-[240px]'>
      <div className='flex flex-col gap-y-4'>
        {Bars.map((e, i) => {
          return (
            <motion.div whileTap={{ scale: 0.9 }} key={i}>
              <Link
                href={`/${e.path}/${username}`}
                className={clsx(
                  'py-2 w-[240px] cursor-pointer flex items-center justify-between px-4',
                  `${e.path === path ? 'bg-slate-300/40' : 'hover:bg-slate-300/40'}`,
                )}
              >
                <div className='flex items-center justify-between gap-x-4'>
                  <div>{e.icon}</div>
                  <div>{e.name}</div>
                </div>
                {e.path === path && (
                  <span className='animate-[ping_1.5s_ease-in-out_infinite] inline-flex h-1.5 w-1.5 rounded-full bg-sky-400 opacity-75'></span>
                )}
              </Link>
            </motion.div>
          )
        })}
      </div>
      <div className='flex flex-col py-2 border-t-[1px] border-slate-500/10'>
        <div className='flex items-center justify-start my-4 gap-x-4 pl-4 hover:bg-slate-300/40 cursor-pointer py-2'>
          <img
            src='https://i.pravatar.cc/150?u=a042581f4e29026024d'
            alt='avatar'
            width={40}
            height={40}
            className='rounded-full'
          />
          <div>{username}</div>
        </div>
        <p className='text-slate-500 text-sm text-center pt-4 border-t-[1px] border-slate-500/10'>WeiChen Lin © 2024</p>
      </div>
    </div>
  )
}
