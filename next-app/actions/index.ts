'use server'

import { getServerSession } from 'next-auth'
import { options } from '@/app/api/auth/[...nextauth]/option'
import { redirect } from 'next/navigation'
import { z } from 'zod'

const pageSchema = z.object({
  name: z.string().min(1),
  email: z.string().email(),
  image: z.string(),
})

export async function GetUser() {
  const session = await getServerSession(options)
  if (!session) {
    redirect('/')
  }

  try {
    const user = pageSchema.parse(session.user)
    return user
  } catch {
    redirect('/')
  }
}
