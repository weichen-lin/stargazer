'use client'

import Rename from './rename'
import EditDescription from './edit-description'
import Switcher from './switcher'
import { useCollectionContext } from '@/app/collections/hooks/useCollectionContext'
import { ArrowDown } from 'lucide-react'
import Link from 'next/link'

export default function Operator() {
  const { collection } = useCollectionContext()
  return (
    <div className='md:pl-0 flex flex-col gap-y-4'>
      <div className='flex items-center justify-start gap-x-6'>
        <Link href='/collections'>
          <ArrowDown className='w-12 h-12 rotate-90 hover:bg-slate-200 p-3 rounded-full mr-2 cursor-pointer' />
        </Link>
        <h1 className='text-2xl font-bold underline'>{collection.name}</h1>
        <Switcher />
      </div>
      <div className='text-slate-500 pl-2'>{collection.description}</div>
      <div className='flex gap-x-3 items-center'>
        <Rename />
        <EditDescription />
      </div>
    </div>
  )
}
