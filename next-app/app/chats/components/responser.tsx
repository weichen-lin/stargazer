'use client'

import { useEffect, useState } from 'react'
import { GetSuggesions } from '@/actions'
import { ISuggestion } from '@/actions/neo4j/repos'
import { Error, HaveSuggestions } from './message'
import { useChatStore } from '@/hooks/chat'

const Responser = ({ query }: { query: string }) => {
  const [answers, setAnswers] = useState<ISuggestion[]>([])
  const [status, setStatus] = useState<number>(200)
  const [isLoading, setResponserLoading] = useState<boolean>(true)
  const { setIsLoading } = useChatStore()

  // useEffect(() => {
  //   const fetchData = async () => {
  //     setIsLoading(true)
  //     try {
  //       const res = await GetSuggesions(query)
  //       if (res.status === 200) {
  //         setAnswers(res.items)
  //       }
  //       setStatus(res.status)
  //     } catch (error) {
  //       console.error('Error fetching suggestions: ', error)
  //       setStatus(500)
  //     } finally {
  //       setIsLoading(false)
  //       setResponserLoading(false)
  //     }
  //   }

  //   fetchData()
  // }, [])

  return (
    <div className='w-full flex flex-col gap-y-4'>
      <div className='flex gap-x-3 items-center'>
        <img src='/icon.jpeg' alt='stargazer' className='rounded-full w-7 h-7' />
        <div className='text-slate-500'>StarGazer</div>
      </div>
      {isLoading && <div className='loader m-4'></div>}
      {!isLoading && answers.length > 0 && <HaveSuggestions suggestions={answers} />}
      {!isLoading && answers.length === 0 && <Error status={status} />}
    </div>
  )
}

export default Responser
