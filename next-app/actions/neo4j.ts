'use server'

import { conn } from '@/actions/adapter'
import { Integer } from 'neo4j-driver'

const fetcher = async <T>(q: string, params: T) => {
  const session = conn.session()
  try {
    const res = await session.executeRead(tx => tx.run(q, params))
    return res.records.map(r => r.toObject())
  } catch (error) {
    console.error(error)
  } finally {
    await session.close()
  }
}

export const getUserProviderInfo = async (name: string): Promise<string> => {
  const q = `
  MATCH (u:User { name: $name })-[:HAS_ACCOUNT]->(a:Account)
  RETURN a{.*}
  `

  const data = await fetcher(q, { name })
  const target = Array.isArray(data) ? data[0] : data
  return target?.a?.providerAccountId ?? null
}

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

export interface IUserSetting {
  openAIKey: string
  githubToken: string
  limit: number
  cosine: number
}

export interface UserInfo extends IUserSetting {
  email: string
  name: string
}

export const getUserInfo = async (params: { email: string }): Promise<UserInfo | null> => {
  const q = `
  MATCH (n:User {email: $email})
  RETURN n;
  `

  const data = await fetcher(q, params)
  const target = Array.isArray(data) ? data[0] : data
  if (target) {
    return {
      email: target?.n.properties.email,
      name: target?.n.properties.name,
      limit: target?.n.properties?.limit ?? 5,
      openAIKey: target?.n.properties?.openAIKey || null,
      githubToken: target?.n.properties?.githubToken || null,
      cosine: target?.n.properties?.cosine ?? 0.8,
    }
  } else {
    return null
  }
}

interface UpdateInfoParams {
  email: string
  limit: number
  openAIKey: string
  githubToken: string
  cosine: number
}

export const updateInfo = async (params: UpdateInfoParams) => {
  const q = `
  MATCH (u:User { email: $email })
  SET u += $data
  RETURN u{.*}
  `
  const session = conn.session()

  try {
    await session.executeWrite(tx => tx.run(q, { data: params, email: params.email }))
    return 'OK'
  } catch (error) {
    console.error(error)
  } finally {
    await session.close()
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

export const updateStarRelation = async (email: string, repo_id: number, isVectorized: boolean) => {
  const q = `
  MATCH (u:User { email: $email })-[s:STARS]->(r:Repository { repo_id: $repo_id })
  SET s.is_vectorized = $isVectorized
  RETURN s{.*}
  `

  const session = conn.session()

  try {
    await session.executeWrite(tx => tx.run(q, { email, repo_id, isVectorized }))
  } catch (error) {
    console.error(error)
  } finally {
    await session.close()
  }
}

export default fetcher
