import { create } from 'zustand'

interface RepoDetailState {
  open: boolean
  setOpen: (open: boolean) => void
  repoID: number
  setRepoID: (repo: number) => void
}

const repoDetailState = create<RepoDetailState>(set => ({
  open: false,
  setOpen: open => set({ open }),
  repoID: 0,
  setRepoID: repoID => set({ repoID }),
}))

export default function useRepoDetail() {
  const { open, setOpen, repoID, setRepoID } = repoDetailState()

  return { open, setOpen, repoID, setRepoID }
}
