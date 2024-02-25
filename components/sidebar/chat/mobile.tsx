'use client'

import { Sheet, SheetContent, SheetTrigger } from '@/components/ui/sheet'
import { Button } from '@/components/ui/button'
import { HamburgerMenuIcon } from '@radix-ui/react-icons'
import Sidebar from '@/components/sidebar/star'
import clsx from 'clsx'
import { ChatSettingDrawer } from '@/components/util/chatSetting'
import { ModeToggle } from '@/components/provider'

export default function MobileSidebar() {
  return (
    <div className={clsx('lg:hidden w-full p-3 gap-y-6 border-b-[1px] border-slate-300 backdrop-blur-md')}>
      <div className='flex justify-between w-full items-center pr-3'>
        <Sheet key='left'>
          <SheetTrigger asChild>
            <Button variant='ghost'>
              <HamburgerMenuIcon />
            </Button>
          </SheetTrigger>
          <SheetContent side='left'>
            <Sidebar path='chats' />
          </SheetContent>
        </Sheet>
        <div className='flex-1 text-xl lg:text-3xl font-semibold w-[200px] pl-3'>Start Chat</div>
        <div className='flex gap-x-4'>
          <ModeToggle />
          <ChatSettingDrawer />
        </div>
      </div>
    </div>
  )
}
