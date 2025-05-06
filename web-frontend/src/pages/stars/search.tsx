import MultipleSelector, {
  type Option,
} from '@/components/ui/multiple-selector';
import { Download, Star, FolderPlus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import useLanguageDistribution from '@/apis/repository/useLanguageDistribution';
import { useState } from 'react';

export default function Search() {
  const { data, isLoading } = useLanguageDistribution();
  const [languages, setLanguages] = useState<Option[]>([]);

  const options: Option[] = data
    .filter((e) => e.language !== '')
    .map((e) => {
      return {
        label: e.language,
        value: e.language,
      };
    });

  return (
    <div className='grid grid-rows-[1fr_70px] w-full overflow-hidden'>
      <div className='flex flex-col gap-y-2 overflow-y-hidden justify-start items-start'>
        <MultipleSelector
          value={[]}
          onChange={(e) => {}}
          options={options}
          placeholder='Select languages you like...'
          emptyIndicator={
            <p className='text-center text-lg leading-10 text-gray-600 dark:text-gray-400 w-[260px]'>
              no results found.
            </p>
          }
          commandProps={{
            className: 'w-[260px]',
          }}
          className='border-slate-900 relative'
          disabled={false}
        />
      </div>
      <div className='flex items-end gap-x-4 justify-start'>
        <Button className=''>
          <div className='flex items-center gap-x-2'>
            <Download className='h-4 w-4' />
            <span>APPLY</span>
          </div>
        </Button>
        <Button variant='outline' className='flex items-center gap-2'>
          <div className='flex gap-x-2 items-center'>
            <FolderPlus className='h-4 w-4' />
            <span>RESET</span>
          </div>
        </Button>
      </div>
    </div>
  );
}
