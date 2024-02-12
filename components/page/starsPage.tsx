'use client'

import { useArrangement } from '@/hooks/stars'
import { GridRepo, ListRepo } from '@/components/repo'
import clsx from 'clsx'
import { MobileBar, DesktopBar } from '@/components/sidebar/star'
import { useSearchParams } from 'next/navigation'

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

export default function Stars({ stars, total }: { stars: Star[]; total: number }) {
  const { arrangement, toggleArrangement } = useArrangement()

  return (
    <div className='flex flex-col gap-y-12 w-full mt-36 lg:mt-0 pb-8'>
      <MobileBar total={total} arrangement={arrangement} toggleArrangement={toggleArrangement} />
      <DesktopBar total={total} arrangement={arrangement} toggleArrangement={toggleArrangement} />
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
