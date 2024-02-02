'use server'

import { PrismaClient, User } from '@prisma/client'
import Stars from '@/pages/stars'

interface UserNamePage {
  username: string
}

export default async function Home({ params }: { params: UserNamePage }) {
  const { username } = params
  const prisma = new PrismaClient()

  return (
    <div className='w-full h-full p-12'>
      <Stars />
    </div>
  )
}
