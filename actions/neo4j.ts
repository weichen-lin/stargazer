'use server'

import { driver, auth } from 'neo4j-driver'

const conn = driver(process.env.NEO4J_URL as string, auth.basic('neo4j', process.env.NEO4J_PASSWORD as string))

const fetcher = async <T>(q: string, params: T) => {
  const session = conn.session()
  try {
    const res = await session.run(q, params)
    return res.records.map(r => r.toObject())
  } catch (error) {
    console.error(error)
  } finally {
    await session.close()
  }
}

export default fetcher
