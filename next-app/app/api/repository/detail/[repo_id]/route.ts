import { RepositoryClient } from '@/client/repository'
import type { NextRequest } from 'next/server'

export const dynamic = 'force-dynamic'
export async function GET(req: NextRequest, { params }: { params: { repo_id: string } }) {
  const { repo_id } = params

  const client = new RepositoryClient()
  const { status_code, data } = await client.getRepoDetail(repo_id)

  return Response.json(data, { status: status_code })
}

export async function DELETE({ params }: { params: { repo_id: string } }) {
  const { repo_id } = params

  const client = new RepositoryClient()

  const { status_code, data } = await client.deleteRepo(repo_id)
  return Response.json(data, { status: status_code })
}
