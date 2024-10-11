import { ICollection } from '@/client/collection'
import { MixIcon } from '@radix-ui/react-icons'
import { useSearchCollection } from '@/app/stars/hook'
import { cn } from '@/lib/utils'
import { LockClosedIcon, LockOpen2Icon } from '@radix-ui/react-icons'
import { Button } from '@/components/ui/button'
import { useStars } from '@/hooks/stars'
import { useFetch } from '@/hooks/util'

export default function ListCollection(props: ICollection) {
  const { id, name, is_public } = props
  const { chosen, setChosen, open, setOpen } = useSearchCollection()
  const { selectedRepo } = useStars()
  const { run, isLoading } = useFetch({
    initialRun: false,
    config: {
      url: `/collection/repos/${id}`,
      method: 'POST',
    },
    onSuccess: e => {
      setOpen(false)
    },
  })

  const addRepos = (repo_ids: number[]) => {
    run({ payload: { repo_ids } })
  }

  return (
    <div
      className={cn(
        'flex justify-between w-full px-4 cursor-pointer py-3 group overflow-x-hidden',
        chosen === id ? 'bg-blue-100' : 'hover:bg-slate-100',
      )}
      onClick={() => {
        setChosen(chosen === id ? null : id)
      }}
      onMouseEnter={() => {}}
      onMouseLeave={() => {}}
    >
      <div className='flex gap-x-4 items-center'>
        <MixIcon className='w-5 h-5' />
        <div>{name}</div>
      </div>
      <div className='flex gap-x-4 items-center'>
        {is_public ? (
          <LockOpen2Icon className='w-5 h-5 text-green-500' />
        ) : (
          <LockClosedIcon className='w-5 h-5 text-red-500' />
        )}
        <div
          className={cn(
            'group-hover:opacity-100 group-hover:translate-x-0',
            'transition-all duration-300 transform',
            chosen === id ? 'translate-x-0' : 'translate-x-[-100%] opacity-0',
          )}
        >
          <Button
            size='sm'
            variant='outline'
            onClick={() => {
              addRepos(selectedRepo)
            }}
            loading={isLoading}
          >
            Add
          </Button>
        </div>
      </div>
    </div>
  )
}
