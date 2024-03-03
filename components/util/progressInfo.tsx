'use client'

import dynamic from 'next/dynamic'
import { Progress } from '@/components/ui/progress'
import { useEffect, useState } from 'react'
import { useChatAlert } from '@/hooks/chat'

const PushUp = dynamic(() => import('@/components/fancyicon/push-up'), { ssr: false })
const Check = dynamic(() => import('@/components/fancyicon/check'), { ssr: false })
const Error = dynamic(() => import('@/components/fancyicon/error'), { ssr: false })

export default function ProgressInfo() {
  const [isEstablish, setIsEstablish] = useState(true)
  const [isFinished, setIsFinished] = useState(false)
  const [errorCode, setErrorCode] = useState<number>(200)
  const [current, setCurrent] = useState(0)
  const [total, setTotal] = useState(0)
  const { startEvent, setCantClose } = useChatAlert()

  useEffect(() => {
    const eventSource = new EventSource('/api/sse')

    eventSource.onmessage = event => {
      setIsEstablish(false)

      const msg = JSON.parse(event.data)
      if (msg.error) {
        setErrorCode(msg.status)
        setCantClose(false)
        eventSource.close()
        setTimeout(() => {
          startEvent(false)
        }, 3500)
      } else {
        if (msg.current === msg.total) {
          setIsFinished(true)
          setCantClose(false)
          setTimeout(() => {
            startEvent(false)
          }, 3500)
          eventSource.close()
        } else {
          setCantClose(true)
          setCurrent(msg.current + 1)
          setTotal(msg.total)
        }
      }
    }

    eventSource.onerror = () => {
      setCantClose(false)
      setIsFinished(true)
      eventSource.close()
      setTimeout(() => {
        startEvent(false)
      }, 2000)
    }

    return () => {
      eventSource.close()
    }
  }, [])

  return (
    <div className='w-full max-w-[450px]'>
      {isEstablish && <Confirming />}
      {!isEstablish && !isFinished && errorCode < 203 && <Progressing current={current} total={total} />}
      {isFinished && <CheckMark />}
      {errorCode > 202 && <ErrorMsg status={errorCode} />}
    </div>
  )
}

const Confirming = () => {
  return (
    <div className='flex justify-center items-center w-full gap-x-4'>
      <PushUp />
      <div className='pt-2'>Confirming the vectorize status ...</div>
    </div>
  )
}

const CheckMark = () => {
  return (
    <div className='flex justify-center items-center w-full gap-x-4'>
      <Check />
      <div className='pt-1'>Vectorize is complete!</div>
    </div>
  )
}

const MsgMap: { [key: number]: string } = {
  404: 'Invalid OpenAI API key.',
}

const ErrorMsg = ({ status }: { status: number }) => {
  const msg = MsgMap[status] || 'Error occurred. Please try again later.'

  return (
    <div className='flex justify-center items-center w-full gap-x-4'>
      <Error />
      <div className='pt-1'>{msg}</div>
    </div>
  )
}

const Progressing = (props: { current: number; total: number }) => {
  const { current, total } = props

  const value = total === 0 ? 0 : (current / total) * 100
  return (
    <div className='w-full flex justify-center gap-x-2 items-center'>
      <Progress value={value} max={100} className='w-1/2' />
      <div className='flex gap-x-2 mx-12'>
        <span>{current}</span>
        <span>/</span>
        <span>{total}</span>
      </div>
    </div>
  )
}
