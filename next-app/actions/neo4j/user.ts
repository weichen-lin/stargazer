'use server'

import { conn } from '@/actions/adapter'
import fetcher from './fetcher'

export const getUserInfo = async (email: string): Promise<IUserConfig> => {
  const q = `
  MATCH (u:User {
    email: $email
  })
  MERGE (u)-[h:HAS_CONFIG]->(c:Config)
  SET c = {
    limit : COALESCE(c.limit , 5),
    cosine: COALESCE(c.cosine, 0.8),
    openai_key: COALESCE(c.openai_key, ''),
    github_token: COALESCE(c.github_token, '')
  }
  RETURN c{.*};
  `
  const session = conn.session()

  try {
    const result = await session.executeWrite(tx => tx.run(q, { email }))
    const target = Array.isArray(result.records) ? result.records[0] : result.records

    const data = target.toObject()
    return {
      limit: data?.c?.limit?.low,
      openAIKey: data?.c?.openai_key,
      githubToken: data?.c?.github_token,
      cosine: data?.c?.cosine,
    }
  } catch (error) {
    console.log('error', error)
    return {
      limit: 5,
      openAIKey: '',
      githubToken: '',
      cosine: 0.8,
    }
  } finally {
    session.close()
  }
}

export interface IUserConfig {
  openAIKey: string
  githubToken: string
  limit: number
  cosine: number
}

export const updateInfo = async (params: IUserConfig & { email: string }): Promise<boolean> => {
  const q = `
  MATCH (n:User {email: $email})
  MERGE (n)-[h:HAS_CONFIG]->(c:Config)
  SET c = {
    limit: COALESCE($limit, 5),
    cosine: COALESCE($cosine, 0.8),
    openai_key: COALESCE($openAIKey, 0.8),
    github_token: COALESCE($githubToken, 0.8)
  }
  RETURN c;
  `
  const session = conn.session()

  try {
    await session.executeWrite(tx => tx.run(q, params))
    return true
  } catch (error) {
    console.error(error)
    return false
  } finally {
    await session.close()
  }
}

interface IUserCrontab {
  status: boolean
  lastTriggerTime: string
  time: number
}

export const getCrontabInfo = async (name: string): Promise<undefined> => {
  const q = `
  MATCH (u:User { name: $name })-[:HAS_CRONTAB]->(c:Crontab)
  RETURN c{.*}, 
  `
  const session = conn.session()

  try {
    const data = await fetcher(q, { name })
    return undefined
  } catch (error) {
    console.error({ error })
  } finally {
    await session.close()
  }
}
