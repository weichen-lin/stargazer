'use server'

import OpenAI from 'openai'
import { PrismaClient } from '@prisma/client'
import pgvector from 'pgvector'

const openai = new OpenAI({
  apiKey: process.env.OPENAI_API_KEY,
})

export const SemanticSearch = async (query: string) => {
  const prisma = new PrismaClient()

  const resp = await openai.embeddings.create({
    model: 'text-embedding-ada-002',
    input: [query],
    encoding_format: 'float',
  })

  const vector = pgvector.toSql(resp.data[0].embedding)

  const items = await prisma.$queryRaw`
      SELECT id, full_name, avatar_url, html_url
      FROM repo_embedding_info 
      WHERE description_embedding <-> ${vector}::vector > 0.75 
      ORDER BY description_embedding <-> ${vector}::vector LIMIT 5`

  return items
}
