'use client'

import { ModeToggle } from '@/components/provider'
import { motion } from 'framer-motion'
import { ArrangeSetting, FixPagination } from '@/components/tab'
import type { Arrangements } from '@/hooks/stars'
import clsx from 'clsx'
import { SyncStars } from '@/components/repo'
import { ChatSettingDialog } from '@/components/util/chatSetting'
import Search from '@/components/search/search'

interface SheetBarProps {
  total: number
  arrangement: string
  toggleArrangement: (arr: Arrangements) => void
}

const DesktopBar = (props: SheetBarProps) => {
  const { total, arrangement, toggleArrangement } = props

  // useEffect(() => {
  //   const handleKeyDown = (event: KeyboardEvent) => {
  //     event.preventDefault()
  //     // 如果是Mac系統
  //     if (navigator.userAgent.toUpperCase().indexOf('MAC') >= 0) {
  //       if (event.metaKey && event.key === 'p') {
  //         // Command + P 被按下
  //         console.log('Mac: Command + P 被按下')
  //         // 在此處放置你想要執行的動作
  //       }
  //     } else {
  //       // 如果是Windows系統
  //       if (event.ctrlKey && event.key === 'p') {
  //         // Ctrl + P 被按下
  //         console.log('Windows: Ctrl + P 被按下')
  //         // 在此處放置你想要執行的動作
  //       }
  //     }
  //   }

  //   document.addEventListener('keydown', handleKeyDown)

  //   return () => {
  //     document.removeEventListener('keydown', handleKeyDown)
  //   }
  // }, [])

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
        <div className='flex items-center gap-x-4'>
          <div className='flex items-center gap-x-4'>
            <Search />
            <FixPagination total={total} />
            <ArrangeSetting arrangement={arrangement} toggle={toggleArrangement} />
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
