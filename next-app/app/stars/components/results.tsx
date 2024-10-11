'use client'

import { useStars } from '@/hooks/stars'
import { Loading, Empty } from './status'
import Repo from './repo'
import SearchCollection from './dialog'
import { useSearchCollection } from '@/app/stars/hook'

export default function Results() {
  const { isSearching, results } = useStars()
  const { open } = useSearchCollection()

  if (isSearching) {
    return <Loading />
  }

  if (!results || results.length === 0) {
    return <Empty />
  }

  return (
    <div className='w-full h-full md:grid md:grid-cols-2 xl:grid-cols-3 3xl:grid-cols-4 gap-4 py-4 flex flex-col content-start'>
      {results.map((item, index) => (
        <Repo key={item.repo_id} {...item} index={index} />
      ))}
      {open && <SearchCollection />}
    </div>
  )
}
