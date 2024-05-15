'use client'

import { useStars } from '@/hooks/stars'
import { Loading, Empty } from './status'
import Repo from './repo'
import { Detail } from '@/components/dialog'
import { useRepoDetail } from '@/hooks/util'

export default function Results() {
  const { isSearching, results } = useStars()
  const { open } = useRepoDetail()

  if (isSearching) {
    return <Loading />
  }

  if (results.length === 0) {
    return <Empty />
  }

  return (
    <div
      className={`w-full h-full md:grid md:grid-cols-2 lg:grid-cols-4 gap-4 py-4 flex flex-col items-center justify-center`}
    >
      {results.map((item, index) => (
        <Repo key={item.repo_id} {...item} index={index} />
      ))}
      {open && <Detail />}
    </div>
  )
}
