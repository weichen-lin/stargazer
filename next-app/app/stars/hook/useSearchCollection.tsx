import { create } from 'zustand'
import { useFetch } from '@/hooks/util'

interface ISearchCollection {
  open: boolean
  setOpen: (open: boolean) => void
  chosen: string | null
  setChosen: (chosen: string | null) => void
}

const searchCollectionState = create<ISearchCollection>(set => ({
  open: false,
  setOpen: open => set({ open }),
  chosen: null,
  setChosen: chosen => set({ chosen }),
}))

export default function useSearchCollection() {
  const { open, setOpen, chosen, setChosen } = searchCollectionState()

  return {
    open,
    setOpen,
    chosen,
    setChosen,
  }
}
