import { create } from 'zustand'

interface ChatAlertProps {
  presence: boolean
  close: () => void
}

const useChatAlert = create<ChatAlertProps>(set => ({
  presence: true,
  close: () => set({ presence: false }),
}))

export default useChatAlert
