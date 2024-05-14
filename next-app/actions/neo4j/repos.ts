'use server'

import fetcher from './fetcher'
import { z } from 'zod'
import { int } from 'neo4j-driver'

export type Repository = {
  repo_id: number
  owner_url: string
  owner_name: string
  open_issues_count: number
  full_name: string
  avatar_url: string
  html_url: string
  description: string
  homepage: string
  stargazers_count: number
  language: string
  last_updated_at: string
}

const UserReposParamsSchema = z.object({
  email: z.string(),
  languages: z.array(z.string()),
  page: z.number(),
  limit: z.number(),
})

type UserReposParams = z.infer<typeof UserReposParamsSchema>

export const getUserRepos = async (params: UserReposParams): Promise<{ total: number; repos: Repository[] }> => {
  const q = `
  MATCH (u:User {email: $email})-[s:STARS]-(r:Repository)
  WHERE r.language IN $languages
  WITH u, COUNT(r) as total
  MATCH (u)-[s:STARS]-(r:Repository)
  WHERE r.language IN $languages
  WITH total, s, r
  ORDER BY s.created_at DESC
  SKIP $limit * ($page - 1)
  LIMIT $limit
  RETURN total, collect(r) as data
  `

  try {
    const p = UserReposParamsSchema.parse(params)

    const data = await fetcher(q, {
      email: p.email,
      languages: p.languages,
      page: int(p.page),
      limit: int(p.limit),
    })

    const target = Array.isArray(data) ? data[0] : data
    const total = target?.total?.low ?? 0

    const repos =
      target?.data?.map((e: any) => {
        const {
          repo_id,
          owner_url,
          owner_name,
          open_issues_count,
          full_name,
          avatar_url,
          html_url,
          description,
          homepage,
          stargazers_count,
          language,
          last_updated_at,
        } = e.properties

        return {
          repo_id: repo_id.low,
          owner_url,
          owner_name,
          open_issues_count: open_issues_count.low,
          full_name,
          avatar_url,
          html_url,
          description,
          homepage,
          stargazers_count: stargazers_count.low,
          language,
          last_updated_at: last_updated_at.toString(),
        }
      }) ?? []

    return {
      total,
      repos,
    }
  } catch (error) {
    console.error(error)
    return {
      total: 0,
      repos: [],
    }
  }
}

export const getUserStarsRelation = async (email: string) => {
  const q = `
  MATCH (u:User { email: $email })-[s:STARS]->(r:Repository)
  RETURN s{.*}
  `

  try {
    const data = await fetcher(q, { email })

    const isVectorized = data?.filter(e => e?.s?.is_vectorized === true)

    return {
      total: data?.length ?? 0,
      vectorized: isVectorized?.length ?? 0,
    }
  } catch (error) {
    console.error(error)
    return {
      total: 0,
      vectorized: 0,
    }
  }
}

export const getUserStarsRelationRepos = async (
  email: string,
): Promise<{ isVectorized: boolean; repo_id: number }[]> => {
  const q = `
  MATCH (u:User { email: $email })-[s:STARS]->(r:Repository)
  RETURN s{.*}, r{.repo_id}
  `

  try {
    const data = await fetcher(q, { email })
    const info = data?.map(e => ({ isVectorized: e?.s?.is_vectorized === true, repo_id: e?.r?.repo_id?.low })) ?? []
    return info.sort((a, b) => {
      if (a.isVectorized === b.isVectorized) {
        return 0
      }
      return a.isVectorized ? -1 : 1
    })
  } catch (error) {
    console.error(error)
    return []
  }
}

export interface ILanguageDistribution {
  language: string
  count: number
}

export const getLanguageDistribution = async (name: string): Promise<ILanguageDistribution[]> => {
  const q = `
  MATCH (u:User { name: $name })-[s:STARS]->(r:Repository)
  WITH r.language as language, COUNT(r) as count
  RETURN language, count
  ORDER BY count DESC
  `

  try {
    const data = await fetcher(q, { name })
    if (data && data?.length > 0) {
      return data.map((e: any) => ({ language: e?.language, count: e?.count?.low }))
    } else {
      return []
    }
  } catch (error) {
    console.error(error)
    return []
  }
}

export const getRepositoriesCount = async (name: string): Promise<number> => {
  const q = `
  MATCH (u: User {name: $name})-[:STARS {is_delete: false}]-(r:Repository)
  RETURN COUNT(r) as count;
  `

  try {
    const data = await fetcher(q, { name })
    if (data && data?.length > 0) {
      return data[0]?.count?.low ?? 0
    }
    return 0
  } catch (error) {
    console.error(error)
    return 0
  }
}

export interface IRepoAtDashboard {
  repo_id: number
  avatar_url: string
  full_name: string
  html_url: string
  open_issues_count: number
  created_at: string
  last_updated_at: string
}

export type ISearchKey = 'last_synced_at' | 'created_at' | 'open_issues_count' | 'last_updated_at'

const SearchQuery: { [key in ISearchKey]: string } = {
  last_synced_at: 's.last_synced_at',
  created_at: 's.created_at',
  open_issues_count: 'r.open_issues_count',
  last_updated_at: 'r.last_updated_at',
}

export const getReposByKey = async (email: string, key: ISearchKey): Promise<IRepoAtDashboard[]> => {
  const q = `
  MATCH (u:User {email: $email})-[s:STARS {is_delete: false}]-(r:Repository)
  RETURN
  r.repo_id as repo_id,
  r.full_name as full_name,
  r.avatar_url as avatar_url,
  r.html_url as html_url,
  r.open_issues_count as open_issues_count,
  s.created_at as created_at,
  r.last_updated_at as last_updated_at,
  s.last_synced_at as last_synced_at
  ORDER BY ${SearchQuery[key]} DESC
  LIMIT 5;
  `

  try {
    const data = await fetcher(q, { email })
    if (data && data?.length > 0) {
      return data.map((e: any) => ({
        repo_id: e?.repo_id?.low,
        full_name: e?.full_name,
        open_issues_count: e?.open_issues_count?.low,
        avatar_url: e?.avatar_url,
        created_at: e?.created_at?.toString(),
        html_url: e?.html_url,
        last_updated_at: e?.last_updated_at,
      }))
    }
    return []
  } catch (error) {
    console.error(error)
    return []
  }
}
