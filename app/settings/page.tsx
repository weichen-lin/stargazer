'use server'

import { z } from 'zod'

const pageSchema = z.object({
  p: z.string(),
})

const parsePage = (a: { p: string }): number => {
  try {
    pageSchema.parse(a)
    return parseInt(a.p)
  } catch {
    return 1
  }
}
export default async function Home() {
  return <div className='w-full h-full flex flex-col lg:flex-row'>123 </div>
}
