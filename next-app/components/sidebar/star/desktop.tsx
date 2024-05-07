'use client'

import { ModeToggle } from '@/components/provider'
import { motion } from 'framer-motion'
import { ArrangeSetting, FixPagination } from '@/components/tab'
import type { Arrangements } from '@/hooks/stars'
import clsx from 'clsx'
import { SyncStars } from '@/components/repo'
import { ChatSettingDialog } from '@/components/util/chatSetting'
import { Search } from '@/components/search'

interface SheetBarProps {
  total: number
  arrangement: string
  toggleArrangement: (arr: Arrangements) => void
}

const DesktopBar = () => {
  // const { total, arrangement, toggleArrangement } = props

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
        <div className='text-xl lg:text-3xl font-semibold w-[200px] pl-3'>My Stars</div>
        <div className='flex items-center gap-x-4'>
          <div className='flex items-center gap-x-4'>
            <Search />
            {/* <FixPagination total={total} />
            <ArrangeSetting arrangement={arrangement} toggle={toggleArrangement} /> */}
          </div>
          <ChatSettingDialog />
          <SyncStars />
          <ModeToggle />
        </div>
      </div>
    </motion.div>
  )
}

export default DesktopBar
