'use server'

import { cn } from '@/lib/utils'
import { GetUser } from '@/actions'
import { CollectionClient } from '@/client/collection'
import { redirect } from 'next/navigation'
import { Owner } from './component'
import { Badge } from '@/components/ui/badge'
import { ArrowDown } from 'lucide-react'
import Link from 'next/link'
import { Rename, EditDescription, Switcher } from './component/operator'
import { CollectionProvider } from '@/app/collections/hooks/useCollectionContext'

export default async function Page(req: {
  params: {
    id: string
  }
}) {
  const { id } = req.params
  const { email } = await GetUser()

  const client = new CollectionClient(email)

  try {
    const getSharedCollection = await client.getCollection(id)
    const { owner, collection, shared_from } = getSharedCollection

    return (
      <div
        className={cn(
          'flex flex-col items-start dark:bg-black',
          'w-full flex-1 3xl:py-8 gap-y-3',
          'lg:content-around p-6',
        )}
      >
        {shared_from && <Owner {...shared_from} />}
        <div className='md:pl-0 flex flex-col gap-y-4'>
          <CollectionProvider collection={collection}>
            <div className='flex items-center justify-start gap-x-6'>
              <Link href='/collections'>
                <ArrowDown className='w-12 h-12 rotate-90 hover:bg-slate-200 p-3 rounded-full mr-2 cursor-pointer' />
              </Link>
              <h1 className='text-2xl font-bold underline'>{collection.name}</h1>
              <h2 className='font-light'>{collection.is_public ? <Badge>Public</Badge> : <Badge>Private</Badge>}</h2>
            </div>
            <div className='text-slate-500 pl-2'>{collection.description}</div>
            <div className='flex gap-x-3 items-center'>
              <Rename />
              <EditDescription />
              <Switcher />
            </div>
          </CollectionProvider>
        </div>
      </div>
    )
  } catch (error) {
    redirect('/404')
  }
}
