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

const testRepo = {
  avatar: 'https://avatars.githubusercontent.com/u/60848391?s=48&v=4',
  html_url: 'https://github.com/weichen-lin/stargazer',
  description: 'A GitHub repository search engine',
  summary:
    'A GitHub repository search engine A GitHub repository search engine A GitHub repository search engine A GitHub repository search engine',
  full_name: 'weichen-lin/stargazer',
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

    if (ref) {
      ref.current?.focus()
    }
  }

  useEffect(() => {
    if (query.length > 0) {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current)
      }

      timeoutRef.current = setTimeout(() => {
        queryRepos(query)
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
