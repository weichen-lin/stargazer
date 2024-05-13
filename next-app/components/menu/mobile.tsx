'use client'

import { Sheet, SheetContent, SheetTrigger } from '@/components/ui/sheet'
import { Button } from '@/components/ui/button'
import { HamburgerMenuIcon } from '@radix-ui/react-icons'
import Sidebar from '@/components/sidebar'
import { ModeToggle } from '@/components/provider'
import clsx from 'clsx'
import { SyncStars } from '@/components/repo'
import { ChatSettingDialog } from '@/components/util/chatSetting'
import { useMenuName } from './util'

const MobileBar = () => {
  const { menuName } = useMenuName()

  return (
    <div
      className={clsx(
        'flex flex-col items-center justify-between',
        'fixed top-0 left-0 z-10 bg-white dark:bg-black',
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
            <Sidebar />
          </SheetContent>
        </Sheet>
        <div className='text-xl lg:text-3xl font-semibold w-[200px] pl-3'>{menuName}</div>
        <div className='flex items-center gap-x-2'>
          <SyncStars />
          <ModeToggle />
          <ChatSettingDialog />
        </div>
      </div>
    </div>
  )
}

export default MobileBar
