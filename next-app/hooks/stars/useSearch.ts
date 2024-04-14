import { create } from 'zustand'

interface ISearch {
  query: string
  open: boolean
  setOpen: (open: boolean) => void
  setQuery: (query: string) => void
}

const searchAtom = create<ISearch>(set => ({
  query: '',
  open: false,
  setOpen: open => set({ open }),
  setQuery: query => set({ query }),
}))

export default function useSearch() {
  const { query, open, setOpen, setQuery } = searchAtom()

  return {
    query,
    open,
    setOpen,
    setQuery,
  }
}
