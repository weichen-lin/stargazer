'use server'

import { conn } from '@/actions/adapter'
import logger from '@/utils/logger'

const fetcher = async <T>(q: string, params: T) => {
  const now = Date.now()
  const session = conn.session()
  try {
    const res = await session.executeRead(tx => tx.run(q, params))
    return res.records.map(r => r.toObject())
  } catch (error) {
    logger.error({ error, q, params, duration: Date.now() - now })
  } finally {
    await session.close()
    logger.info({ q, params, duration: Date.now() - now })
  }
}

const writer = async <T>(q: string, params: T) => {
  const now = Date.now()
  const session = conn.session()
  try {
    const res = await session.executeWrite(tx => tx.run(q, params))
    return res.records.map(r => r.toObject())
  } catch (error) {
    logger.error({ error, q, params, duration: Date.now() - now })
  } finally {
    await session.close()
    logger.info({ q, params, duration: Date.now() - now })
  }
}

export { fetcher, writer }
