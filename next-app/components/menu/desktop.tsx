'use client'

import { ModeToggle } from '@/components/provider'
import { motion } from 'framer-motion'
import clsx from 'clsx'
import { Search } from '@/components/search'
import { useMenuName } from './util'
import { Detail } from '@/components/dialog'
import { useRepoDetail } from '@/hooks/util'

const DesktopBar = () => {
  const { menuName } = useMenuName()
  const { open } = useRepoDetail()

  return (
    <motion.div
      initial={{ x: 80 }}
      animate={{ x: 0 }}
      className={clsx(
        'flex-col items-center justify-between hidden lg:flex',
        'p-3 gap-y-6 border-b-[1px] border-slate-300 backdrop-blur-md',
        'bg-white dark:bg-black dark:border-slate-800 dark:text-white',
      )}
    >
      <div className='flex justify-between w-full items-center'>
        <div className='text-xl lg:text-3xl font-semibold w-[200px] pl-3'>{menuName}</div>
        <div className='flex items-center gap-x-4'>
          <Search />
          <ModeToggle />
          {open && <Detail />}
        </div>
      </div>
    </motion.div>
  )
}

export default DesktopBar
