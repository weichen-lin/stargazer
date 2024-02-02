'use client'
import { motion } from 'framer-motion'
import { Star, Chat, Wrench } from '@phosphor-icons/react'
import clsx from 'clsx'
import Link from 'next/link'

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
  )
}
