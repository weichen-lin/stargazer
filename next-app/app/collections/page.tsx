'use client'

import { cn } from '@/lib/utils'
import { Operator, Collections } from './components'

export default function CollectionPage() {
  return (
    <div
      className={cn(
        'flex flex-col items-start dark:bg-black',
        'w-full flex-1 3xl:py-8 gap-y-3',
        'lg:content-around p-6',
      )}
    >
      <div className='md:pl-3'>
        <Operator />
      </div>
      <Collections />
    </div>
  )
}
