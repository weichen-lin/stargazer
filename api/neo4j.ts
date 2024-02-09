import { driver, auth } from 'neo4j-driver'

const conn = driver(
  'neo4j+s://ac6f9976.databases.neo4j.io',
  auth.basic('neo4j', 'O2Mara0PocLB8L1mR7K_8btVeUMnnn7p5L3fOgn2bFA'),
)

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
