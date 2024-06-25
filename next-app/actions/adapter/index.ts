import { isInt, integer, Driver } from 'neo4j-driver'
import type { Adapter } from '@auth/core/adapters'
import { driver, auth } from 'neo4j-driver'

export const conn = driver(
  process.env.NEO4J_URL ?? 'neo4j://neo4j:7687',
  auth.basic('neo4j', process.env.NEO4J_PASSWORD ?? 'password'),
)

export function Neo4jAdapter(driver: Driver): Adapter {
  const { read, write } = client(driver)

  return {
    async createUser(data) {
      const id = crypto.randomUUID()
      const user = { ...data, id}
      await write(`CREATE (u:User $data)`, { data: { ...user } })
      return user
    },

    async getUser(id) {
      return await read(`MATCH (u:User { id: $id }) RETURN u{.*}`, {
        id,
      })
    },

    async getUserByEmail(email) {
      return await read(`MATCH (u:User { email: $email }) RETURN u{.*}`, {
        email,
      })
    },

    async getUserByAccount(provider_providerAccountId) {
      return await read(
        `MATCH (u:User)-[:HAS_ACCOUNT]->(a:Account {
           providerAccountId: $providerAccountId
         })
         RETURN u{.*}`,
        provider_providerAccountId,
      )
    },

    async updateUser(data) {
      return (
        await write(
          `MATCH (u:User { id: $id })
           SET u += $data
           RETURN u{.*}`,
          data,
        )
      ).u
    },

    async deleteUser(id) {
      return await write(
        `MATCH (u:User { id: $id })
         WITH u, u{.*} AS properties
         DETACH DELETE u
         RETURN properties`,
        { id },
      )
    },

    async linkAccount(data) {
      const { userId, ...a } = data
      await write(
        `MATCH (u:User { id: $userId })
         MERGE (a:Account {
           providerAccountId: $a.providerAccountId,
           provider: $a.provider
         }) 
         SET a += $a
         MERGE (u)-[:HAS_ACCOUNT]->(a)`,
        { userId, a },
      )
      return data
    },

    async unlinkAccount(provider_providerAccountId) {
      return await write(
        `MATCH (u:User)-[:HAS_ACCOUNT]->(a:Account {
           providerAccountId: $providerAccountId,
           provider: $provider
         })
         WITH u, a, properties(a) AS properties
         DETACH DELETE a
         RETURN properties { .*, userId: u.id }`,
        provider_providerAccountId,
      )
    },

    async createSession(data) {
      const { userId, ...s } = format.to(data)
      await write(
        `MATCH (u:User { id: $userId })
         CREATE (s:Session)
         SET s = $s
         CREATE (u)-[:HAS_SESSION]->(s)`,
        { userId, s },
      )
      return data
    },

    async getSessionAndUser(sessionToken) {
      const result = await write(
        `OPTIONAL MATCH (u:User)-[:HAS_SESSION]->(s:Session { sessionToken: $sessionToken })
         WHERE s.expires <= datetime($now)
         DETACH DELETE s
         WITH count(s) AS c
         MATCH (u:User)-[:HAS_SESSION]->(s:Session { sessionToken: $sessionToken })
         RETURN s { .*, userId: u.id } AS session, u{.*} AS user`,
        { sessionToken, now: new Date().toISOString() },
      )

      if (!result?.session || !result?.user) return null

      return {
        session: format.from<any>(result.session),
        user: format.from<any>(result.user),
      }
    },

    async updateSession(data) {
      return await write(
        `MATCH (u:User)-[:HAS_SESSION]->(s:Session { sessionToken: $sessionToken })
         SET s += $data
         RETURN s { .*, userId: u.id }`,
        data,
      )
    },

    async deleteSession(sessionToken) {
      return await write(
        `MATCH (u:User)-[:HAS_SESSION]->(s:Session { sessionToken: $sessionToken })
         WITH u, s, properties(s) AS properties
         DETACH DELETE s
         RETURN properties { .*, userId: u.id }`,
        { sessionToken },
      )
    },

    async createVerificationToken(data) {
      await write(
        `MERGE (v:VerificationToken {
           identifier: $identifier,
           token: $token
         })
         SET v = $data`,
        data,
      )
      return data
    },

    async useVerificationToken(data) {
      const result = await write(
        `MATCH (v:VerificationToken {
           identifier: $identifier,
           token: $token
         })
         WITH v, properties(v) as properties
         DETACH DELETE v
         RETURN properties`,
        data,
      )
      return format.from<any>(result?.properties)
    },
  }
}

// https://github.com/honeinc/is-iso-date/blob/master/index.js
const isoDateRE =
  /(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+([+-][0-2]\d:[0-5]\d|Z))|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d([+-][0-2]\d:[0-5]\d|Z))|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d([+-][0-2]\d:[0-5]\d|Z))/

function isDate(value: any) {
  return value && isoDateRE.test(value) && !isNaN(Date.parse(value))
}

export const format = {
  /** Takes a plain old JavaScript object and turns it into a Neo4j compatible object */
  to(object: Record<string, any>) {
    const newObject: Record<string, unknown> = {}
    for (const key in object) {
      const value = object[key]
      if (value instanceof Date) newObject[key] = value.toISOString()
      else newObject[key] = value
    }
    return newObject
  },
  /** Takes a Neo4j object and returns a plain old JavaScript object */
  from<T = Record<string, unknown>>(object?: Record<string, any>): T | null {
    const newObject: Record<string, unknown> = {}
    if (!object) return null
    for (const key in object) {
      const value = object[key]
      if (isDate(value)) {
        newObject[key] = new Date(value)
      } else if (isInt(value)) {
        if (integer.inSafeRange(value)) newObject[key] = value.toNumber()
        else newObject[key] = value.toString()
      } else {
        newObject[key] = value
      }
    }

    return newObject as T
  },
}

export function client(conn: Driver) {
  return {
    /** Reads values from the database */
    async read<T>(statement: string, values?: any): Promise<T | null> {
      const session = conn.session()
      try {
        const result = await session.executeRead(tx => tx.run(statement, values))
        return format.from<T>(result?.records[0]?.get(0)) ?? null
      } catch (error) {
        console.log('error at read transaction', error)
        return null
      } finally {
        session.close()
      }
    },
    /**
     * Reads/writes values from/to the database.
     * Properties are available under `$data`
     */
    async write<T extends Record<string, any>>(statement: string, values: T): Promise<any> {
      const session = conn.session()
      try {
        const result = await session.executeWrite(tx => tx.run(statement, values))
        return format.from<T>(result?.records[0]?.get(0))
      } catch (error) {
        console.log('error at write transaction', statement, error)
        return null
      } finally {
        session.close()
      }
    },
  }
}
