'use server'

import { fetcher, writer } from './fetcher'
import { z } from 'zod'

const IdNameSchema = z.object({
  id: z.string(),
  name: z.string(),
})

type IdName = z.infer<typeof IdNameSchema>

export const getTags = async (email: string): Promise<IdName[]> => {
  const q = `
  MATCH (u:User)-[:HAS_TAG]->(t:Tag)
  return t.id as id, t.name as name
  `

  try {
    const data = await fetcher(q, { email })
    const tags = data?.map((tag: any) => IdNameSchema.parse(tag)) ?? []

    return tags
  } catch (error) {
    console.error(error)
    return []
  }
}

const createTagSchema = z.object({
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
    MATCH (u)-[s:STARS]-(r:Repository {repo_id: $repo_id})
    SET s.last_modified_at = datetime()
    WITH t, r
    MERGE (r)-[:TAGGED_BY]->(t)
    RETURN t.id as id, t.name as name
    `

  const data = await writer(q, params)
  const tag = Array.isArray(data) && IdNameSchema.parse(data[0])

  return tag
}

export const getTagsByRepo = async (email: string, repo_id: number): Promise<IdName[]> => {
  const q = `
  MATCH (u:User {email: $email})-[:HAS_TAG]->(t:Tag)<-[:TAGGED_BY]-(r:Repository {repo_id: $repo_id})
  RETURN t.id as id, t.name as name
  `

  try {
    const data = await fetcher(q, { email, repo_id })
    const tags = data?.map((tag: any) => IdNameSchema.parse(tag)) ?? []

    return tags
  } catch (error) {
    console.error(error)
    return []
  }
}

export const deleteTagByRepo = async (email: string, repo_id: number, tag_id: string) => {
  const q = `
  MATCH (u:User {email: $email})-[s:STARS]->(r:Repository {repo_id: $repo_id})
  SET s.last_modified_at = datetime()
  WITH u, s, r
  MATCH (u)-[:HAS_TAG]->(t:Tag {id: $tag_id})<-[tb:TAGGED_BY]-(r)
  DELETE tb
  `

  await writer(q, { email, repo_id, tag_id })
}
