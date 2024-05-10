'use server'

import fetcher from './fetcher'
import { Integer, DateTime } from 'neo4j-driver'

interface UserReposParams {
  username: string
  page: Integer
  limit: Integer
}

export type Repository = {
  id: string
  full_name: string
  owner: {
    avatar_url: string
  }
  html_url: string
  description: string
  homepage: string
  stargazers_count: number
  language: string
}

export const getUserRepos = async (params: UserReposParams): Promise<{ total: number; stars: Repository[] }> => {
  const q = `
  MATCH (n:User {name: $username})-[:STARS]->(r:Repository)
  WITH count(r) as total
  MATCH (n:User {name: $username})-[:STARS]->(r:Repository)
  WITH total, r
  ORDER BY r.created_at DESC
  SKIP $limit * ($page - 1)
  LIMIT $limit
  WITH total, collect(r) as l
  RETURN total, l;
  `

  const data = await fetcher(q, params)

  const target = Array.isArray(data) ? data[0] : data
  const total = target?.total?.low ?? 0

  const stars = target?.l
    ? target?.l.map((e: any) => {
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

  return {
    total,
    stars,
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

export const getReposByKey = async (name: string, key: string): Promise<IRepoAtDashboard[]> => {
  const q = `
  MATCH (u: User {name: $name})-[:STARS {is_delete: false}]-(r:Repository)
  RETURN 
  r.repo_id as repo_id, 
  r.full_name as full_name, 
  r.avatar_url as avatar_url, 
  r.html_url as html_url,
  r.open_issues_count as open_issues_count,
  r.created_at as created_at, 
  r.last_updated_at as last_updated_at
  ORDER BY r.${key} DESC
  LIMIT 5;
  `

  try {
    const data = await fetcher(q, { name })
    if (data && data?.length > 0) {
      return data.map((e: any) => ({
        repo_id: e?.repo_id?.low,
        full_name: e?.full_name,
        open_issues_count: e?.open_issues_count?.low,
        avatar_url: e?.avatar_url,
        created_at: e?.created_at.toString(),
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