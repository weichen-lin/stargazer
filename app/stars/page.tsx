'use server'

import { redirect } from 'next/navigation'
import { getUserRepos } from '@/actions/neo4j'
import { int } from 'neo4j-driver'
import { z } from 'zod'
import Stars from '@/components/page/starsPage'
import { getServerSession } from 'next-auth'
import { options } from '@/app/api/auth/[...nextauth]/option'

const pageSchema = z.object({
  p: z.string(),
})

const parsePage = (a: { p: string }): number => {
  try {
    pageSchema.parse(a)
    return parseInt(a.p)
  } catch {
    return 1
  }
}

export default async function Home({ searchParams }: { searchParams: { p: string } }) {
  const session = await getServerSession(options)
  if (!session) {
    redirect('/')
  }
  const name = (session as any)?.user?.name ?? ''
  const page = parsePage(searchParams as any)
  const { total, stars } = await getUserRepos({ username: name, page: int(page), limit: int(20) })

  return (
    <div className='w-full h-full flex flex-col lg:flex-row'>
      <Stars stars={stars} total={total} />
    </div>
  )
}
