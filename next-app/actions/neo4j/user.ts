import { conn } from '@/actions/adapter'
import fetcher from './fetcher'

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

export const getUserProviderInfo = async (name: string): Promise<string> => {
  const q = `
  MATCH (u:User { name: $name })-[:HAS_ACCOUNT]->(a:Account)
  RETURN a{.*}
  `

  const data = await fetcher(q, { name })
  const target = Array.isArray(data) ? data[0] : data
  return target?.a?.providerAccountId ?? null
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
