'use client'

import { ModeToggle } from '@/components/provider'
import { motion } from 'framer-motion'
import { ArrangeSetting, FixPagination } from '@/components/tab'
import type { Arrangements } from '@/hooks/stars'
import clsx from 'clsx'

interface SheetBarProps {
  total: number
  arrangement: string
  toggleArrangement: (arr: Arrangements) => void
}

const DesktopBar = (props: SheetBarProps) => {
  const { total, arrangement, toggleArrangement } = props

  return (
    <motion.div
      initial={{ x: 80 }}
      animate={{ x: 0 }}
      className={clsx(
        'flex-col items-center justify-between fixed z-10 left-[260px] hidden lg:flex',
        'w-[calc(100%-260px)] p-3 gap-y-6 border-b-[1px] border-slate-300 backdrop-blur-md',
      )}
    >
      <div className='flex justify-between w-full items-center'>
        <div className='text-xl lg:text-3xl font-semibold w-[200px] pl-3'>My Stars</div>
        <div className='flex items-center gap-x-12'>
          <div className='flex items-center gap-x-12'>
            <FixPagination total={total} />
            <ArrangeSetting arrangement={arrangement} toggle={toggleArrangement} />
          </div>
          <ModeToggle />
        </div>
      </div>
    </motion.div>
  )
}

export default DesktopBar
