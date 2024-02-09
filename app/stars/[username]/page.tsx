'use server'

import Stars from '@/pages/stars'
import { redirect } from 'next/navigation'
import { Neo4jfetcher, getUserRepos } from '@/api'
import { int } from 'neo4j-driver'
import { z } from 'zod'

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

interface UserNamePage {
  username: string
}

export default async function Home({ params, searchParams }: { params: UserNamePage; searchParams: { p: string } }) {
  const { username } = params
  const page = parsePage(searchParams as any)

  const data = await Neo4jfetcher(getUserRepos, { username: username, page: int(page), limit: int(20) })

  if (!data) {
    redirect('/404')
  }

  const target = Array.isArray(data) ? data[0] : data
  const total = target.total.low ?? 0

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
      <Stars stars={stars} total={total} current={page} page={`/stars/${username}?p=`} />
    </div>
  )
}
