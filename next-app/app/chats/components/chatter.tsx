'use client'

import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'
import { PaperPlaneIcon } from '@radix-ui/react-icons'
import clsx from 'clsx'
import { useChat, useChatAlert } from '@/hooks/chat'
import { useRef, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import ChatAlert from './chatAlert'
import { useUser } from '@/context'
import Responser from './responser'

export default function Chatter() {
  const { isLoading, isDisabled, messages, ref, text, onFoucs, onBlur, sendMessage, handleKeyDown, handleTextValue } =
    useChat()

  const { presence } = useChatAlert()

  const messagesEndRef = useRef<HTMLDivElement>(null)

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  useEffect(() => {
    scrollToBottom()
  }, [messages.length])

  return (
    <div className='flex-1 flex flex-col justify-between p-4 w-full overflow-y-hidden'>
      <AnimatePresence>
        {presence && (
          <motion.div initial={{ opacity: 1, x: 0 }} exit={{ opacity: 0, x: 100, transition: { duration: 0.5 } }}>
            <ChatAlert />
          </motion.div>
        )}
      </AnimatePresence>
      <div className='flex-1 py-4 overflow-y-auto gap-y-6 flex flex-col w-full lg:w-3/4 xl:2/3 lg:mx-auto'>
        {messages.map((message, i) => (
          <div key={`chat-message-group-${i}`} className='flex flex-col gap-y-8'>
            <UserRequest key={`user-request-msg-${i}`} message={message} />
            <Responser key={`stargazer-response-msg-${i}`} query={message} />
          </div>
        ))}
        <div ref={messagesEndRef}></div>
      </div>
      <div className='relative w-full lg:w-3/4 xl:2/3 lg:mx-auto'>
        <Textarea
          placeholder='Message to StarGazer'
          className='pr-20 lg:text-base'
          ref={ref}
          onFocus={onFoucs}
          onBlur={onBlur}
          onKeyDown={handleKeyDown}
          value={text}
          onChange={handleTextValue}
          disabled={isDisabled || isLoading}
        />
        <Button
          className='absolute right-3 bottom-3 flex items-center'
          loading={isLoading}
          disabled={isDisabled}
          onClick={sendMessage}
        >
          <PaperPlaneIcon />
        </Button>
      </div>
    </div>
  )
}

const UserRequest = ({ message }: { message: string }) => {
  const { image } = useUser()

  return (
    <div className='flex flex-col gap-y-2 items-end w-full mx-auto'>
      <div className='flex gap-x-3 items-center justify-end'>
        <div className='text-slate-500'>You</div>
        <img src={image} alt='stargazer' className='rounded-full w-7 h-7' />
      </div>
      <div
        className={clsx(
          'rounded-lg shadow-md w-full flex gap-x-3 xl:w-2/3',
          'border-[1px] border-slate-100/40 p-3 text-wrap',
          'text-slate-700 border-[1px] border-slate-300',
        )}
      >
        {message}
      </div>
    </div>
  )
}
