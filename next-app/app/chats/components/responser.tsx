'use client'

import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'
import { PaperPlaneIcon, GitHubLogoIcon } from '@radix-ui/react-icons'
import clsx from 'clsx'
import { useChat, useChatAlert } from '@/hooks/chat'
import { useRef, useEffect, useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import ChatAlert from './chatAlert'
import { ISuggestion, GetSuggesions } from '@/actions'
import { useUser } from '@/context'
import { Empty, Error, HaveSuggestions } from './message'
import { useChatStore } from '@/hooks/chat'
import { memo } from 'react'

const useSuggestions = (query: string) => {
  const [answers, setAnswers] = useState<ISuggestion[] | boolean>(false)
  const [isLoading, setResponserLoading] = useState<boolean>(true)
  const { setIsLoading } = useChatStore()

  useEffect(() => {
    const fetchData = async () => {
      setIsLoading(true)
      try {
        const res = await GetSuggesions(query)
        setAnswers(res)
      } catch (error) {
        console.log(error)
      } finally {
        setIsLoading(false)
        setResponserLoading(false)
      }
    }

    fetchData()
  }, [])

  return { answers, isLoading }
}

const Responser = ({ query }: { query: string }) => {
  const [answers, setAnswers] = useState<ISuggestion[] | boolean>(false)
  const [isLoading, setResponserLoading] = useState<boolean>(true)
  const { setIsLoading } = useChatStore()

  useEffect(() => {
    const fetchData = async () => {
      setIsLoading(true)
      try {
        const res = await GetSuggesions(query)
        setAnswers(res)
      } catch (error) {
        console.log(error)
      } finally {
        setIsLoading(false)
        setResponserLoading(false)
      }
    }

    fetchData()
  }, [])

  return (
    <div className='w-full flex flex-col gap-y-2 lg:w-2/3'>
      <div className='flex gap-x-3 items-center'>
        <img src='/icon.jpeg' alt='stargazer' className='rounded-full w-7 h-7' />
        <div className='text-slate-500'>StarGazer</div>
      </div>
      {isLoading && <div className='loader m-4'></div>}
      {!isLoading && Array.isArray(answers) && <HaveSuggestions suggestions={answers} />}
      {!isLoading && !Array.isArray(answers) && <Error />}
    </div>
  )
}

export default memo(Responser)
