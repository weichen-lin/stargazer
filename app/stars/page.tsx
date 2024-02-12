'use server'

import { redirect } from 'next/navigation'
import { Neo4jfetcher, getUserRepos } from '@/actions'
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

  const data = await Neo4jfetcher(getUserRepos, { username: name, page: int(page), limit: int(20) })

  const target = Array.isArray(data) ? data[0] : data
  const total = target?.total?.low ?? 0

  const stars = target?.limitedRepositories
    ? target?.limitedRepositories.map((e: any) => {
        return {
          id: e.properties.id,
          full_name: e.properties.full_name,
          owner: {
            avatar_url: e.properties.avatar_url,
          },
          html_url: e.properties.html_url,
          description: e.properties.description,
          homepage: e.properties.homepage ?? '',
          stargazers_count: e.properties.stargazers_count.low,
          language: e.properties.language,
        }
      })
    : []

  return (
    <div className='w-full h-full flex flex-col lg:flex-row'>
      <Stars stars={stars} total={total} />
    </div>
  )
}
