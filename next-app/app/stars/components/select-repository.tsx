import { Plus, FolderInput, Trash2 } from 'lucide-react'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { useStars } from '@/hooks/stars'

export default function SelectRepository() {
  const { selectedRepo, setSelectedRepo } = useStars()

  return (
    selectedRepo.length > 0 && (
      <div className='bg-slate-200 h-10 p-2 rounded-lg flex gap-x-4 items-center'>
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger asChild>
              <Plus
                className='w-7 h-7 rotate-45 hover:bg-slate-300 p-1 rounded-full cursor-pointer'
                onClick={() => {
                  setSelectedRepo([])
                }}
              />
            </TooltipTrigger>
            <TooltipContent>
              <p>Cancel Select</p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
        <div className='text-sm'>{selectedRepo.length} Selected</div>
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger asChild>
              <FolderInput className='w-7 h-7 hover:bg-slate-300 p-[5px] rounded-full cursor-pointer' />
            </TooltipTrigger>
            <TooltipContent>
              <p>Add to collection</p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger asChild>
              <Trash2 className='w-7 h-7 hover:bg-slate-300 p-[5px] rounded-full cursor-pointer' />
            </TooltipTrigger>
            <TooltipContent>
              <p>Delete Repositories</p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </div>
    )
  )
}
