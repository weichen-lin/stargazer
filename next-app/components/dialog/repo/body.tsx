'use client'

import { Textarea } from '@/components/ui/textarea'
import { IRepository } from '@/client/repository'
import Tags from './tags'

export default function Body(props: IRepository) {
  const { repo_id, owner_name, repo_name } = props
  return (
    <div className='flex-1 overflow-y-auto p-3'>
      <Tags repo_id={repo_id} />
      <div className='flex flex-col mb-4 gap-y-2'>
        <label className='text-sm text-slate-500 dark:text-slate-400'>Summary</label>
        {/* <Textarea value={gpt_summary} className='w-full' rows={5} onChange={() => {}} /> */}
      </div>
      <img
        src={`https://api.star-history.com/svg?repos=${owner_name}/${repo_name}&type=Date&theme=light`}
        className='max-h-[350px] mx-auto'
      ></img>
    </div>
  )
}
