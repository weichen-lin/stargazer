'use server'

import { cn } from '@/lib/utils'
import { CollectionClient } from '@/client/collection'
import { redirect } from 'next/navigation'
import { Owner } from './component'
import Operator from './component/operator'
import { CollectionProvider } from '../hooks/useCollectionContext'
import Results from './component/collect-repos/results'

export default async function Page(req: {
  params: {
    id: string
  }
}) {
  const { id } = req.params

  const client = new CollectionClient()

  try {
    const { data: getSharedCollection } = await client.getCollection(id)
    const { collection, shared_from } = getSharedCollection

    return (
      <div
        className={cn(
          'flex flex-col items-start dark:bg-black',
          'w-full flex-1 3xl:py-8 gap-y-3',
          'lg:content-around p-6',
        )}
      >
        {shared_from && <Owner {...shared_from} />}
        <CollectionProvider initCollection={collection}>
          <Operator />
          <div className='w-full flex flex-col items-center justify-center flex-1 overflow-y-auto mb-8'>
            <Results />
          </div>
        </CollectionProvider>
      </div>
    )
  } catch (error) {
    redirect('/404')
  }
}
