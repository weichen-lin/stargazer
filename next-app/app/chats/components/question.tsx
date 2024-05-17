'use client'

import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'
import { PaperPlaneIcon, GitHubLogoIcon } from '@radix-ui/react-icons'
import clsx from 'clsx'
import { useChat, useChatAlert } from '@/hooks/chat'
import { useRef, useEffect, useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import ChatAlert from './chatAlert'
import { ISuggestion } from '@/actions'
import { useUser } from '@/context'

const UserRequest = ({ message }: { message: string }) => {
  const { image } = useUser()

  return (
    <div className='flex flex-col gap-y-2 justify-end w-full lg:w-2/3 lg:mx-auto'>
      <div className='flex gap-x-3 items-center justify-end'>
        <div className='text-slate-500'>You</div>
        <img src={image} alt='stargazer' className='rounded-full w-7 h-7' />
      </div>
      <div
        className={clsx(
          'rounded-lg shadow-md w-full flex gap-x-3',
          'border-[1px] border-slate-100/40 p-3 text-wrap',
          'text-slate-700 border-[1px] border-slate-300',
        )}
      >
        {message}
      </div>
    </div>
  )
}
