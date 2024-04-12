import { create } from 'zustand'

interface ChatAlertProps {
  presence: boolean
  cantClose: boolean
  isStartEvent: boolean
  startEvent: (s: boolean) => void
  close: () => void
  setCantClose: (c: boolean) => void
}

const useChatAlert = create<ChatAlertProps>(set => ({
  presence: true,
  cantClose: false,
  isStartEvent: false,
  startEvent: (s: boolean) => set(prev => ({ ...prev, isStartEvent: s })),
  close: () => set(prev => ({ ...prev, presence: false })),
  setCantClose: (cantclose: boolean) => set(prev => ({ ...prev, cantClose: cantclose })),
}))

export default useChatAlert
