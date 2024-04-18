import { Sheet, SheetContent, SheetTrigger } from '@/components/ui/sheet'
import { Button } from '@/components/ui/button'
import { HamburgerMenuIcon } from '@radix-ui/react-icons'
import Sidebar from '@/components/sidebar/star'
import { ModeToggle } from '@/components/provider'
import { motion } from 'framer-motion'
import { ArrangeSetting, FixPagination } from '@/components/tab'
import type { Arrangements } from '@/hooks/stars'
import clsx from 'clsx'
import { SyncStars } from '@/components/repo'
import { ChatSettingDialog } from '@/components/util/chatSetting'

interface SheetBarProps {
  total: number
  arrangement: string
  toggleArrangement: (arr: Arrangements) => void
}

const MobileBar = (props: SheetBarProps) => {
  const { total, arrangement, toggleArrangement } = props

  return (
    <div
      className={clsx(
        'flex flex-col items-center justify-between',
        'fixed top-0 left-0 z-10',
        'lg:hidden w-full p-3 gap-y-6 border-b-[1px] border-slate-300 backdrop-blur-md',
      )}
    >
      <div className='flex justify-between w-full items-center'>
        <Sheet key='left'>
          <SheetTrigger asChild>
            <Button variant='ghost'>
              <HamburgerMenuIcon />
            </Button>
          </SheetTrigger>
          <SheetContent side='left'>
            <Sidebar path='stars' />
          </SheetContent>
        </Sheet>
        <div className='text-xl lg:text-3xl font-semibold w-[200px] pl-3'>My Stars</div>
        <div className='flex items-center gap-x-2'>
          <ArrangeSetting arrangement={arrangement} toggle={toggleArrangement} />
          <SyncStars />
          <ModeToggle />
          <ChatSettingDialog />
        </div>
      </div>
      {total > 0 && (
        <motion.div initial={{ x: 80 }} animate={{ x: 0 }} className='flex justify-between items-center w-full'>
          <FixPagination total={total} />
        </motion.div>
      )}
    </div>
  )
}

export default MobileBar
