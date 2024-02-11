'use client'

import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'
import { PaperPlaneIcon, GitHubLogoIcon } from '@radix-ui/react-icons'
import clsx from 'clsx'
import { useChat, SuggestionProps } from '@/hooks/chat'
import { useRef, useEffect } from 'react'

export default function Chatter() {
  const { isLoading, isDisabled, messages, ref, text, onFoucs, onBlur, handleKeyDown, handleTextValue } = useChat()

  const messagesEndRef = useRef<HTMLDivElement>(null)

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  return (
    <div className='h-full flex flex-col justify-between p-4 w-full'>
      <div className='flex-1 py-2 overflow-y-auto gap-y-6 flex flex-col'>
        {messages.map((message, i) => {
          if (message.type === 'question') {
            return <UserRequest key={i} message={message.value} />
          } else if (message.type === 'suggest') {
            return <StarGazerResponse key={i} suggestions={message.value} />
          } else if (message.type === 'error') {
            return <UserRequest key={i} message={message.value} />
          }
        })}
        <div ref={messagesEndRef}></div>
      </div>
      <div className='relative lg:w-2/3 lg:mx-auto'>
        <Textarea
          placeholder='Message to StarGazer'
          className='pr-20 lg:text-base'
          ref={ref}
          onFocus={onFoucs}
          onBlur={onBlur}
          onKeyDown={handleKeyDown}
          value={text}
          onChange={handleTextValue}
          disabled={isDisabled}
        />
        <Button className='absolute right-3 bottom-3 flex items-center' loading={isLoading} disabled={isDisabled}>
          <PaperPlaneIcon />
        </Button>
      </div>
    </div>
  )
}

const StarGazerResponse = (props: { suggestions: SuggestionProps[] }) => {
  const { suggestions } = props
  return (
    <div className='w-full lg:w-2/3 lg:mx-auto'>
      <div className='flex flex-col items-start gap-y-2 w-full sm:max-w-[360px]'>
        <div className='flex gap-x-3 items-center'>
          <img src='/icon.png' alt='stargazer' className='rounded-full w-7 h-7' />
          <div className='font-semibold'>StarGazer</div>
        </div>
        <div className='px-3.5 py-3 flex flex-col w-full gap-y-5'>
          {suggestions.map((e, i) => (
            <Suggestions key={i} {...e} />
          ))}
        </div>
      </div>
    </div>
  )
}

const UserRequest = ({ message }: { message: string }) => {
  return (
    <div className='flex justify-end w-full lg:w-2/3 lg:mx-auto'>
      <div className='flex flex-col items-end gap-y-2 w-full sm:max-w-[480px]'>
        <div className='flex gap-x-3 items-center justify-end'>
          <div className='font-semibold'>You</div>
          <img src='/icon.png' alt='stargazer' className='rounded-full w-7 h-7' />
        </div>
        <div className='px-3.5 py-3 flex flex-col w-full gap-y-5'>
          <div
            className={clsx(
              'rounded-md drop-shadow-md shadow-md w-full flex gap-x-3 items-center justify-between',
              'border-[1px] border-slate-100/40 p-3 text-wrap',
              'bg-blue-100 text-slate-700',
            )}
          >
            {message}
          </div>
        </div>
      </div>
    </div>
  )
}

const Suggestions = (props: SuggestionProps) => {
  const { avatar_url, full_name, html_url } = props
  return (
    <div className='rounded-md drop-shadow-md shadow-md w-full flex gap-x-3 items-center justify-between border-[1px] border-slate-100/40 p-2'>
      <div className='flex gap-x-2 items-center'>
        <div className='w-10 h-10'>
          <img src={avatar_url} alt='stargazer' className='rounded-full w-full h-full' />
        </div>
        <div>{full_name}</div>
      </div>
      <a className='p-2 rounded-full hover:bg-slate-300/30' href={html_url} target='_blank'>
        <GitHubLogoIcon />
      </a>
    </div>
  )
}
