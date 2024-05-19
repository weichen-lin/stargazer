import { getRepoDetail, deleteRepo, getTagsByRepo } from '@/actions/neo4j'
import { NextRequest } from 'next/server'
import { GetUser } from '@/actions'

export const dynamic = 'force-dynamic'
export async function GET(request: NextRequest, { params }: { params: { repo_id: string } }) {
  const { email } = await GetUser()
  const { repo_id } = params

  try {
    const detail = await getRepoDetail(email, repo_id)

    return new Response(JSON.stringify(detail), {
      headers: {
        'Content-Type': 'application/json',
      },
    })
  } catch (error) {
    return new Response('error', {
      status: 400,
    })
  }
}

export async function DELETE(request: NextRequest, { params }: { params: { repo_id: string } }) {
  const { email } = await GetUser()
  const { repo_id } = params

  try {
    const detail = await deleteRepo(email, repo_id)

    return new Response(JSON.stringify(detail), {
      headers: {
        'Content-Type': 'application/json',
      },
    })
  } catch (error) {
    return new Response('error', {
      status: 400,
    })
  }
}
