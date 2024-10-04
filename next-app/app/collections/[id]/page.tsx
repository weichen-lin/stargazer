'use server'

import { cn } from '@/lib/utils'
import { GetUser } from '@/actions'
import { CollectionClient } from '@/client/collection'
import { redirect } from 'next/navigation'
import { Owner } from './component'

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
    // 9024c830-81aa-4018-b6ad-c8984740cc37

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
        <div className='md:pl-3'></div>
        123
      </div>
    )
  } catch (error) {
    redirect('/404')
  }
}
