import { useState, useRef, KeyboardEvent } from 'react'
import { GetSuggesions, ISuggestion } from '@/actions'

type MessageKeyType = 'question' | 'suggest' | 'error'

interface MessageValueType {
  question: string
  suggest: ISuggestion[]
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

  const sendMessage = async (text: string) => {
    setText('')
    addMessage(text)
    setIsloading(true)
    const res = await GetSuggesions(text)
    setIsloading(false)
    return res
  }

  const handleButtonOnClick = async () => {
    const res = await sendMessage(text)

    if (res.length === 0) {
      setMessages(prev => [...prev, { type: 'error', value: '没有找到相关的项目' }])
    } else {
      setMessages(prev => [...prev, { type: 'suggest', value: res }])
    }
  }

  const handleKeyDown = async (event: KeyboardEvent<HTMLTextAreaElement>) => {
    if (event.key === 'Enter' && event.shiftKey) {
      return
    } else if (event.key === 'Enter') {
      event.preventDefault()
      const res = await sendMessage(text)

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

  return {
    messages,
    addMessage,
    handleButtonOnClick,
    isLoading,
    isDisabled,
    ref,
    text,
    onFoucs,
    onBlur,
    handleKeyDown,
    handleTextValue,
  }
}
