'use server'

import { PrismaClient, User } from '@prisma/client'
import Stars from '@/pages/stars'
import { redirect } from 'next/navigation'

interface UserNamePage {
  username: string
}

export default async function Home({ params }: { params: UserNamePage }) {
  const { username } = params
  const prisma = new PrismaClient()

  const user: User | null = await prisma.user.findFirst({
    where: {
      name: username,
    },
  })

  // if (!user) {
  //   redirect('/404')
  // }

  return (
    <div className='w-full h-full p-12'>
      <Stars />
    </div>
  )
}
