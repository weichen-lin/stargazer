import { create } from 'zustand'

interface ISearchCollection {
  open: boolean
  setOpen: (open: boolean) => void
}

const searchCollectionState = create<ISearchCollection>(set => ({
  open: true,
  setOpen: open => set({ open }),
}))

export default function useSearchCollection() {
  const { open, setOpen } = searchCollectionState()

  return {
    open,
    setOpen,
  }
}
