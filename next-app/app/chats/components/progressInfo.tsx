'use client'

import dynamic from 'next/dynamic'
import { Progress } from '@/components/ui/progress'
import { useEffect, useState } from 'react'
import { useChatAlert } from '@/hooks/chat'
import { EventSourcePolyfill } from 'event-source-polyfill'
import { generateAccessToken } from '@/actions/util'
import { useUser } from '@/context'
import { z } from 'zod'

const PushUp = dynamic(() => import('@/components/fancyicon/push-up'), { ssr: false })
const Check = dynamic(() => import('@/components/fancyicon/check'), { ssr: false })
const Error = dynamic(() => import('@/components/fancyicon/error'), { ssr: false })

const schema = z.object({
  total: z.number().default(0),
  current: z.number().default(0),
})

type ISchema = z.infer<typeof schema>

const parseEventData = (data: string): ISchema => {
  const o = JSON.parse(data)
  return schema.parse(o)
}

const producer_url = process.env.NEXT_PUBLIC_PRODUCER_URL ?? '/producer'

export default function ProgressInfo() {
  const [isEstablish, setIsEstablish] = useState(true)
  const [isFinished, setIsFinished] = useState(false)
  const [status, setStatus] = useState<ISchema>({ total: 0, current: 0 })
  const { startEvent, setCantClose } = useChatAlert()
  const { email } = useUser()

  useEffect(() => {
    let eventSource: EventSourcePolyfill | null = null

    const startSyncStars = async () => {
      try {
        const token = await generateAccessToken(email)
        eventSource = new EventSourcePolyfill(`${producer_url}/sync_user_stars`, {
          headers: {
            Authorization: token,
          },
        })

        setIsEstablish(false)
        setCantClose(false)

        eventSource.onmessage = event => {
          const o = parseEventData(event.data)
          setStatus(prev => ({ ...prev, ...o }))
        }

        eventSource.onerror = e => {
          eventSource?.close()
          setIsFinished(true)
          setCantClose(true)

          setTimeout(() => {
            startEvent(false)
          }, 1500)
        }
      } catch (e) {
        return null
      }
    }

    startSyncStars()

    return () => {
      eventSource && eventSource?.close()
    }
  }, [])

  return (
    <div className='w-full max-w-[450px]'>
      {isEstablish && <Confirming />}
      {!isEstablish && !isFinished && <Progressing {...status} />}
      {isFinished && <CheckMark />}
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

const Progressing = (props: ISchema) => {
  const { current, total } = props
  console.log({ current, total })
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
