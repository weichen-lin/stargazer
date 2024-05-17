'use server'

import { fetcher, writer } from './fetcher'
import { z } from 'zod'

export const getTags = async (email: string) => {
  const q = `
  MATCH (u:User)-[:HAS_TAG]->(t:Tag)
  return t.id as id, t.name as name
  `

  const data = await fetcher(q, { email })
  return data
}

export const createTagSchema = z.object({
  email: z.string(),
  name: z.string(),
  repo_id: z.number(),
})

export type CreateTag = z.infer<typeof createTagSchema>

export const createTag = async (props: CreateTag) => {
  const params = createTagSchema.parse(props)

  const q = `
    MATCH (u:User {email: $email})
    MERGE (t:Tag {name: $name})
    ON CREATE SET t.id = randomUUID()
    MERGE (u)-[:HAS_TAG]->(t)
    WITH t
    MATCH (r:Repository {repo_id: $repo_id})
    WITH t, r
    MERGE (r)-[:TAGGED_BY]->(t)
    RETURN t.id as id, t.name as name
    `

  const data = await writer(q, params)
  return data
}
