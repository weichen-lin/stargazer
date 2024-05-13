'use client'

import { useArrangement } from '@/hooks/dashboard'
import { GridRepo, ListRepo } from '@/components/repo'
import clsx from 'clsx'
import { MobileBar, DesktopBar } from '@/components/sidebar/dashboard'
import type { Repository } from '@/actions/neo4j'
import Image from 'next/image'
import { SymbolIcon } from '@radix-ui/react-icons'
import { Suspense } from 'react'
import { LanguagePie } from '@/components/charts'

export default function Stars({ stars, total }: { stars: Repository[]; total: number }) {
  const { arrangement, toggleArrangement } = useArrangement()

  return (
    <div className='flex flex-col gap-y-12 w-full mt-36 lg:mt-0 pb-8 h-full'>
      <MobileBar total={total} arrangement={arrangement} toggleArrangement={toggleArrangement} />
      <DesktopBar total={total} arrangement={arrangement} toggleArrangement={toggleArrangement} />
    </div>
  )
}
