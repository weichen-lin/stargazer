import { create } from 'zustand'
import { useRef, useEffect } from 'react'
import { useFetch } from '@/hooks/util'
import { IRepository } from '@/client/repository'

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

export interface ISearchRepo {
  avatar: string
  html_url: string
  description: string
  summary: string
  full_name: string
}

export default function useSearch() {
  const { query, open, setOpen, setQuery } = searchAtom()
  const { data, isLoading, run } = useFetch<IRepository[]>({
    initialRun: false,
    config: {
      url: '/repository/full-text-search',
      method: 'GET',
    },
  })
  const ref = useRef<HTMLInputElement>(null)
  const timeoutRef = useRef<NodeJS.Timeout | null>(null)

  const openDialog = () => {
    setOpen(true)
    if (ref) {
      ref.current?.focus()
    }
  }

  const clear = () => {
    setQuery('')
  }

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        setOpen(false)
      }

      if (navigator.userAgent.toUpperCase().indexOf('MAC') >= 0) {
        if (event.metaKey && event.key === 'k') {
          openDialog()
        }
      } else {
        if (event.ctrlKey && event.key === 'k') {
          openDialog()
        }
      }
    }

    document.addEventListener('keydown', handleKeyDown)

    return () => {
      document.removeEventListener('keydown', handleKeyDown)
    }
  }, [])

  useEffect(() => {
    if (query.length > 0) {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current)
      }

      timeoutRef.current = setTimeout(() => {
        run({ params: { query } })
        if (ref) {
          ref.current?.focus()
        }
      }, 500)
    }
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current)
      }
    }
  }, [query])

  return {
    query,
    open,
    data,
    ref,
    clear,
    isLoading,
    setOpen,
    setQuery,
  }
}
