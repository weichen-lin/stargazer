'use client'

import { ISuggestion } from '@/actions'
import { motion } from 'framer-motion'
import clsx from 'clsx'
import { useRepoDetail } from '@/hooks/util'

function FullTextIndex(text: string, target: string): number[] | null {
  let left = 0
  let right = 0
  const result: number[] = []
  const t = text.toLowerCase()
  const q = target.toLowerCase()

  while (right < t.length && left < q.length) {
    if (t[right] === q[left]) {
      result.push(right)
      left++
    }
    right++
  }

  return result.length === q.length ? result : null
}

const FullTextSearchResult = (props: ISuggestion & { query: string; close: () => void }) => {
  const { repo_id, avatar_url, full_name, description, query, close } = props
  const { setRepoID, setOpen } = useRepoDetail()

  const full_name_index = FullTextIndex(full_name, query)
  const description_index = FullTextIndex(description ?? '', query)

  return (
    <motion.div
      whileHover={{ scale: 1.02, x: -5 }}
      className={clsx(
        'flex flex-col items-start gap-y-2',
        'hover:border-[1px] hover:border-slate-300 hover:shadow-md',
        'p-4 bg-white rounded-lg transition-colors w-[95%] mx-auto cursor-pointer',
      )}
      onClick={() => {
        setRepoID(repo_id)
        setOpen(true)
        close()
      }}
    >
      <div className='flex gap-x-2 items-center'>
        <img src={avatar_url} alt={full_name} className='w-6 h-6 rounded-md' />
        <div className='font-semibold text-blue-700 line-clamp-2'>
          {full_name_index ? (
            <>
              {full_name.split('').map((char, i) => {
                if (full_name_index.includes(i)) {
                  return (
                    <span className='bg-yellow-100' key={`char_full_name_${i}`}>
                      {char}
                    </span>
                  )
                }
                return char
              })}
            </>
          ) : (
            full_name
          )}
        </div>
      </div>
      <div className='w-full text-ellipsis text-slate-400 text-sm break-words'>
        {description_index && description ? (
          <>
            {description.split('').map((char, i) => {
              if (description_index.includes(i)) {
                return (
                  <span className='bg-yellow-100' key={`char_description_${i}`}>
                    {char}
                  </span>
                )
              }
              return char
            })}
          </>
        ) : (
          description ?? ''
        )}
      </div>
    </motion.div>
  )
}

export default FullTextSearchResult
