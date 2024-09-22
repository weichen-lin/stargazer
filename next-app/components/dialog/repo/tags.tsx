'use client'

import { useState } from 'react'
import MultipleSelector, { Option } from '@/components/ui/multiple-selector'
import { useFetch } from '@/hooks/util'
import { ITag } from '@/client/tag'

const Tags = (props: { repo_id: number }) => {
  const { repo_id } = props
  const [currentTags, setCurrentTags] = useState<{ value: string; label: string }[]>([])
  const { isLoading } = useFetch<ITag[]>({
    initialRun: true,
    config: {
      url: `/tags/${repo_id}`,
      method: 'GET',
    },
    onSuccess: data => {
      const tags: Option[] = Array.isArray(data) ? data.map(tag => ({ value: tag.name, label: tag.name })) : []
      setCurrentTags(prev => [...tags])
    },
  })

  const { run: createTag } = useFetch<string>({
    initialRun: false,
    config: {
      url: `/tags`,
      method: 'POST',
    },
  })

  const { run: deleteTag } = useFetch<string>({
    initialRun: false,
    config: {
      url: `/tags`,
      method: 'DELETE',
    },
  })

  return (
    <div className='flex w-full flex-col'>
      {isLoading ? (
        <div className='w-full h-9 bg-slate-300/40 rounded-xl animate-pulse'></div>
      ) : (
        <MultipleSelector
          creatable
          className='rounded-xl border bg-card text-card-foreground shadow'
          onChange={async value => {
            const remove = currentTags.filter(t => !value.find(v => v.value === t.value))
            const add = value.filter(v => !currentTags.find(t => t.value === v.value))

            await Promise.all([
              add.map(async tag => {
                await createTag({
                  payload: {
                    name: tag.value,
                    repo_id,
                  },
                })
                setCurrentTags(prev => [...prev, { value: tag.value, label: tag.value }])
              }),
              remove.map(async tag => {
                await deleteTag({
                  payload: {
                    name: tag.value,
                    repo_id,
                  },
                })
                setCurrentTags(currentTags.filter(t => t.value !== tag.value))
              }),
            ])
          }}
          defaultOptions={[]}
          value={currentTags}
          maxSelected={5}
          placeholder='trying to search tabs...'
          loadingIndicator={<p className='py-2 text-center text-lg leading-10 text-muted-foreground'>loading...</p>}
          emptyIndicator={
            <p className='w-full text-center text-lg leading-10 text-muted-foreground'>no results found.</p>
          }
        />
      )}
    </div>
  )
}

export default Tags
