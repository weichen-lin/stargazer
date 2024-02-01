'use server'

import { PrismaClient, User } from '@prisma/client'

export default async function Home() {
  const users = await prisma.user.findMany()
  return (
    <div className='bg-red-300'>
      {users.map((e: User) => {
        return <div key={e.id}>{e.name}</div>
      })}
    </div>
  )
}
