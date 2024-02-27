'use client'

import { useArrangement } from '@/hooks/stars'
import { GridRepo, ListRepo } from '@/components/repo'
import clsx from 'clsx'
import { MobileBar, DesktopBar } from '@/components/sidebar/star'
import type { Repository } from '@/actions/neo4j'
import Image from 'next/image'
import { SymbolIcon } from '@radix-ui/react-icons'

export default function Stars({ stars, total }: { stars: Repository[]; total: number }) {
  const { arrangement, toggleArrangement } = useArrangement()

  return (
    <div className='flex flex-col gap-y-12 w-full mt-36 lg:mt-0 pb-8 h-full'>
      <MobileBar total={total} arrangement={arrangement} toggleArrangement={toggleArrangement} />
      <DesktopBar total={total} arrangement={arrangement} toggleArrangement={toggleArrangement} />
      <div
        className={clsx(
          'w-full flex flex-col flex-wrap gap-4 flex-1 p-4 items-center md:items-start lg:p-8 lg:mt-16',
          arrangement === 'grid' ? 'md:flex-row' : 'lg:justify-start',
        )}
      >
        {total > 0 ? (
          stars?.map((e, i) =>
            arrangement === 'grid' ? <GridRepo key={i} {...e} index={i} /> : <ListRepo key={i} {...e} index={i} />,
          )
        ) : (
          <div className='w-full flex flex-col items-center justify-center gap-y-4 max-w-[400px] h-full md:max-w-full'>
            <Image src='/empty.png' width={250} height={250} alt='empty' />
            <div className='text-gray-500 text-center'>
              It seems that Stargazer doesn't have information about your GitHub stars.
            </div>
            <div className='text-gray-500 text-center'>
              Please click on the symbol in the upper right corner
              <SymbolIcon className='w-5 h-5 inline-block mx-2' />
              to begin syncing the data.
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
