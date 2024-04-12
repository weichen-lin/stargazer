'use client'

import { GridIcon, ListBulletIcon } from '@radix-ui/react-icons'
import { motion } from 'framer-motion'
import { Arrangements } from '@/hooks/stars'
import clsx from 'clsx'

interface ArrangeMentsProps {
  arrangement: string
  toggle: (arrangement: Arrangements) => void
}

export default function ArrangeSetting(props: ArrangeMentsProps) {
  const { arrangement, toggle } = props
  return (
    <div className='flex gap-x-1 lg:gap-x-4 justify-between items-center bg-slate-300/20 shadow-md p-1 lg:p-2 rounded-md'>
      <motion.button
        whileTap={{ scale: arrangement === 'grid' ? 1 : 0.9 }}
        className={clsx('text-gray-700 flex gap-x-3 py-1 px-2 rounded-md', `${arrangement === 'grid' && 'bg-white'}`)}
        onClick={() => toggle('grid')}
      >
        <GridIcon />
      </motion.button>
      <motion.button
        whileTap={{ scale: arrangement === 'list' ? 1 : 0.9 }}
        className={clsx('text-gray-700 flex gap-x-3 py-1 px-2 rounded-md', `${arrangement === 'list' && 'bg-white'}`)}
        onClick={() => toggle('list')}
      >
        <ListBulletIcon />
      </motion.button>
    </div>
  )
}
