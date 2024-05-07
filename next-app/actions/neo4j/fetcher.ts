'use server'

import { conn } from '@/actions/adapter'

const fetcher = async <T>(q: string, params: T) => {
  const session = conn.session()
  try {
    const res = await session.executeRead(tx => tx.run(q, params))
    return res.records.map(r => r.toObject())
  } catch (error) {
    console.error('error', error)
  } finally {
    await session.close()
  }
}

export default fetcher
