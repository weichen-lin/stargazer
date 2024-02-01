'use server'

import { PrismaClient, User } from '@prisma/client'

interface UserNamePage {
  username: string
}

export default async function Home({ params }: { params: UserNamePage }) {
  const { username } = params
  const prisma = new PrismaClient()
  const users = await prisma.user.findMany()

  return (
    <div className=''>
      {users.map((e: User) => {
        return <div key={e.id}>{e.name}</div>
      })}
    </div>
  )
}
