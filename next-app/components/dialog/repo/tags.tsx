'use client'

import { useEffect, useState, useRef } from 'react'
import MultipleSelector, { Option } from '@/components/ui/multiple-selector'
import { getTags, getTagsByRepo, createTag, deleteTagByRepo } from '@/actions/neo4j'
import { useUser } from '@/context'

const Tags = (props: { repo_id: number }) => {
  const { repo_id } = props
  const { email } = useUser()
  const [getTagging, setGetTagging] = useState(true)
  const [currentTags, setCurrentTags] = useState<Option[]>([])

  useEffect(() => {
    const getCurrentTags = async () => {
      const tags = await getTagsByRepo(email, repo_id)
      setCurrentTags(tags.map(tag => ({ value: tag.id, label: tag.name })))
      setGetTagging(false)
    }

    getCurrentTags()
  }, [])

  const tags: Option[] = currentTags.map(tag => ({ value: tag.value, label: tag.label }))

  return (
    <div className='flex w-full flex-col'>
      {getTagging ? (
        <div className='w-full h-9 bg-slate-300 rounded-sm animate-pulse'></div>
      ) : (
        <MultipleSelector
          creatable
          onSearch={async () => {
            const tags = await getTags(email)
            return tags
              .map(tag => ({ value: tag.id, label: tag.name }))
              .filter(tag => !currentTags.find(t => t.value === tag.value))
          }}
          onChange={async value => {
            const remove = currentTags.filter(t => !value.find(v => v.value === t.value))
            const add = value.filter(v => !currentTags.find(t => t.value === v.value))

            await Promise.all([
              add.map(async tag => {
                const t = await createTag({ email, name: tag.label, repo_id })
                if (t) {
                  setCurrentTags(prev => [...prev, { value: t.id, label: t.name }])
                }
              }),
              remove.map(async tag => {
                await deleteTagByRepo(email, repo_id, tag.value)
                setCurrentTags(currentTags.filter(t => t.value !== tag.value))
              }),
            ])
          }}
          defaultOptions={[]}
          value={tags}
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
