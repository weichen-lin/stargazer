'use server'

import { randomUUID } from 'crypto'
import axios from 'axios'
import { getServerSession } from 'next-auth'
import { options } from '@/app/api/auth/[...nextauth]/option'
import { redirect } from 'next/navigation'
import { z } from 'zod'

const userSchema = z.object({
  name: z.string().min(1),
  email: z.string().email(),
  image: z.string(),
})

const fetcher = async () => {
  let isHealth: boolean = false

  const session = await getServerSession(options)

  return {}
}
