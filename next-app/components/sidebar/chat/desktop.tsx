'use client'

import { motion } from 'framer-motion'
import clsx from 'clsx'
import { ChatSettingDialog } from '@/components/util/chatSetting'
import { ModeToggle } from '@/components/provider'

const DesktopBar = () => {
  return (
    <div
      className={clsx(
        'items-center justify-between z-10 left-[260px] hidden lg:block w-full',
        'p-3 border-b-[1px] border-slate-300 backdrop-blur-md',
      )}
    >
      <div className='flex justify-between w-full items-center'>
        <motion.div initial={{ x: 80 }} animate={{ x: 0 }} className='text-xl lg:text-3xl font-semibold w-[200px] pl-3'>
          Start Chat
        </motion.div>
        <div className='flex gap-x-4'>
          <ModeToggle />
          <ChatSettingDialog />
        </div>
      </div>
    </div>
  )
}

export default DesktopBar
