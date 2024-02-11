import { Sheet, SheetContent, SheetTrigger } from '@/components/ui/sheet'
import { Button } from '@/components/ui/button'
import { HamburgerMenuIcon } from '@radix-ui/react-icons'
import Sidebar from '@/components/sidebar/star'
import { ModeToggle } from '@/components/theme'
import { motion } from 'framer-motion'
import { ArrangeSetting, FixPagination } from '@/components/tab'
import type { Arrangements } from '@/hooks/stars'
import clsx from 'clsx'

interface SheetBarProps {
  total: number
  current: number
  arrangement: string
  toggleArrangement: (arr: Arrangements) => void
  page: string
}

const MobileBar = (props: SheetBarProps) => {
  const { total, current, arrangement, toggleArrangement, page } = props
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
            <Sidebar path='stars' username='username' />
          </SheetContent>
        </Sheet>
        <div className='text-xl lg:text-3xl font-semibold w-[200px] pl-3'>My Stars</div>
        <div className='flex items-center gap-x-6'>
          <ArrangeSetting arrangement={arrangement} toggle={toggleArrangement} />
          <ModeToggle />
        </div>
      </div>
      <motion.div initial={{ x: 80 }} animate={{ x: 0 }} className='flex justify-between items-center w-full'>
        <FixPagination total={total} current={current} page={page} />
      </motion.div>
    </div>
  )
}

export default MobileBar
