'use client'

import Image from 'next/image'
import clsx from 'clsx'
import { Button } from '@/components/ui/button'
import { signIn } from 'next-auth/react'
import { GitHubLogoIcon, HomeIcon } from '@radix-ui/react-icons'
import { ModeToggle } from '@/components/provider'
import { useSession } from 'next-auth/react'
import dynamic from 'next/dynamic'

const TypeWritter = dynamic(() => import('@/components/ui/typewriter-effect'), { ssr: false })

export default function Home() {
  const words = [
    {
      text: 'Start',
      className: 'text-xl lg:text-3xl',
    },
    {
      text: 'managing',
      className: 'text-xl lg:text-3xl',
    },
    {
      text: 'your',
      className: 'text-xl lg:text-3xl',
    },
    {
      text: 'stars',
      className: 'text-xl lg:text-3xl',
    },
    {
      text: 'with',
      className: 'text-xl lg:text-3xl',
    },
    {
      text: 'StarGazer.',
      className: 'text-xl text-blue-500 dark:text-blue-500 lg:text-3xl',
    },
  ]

  const { data: session } = useSession()

  const handleSignIn = async () => {
    const result = await signIn('github', { callbackUrl: '/stars' })
    if (result?.error) {
      console.error('Sign in failed:', result.error)
    }
  }

  return (
    <main className='h-screen p-6 overflow-hidden'>
      <div className={clsx('w-full max-w-[1024px] mx-auto h-full', 'flex flex-col justify-between')}>
        <div className='flex items-center justify-between'>
          <div className='flex gap-x-8 items-center'>
            <Image src='/icon.png' width={60} height={60} className='rounded-full' alt='stargazer logo' />
            <div className='text-4xl'>StarGazer</div>
          </div>
          <ModeToggle />
        </div>
        <div className='flex justify-between items-start md:items-center flex-col lg:flex-row gap-6'>
          <div className='flex flex-col'>
            <TypeWritter words={words} cursorClassName='h-7 lg:h-10' />
            <p className='leading-7'>
              AI is reshaping Star Management beyond Github Stars. Explore limitless potential in AI-driven star
              management – just the start!
            </p>
          </div>
          <Image src='/home.png' width={400} height={400} className='mx-auto' alt='home pic' />
        </div>
        {session ? (
          <Button className='flex gap-x-4 max-w-[320px] mx-auto' onClick={handleSignIn}>
            <HomeIcon />
            Dashboard
          </Button>
        ) : (
          <Button className='flex gap-x-4 max-w-[320px] mx-auto' onClick={handleSignIn}>
            <GitHubLogoIcon />
            Sign in with Github
          </Button>
        )}
        <div className='text-center'>© WeiChen Lin 2024</div>
      </div>
    </main>
  )
}
