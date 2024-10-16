import TagClient from '@/client/tag'
import type { NextRequest } from 'next/server'

export const dynamic = 'force-dynamic'
export async function GET(req: NextRequest, { params }: { params: { repo_id: string } }) {
  const { repo_id } = params

  const client = new TagClient()

  const { data, status_code } = await client.getTags(repo_id)
  return Response.json(data, { status: status_code })
}

// export async function DELETE(request: NextRequest, { params }: { params: { repo_id: string } }) {
//   const { email } = await GetUser()
//   const { repo_id } = params

//   try {
//     const detail = await deleteRepo(email, repo_id)

//     return new Response(JSON.stringify(detail), {
//       headers: {
//         'Content-Type': 'application/json',
//       },
//     })
//   } catch (error) {
//     return new Response('error', {
//       status: 400,
//     })
//   }
// }
