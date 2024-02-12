import * as React from 'react'
import { MinusIcon, PlusIcon, GearIcon } from '@radix-ui/react-icons'
import { Button } from '@/components/ui/button'
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/components/ui/drawer'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'

export function ChatSettingDialog() {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant='outline' size='icon'>
          <GearIcon className='h-4 w-4' />
        </Button>
      </DialogTrigger>
      <DialogContent className='sm:max-w-[425px]'>
        <DialogHeader>
          <DialogTitle>Edit profile</DialogTitle>
          <DialogDescription>Make changes to your profile here. Click save when you&apos;re done.</DialogDescription>
        </DialogHeader>
        <div className='grid gap-4 py-4'>
          <div className='grid grid-cols-4 items-center gap-4'></div>
        </div>
        <DialogFooter>
          <Button type='submit'>Save changes</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

export function ChatSettingDrawer() {
  const [goal, setGoal] = React.useState(350)

  function onClick(adjustment: number) {
    setGoal(Math.max(200, Math.min(400, goal + adjustment)))
  }

  return (
    <Drawer>
      <DrawerTrigger asChild>
        <GearIcon />
      </DrawerTrigger>
      <DrawerContent>
        <div className='mx-auto w-full max-w-sm'>
          <DrawerHeader>
            <DrawerTitle>Move Goal</DrawerTitle>
            <DrawerDescription>Set your daily activity goal.</DrawerDescription>
          </DrawerHeader>
          <div className='p-4 pb-0'>
            <div className='flex items-center justify-center space-x-2'>
              <Button
                variant='outline'
                size='icon'
                className='h-8 w-8 shrink-0 rounded-full'
                onClick={() => onClick(-10)}
                disabled={goal <= 200}
              >
                <MinusIcon className='h-4 w-4' />
                <span className='sr-only'>Decrease</span>
              </Button>
              <div className='flex-1 text-center'>
                <div className='text-7xl font-bold tracking-tighter'>{goal}</div>
                <div className='text-[0.70rem] uppercase text-muted-foreground'>Calories/day</div>
              </div>
              <Button
                variant='outline'
                size='icon'
                className='h-8 w-8 shrink-0 rounded-full'
                onClick={() => onClick(10)}
                disabled={goal >= 400}
              >
                <PlusIcon className='h-4 w-4' />
                <span className='sr-only'>Increase</span>
              </Button>
            </div>
            <div className='mt-3 h-[120px]'>123</div>
          </div>
          <DrawerFooter>
            <Button>Submit</Button>
            <DrawerClose asChild>
              <Button variant='outline'>Cancel</Button>
            </DrawerClose>
          </DrawerFooter>
        </div>
      </DrawerContent>
    </Drawer>
  )
}
