'use client'

import { useArrangement } from '@/hooks/stars'
import { GridRepo, ListRepo } from '@/components/repo'
import clsx from 'clsx'
import { MobileBar, DesktopBar } from '@/components/sidebar/star'

interface Star {
  id: number
  full_name: string
  owner: {
    avatar_url: string
  }
  html_url: string
  description: string
  homepage: string
  stargazers_count: number
  language: string
}

export default function Stars({
  stars,
  total,
  current,
  page,
}: {
  stars: Star[]
  total: number
  current: number
  page: string
}) {
  const { arrangement, toggleArrangement } = useArrangement()
  return (
    <div className='flex flex-col gap-y-12 w-full mt-36 lg:mt-0 pb-8'>
      <MobileBar
        total={total}
        current={current}
        arrangement={arrangement}
        toggleArrangement={toggleArrangement}
        page={page}
      />
      <DesktopBar
        total={total}
        current={current}
        arrangement={arrangement}
        toggleArrangement={toggleArrangement}
        page={page}
      />
      <div
        className={clsx(
          'w-full flex flex-col flex-wrap gap-4 flex-1 p-4 items-center md:items-start lg:p-8 lg:mt-16',
          arrangement === 'grid' ? 'md:flex-row' : 'lg:justify-start',
        )}
      >
        {stars?.map((e, i) =>
          arrangement === 'grid' ? <GridRepo key={i} {...e} index={i} /> : <ListRepo key={i} {...e} index={i} />,
        )}
      </div>
    </div>
  )
}
