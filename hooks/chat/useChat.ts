import { useState, useRef, KeyboardEvent } from 'react'
import search from '@/actions/semantic'

type MessageKeyType = 'question' | 'suggest' | 'error'

export interface SuggestionProps {
  avatar_url: string
  html_url: string
  full_name: string
}

interface MessageValueType {
  question: string
  suggest: SuggestionProps[]
  error: string
}

type IMassages = {
  [K in MessageKeyType]: {
    type: K
    value: MessageValueType[K]
  }
}[MessageKeyType]

export default function useChat() {
  const [messages, setMessages] = useState<IMassages[]>([])
  const [isLoading, setIsloading] = useState<boolean>(false)
  const [isDisabled, setIsDisabled] = useState<boolean>(false)
  const [text, setText] = useState<string>('')

  const _Search = async (q: string) => {
    const res = await search(q)
    return res
  }

  const ref = useRef<HTMLTextAreaElement>(null)

  const onFoucs = () => {
    setIsDisabled(false)
  }

  const onBlur = () => {
    setIsDisabled(false)
  }

  const handleTextValue = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
    setText(event.target.value)
  }

  const handleKeyDown = async (event: KeyboardEvent<HTMLTextAreaElement>) => {
    if (event.key === 'Enter' && event.shiftKey) {
      return
    } else if (event.key === 'Enter') {
      // 阻止 Enter 键的默认行为（即换行）
      event.preventDefault()
      setText('')
      addMessage(text)
      const res = (await _Search(text)) as SuggestionProps[]

      if (res.length === 0) {
        setMessages(prev => [...prev, { type: 'error', value: '没有找到相关的项目' }])
      } else {
        setMessages(prev => [...prev, { type: 'suggest', value: res }])
      }
    }
  }

  const addMessage = (message: string) => {
    setMessages(prev => [...prev, { type: 'question', value: message }])
  }

  return { messages, addMessage, isLoading, isDisabled, ref, text, onFoucs, onBlur, handleKeyDown, handleTextValue }
}
