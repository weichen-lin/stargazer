import { create } from 'zustand'
import { useState, useRef, useEffect } from 'react'
import { GetFullTextSearch, ISuggestion } from '@/actions'

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
  const [loading, setLoading] = useState(false)
  const [repos, setRepos] = useState<ISuggestion[]>([])
  const ref = useRef<HTMLInputElement>(null)
  const timeoutRef = useRef<NodeJS.Timeout | null>(null)

  const queryRepos = async (query: string) => {
    setLoading(true)
    const repos = await GetFullTextSearch(query)
    setRepos(repos)
    setLoading(false)
  }

  const openDialog = () => {
    setOpen(true)
    if (ref) {
      ref.current?.focus()
    }
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
        queryRepos(query)
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
    repos,
    ref,
    loading,
    setOpen,
    setQuery,
    queryRepos,
  }
}
