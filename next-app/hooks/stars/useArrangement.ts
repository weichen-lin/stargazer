import { useState } from 'react'

export type Arrangements = 'list' | 'grid'

export default function useArrangement() {
  const [arrangement, setArrangement] = useState<Arrangements>('grid')

  const toggleArrangement = (arr: Arrangements) => {
    setArrangement(arr)
  }

  return { arrangement, toggleArrangement }
}
