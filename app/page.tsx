'use client'

import Image from 'next/image'
import clsx from 'clsx'
import { NextUIProvider } from '@nextui-org/system'
import { Input } from '@nextui-org/input'
import { Button } from '@nextui-org/button'
import { ArrowLongRightIcon } from '@heroicons/react/16/solid'

export default function Home() {
  return (
    <NextUIProvider>
      <main className='h-screen p-6 bg-[#f2f0e8] overflow-hidden'>
        <div className={clsx('w-full max-w-[1024px] mx-auto h-full', 'flex flex-col justify-between')}>
          <div className='flex items-center gap-x-8 justify-start'>
            <Image src='/icon.png' width={60} height={60} className='rounded-full' alt='stargazer logo' />
            <div className='text-gray-700 text-4xl'>StarGazer</div>
          </div>
          <div className='flex justify-between items-center flex-col md:flex-row'>
            <div className='flex flex-col gap-y-8'>
              <div className='text-gray-700 text-2xl md:text-4xl'>Start Managing Your Stars Today.</div>
              <div className='text-gray-700 text-md md:text-lg'>
                AI is reshaping Star Management beyond Github Stars. Explore limitless potential in AI-driven star
                management – just the start!
              </div>
            </div>
            <Image src='/home.png' width={400} height={400} className='' alt='home pic' />
          </div>
          <div className='w-[320px] mx-auto flex gap-x-4 items-center md:my-36'>
            <Input
              key='username'
              type='text'
              label='Username'
              labelPlacement='inside'
              radius='md'
              classNames={{
                label: 'text-black',
                input: ['text-lg', 'text-black/90', 'placeholder:text-black dark:placeholder:text-white/60'],
                innerWrapper: '',
                inputWrapper: ['shadow-xl', 'border-[2px] border-slate-700/30'],
              }}
            />
            <Button
              isIconOnly
              size='lg'
              color='secondary'
              spinner={
                <svg
                  className='animate-spin h-5 w-5 text-current'
                  fill='none'
                  viewBox='0 0 24 24'
                  xmlns='http://www.w3.org/2000/svg'
                >
                  <circle className='opacity-25' cx='12' cy='12' r='10' stroke='currentColor' strokeWidth='4' />
                  <path
                    className='opacity-75'
                    d='M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z'
                    fill='currentColor'
                  />
                </svg>
              }
            >
              <ArrowLongRightIcon className='w-6 h-6' />
            </Button>
          </div>
          <div className='text-black text-center'>© WeiChen Lin 2024</div>
        </div>
      </main>
    </NextUIProvider>
  )
}
