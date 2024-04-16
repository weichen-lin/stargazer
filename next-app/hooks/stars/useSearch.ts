import { create } from 'zustand'
import { useState } from 'react'

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
  const [repos, setRepos] = useState<ISearchRepo[]>([testRepo])

  return {
    query,
    open,
    repos,
    setOpen,
    setQuery,
  }
}
