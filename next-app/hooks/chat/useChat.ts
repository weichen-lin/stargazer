import { useState, useRef, KeyboardEvent } from 'react'
import { create } from 'zustand'

interface IChatStore {
  isLoading: boolean
  isDisabled: boolean
  setIsLoading: (isLoading: boolean) => void
  setIsDisabled: (isDisabled: boolean) => void
}

export const useChatStore = create<IChatStore>(set => ({
  isLoading: false,
  isDisabled: false,
  setIsLoading: (isLoading: boolean) => set({ isLoading }),
  setIsDisabled: (isDisabled: boolean) => set({ isDisabled }),
}))

export default function useChat() {
  const { isLoading, isDisabled, setIsLoading, setIsDisabled } = useChatStore()

  const [messages, setMessages] = useState<string[]>(['test'])
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

  const sendMessage = () => {
    setText('')
    setMessages(prev => [...prev, text.trim()])
    setIsLoading(true)
  }

  const handleKeyDown = async (event: KeyboardEvent<HTMLTextAreaElement>) => {
    if (event.key === 'Enter' && event.shiftKey) {
      return
    } else if (event.key === 'Enter') {
      event.preventDefault()
      sendMessage()
    }
  }

  return {
    messages,
    isLoading,
    isDisabled,
    ref,
    text,
    sendMessage,
    onFoucs,
    onBlur,
    handleKeyDown,
    handleTextValue,
  }
}
