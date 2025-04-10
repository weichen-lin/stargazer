import MultipleSelector from '@/components/ui/multiple-selector';

export default function Search() {
  return (
    <div className='flex flex-col gap-2 items-center'>
      <MultipleSelector
        value={[]}
        onChange={(e) => {}}
        defaultOptions={[]}
        placeholder='Select languages you like...'
        emptyIndicator={
          <p className='text-center text-lg leading-10 text-gray-600 dark:text-gray-400'>
            no results found.
          </p>
        }
        className='border-slate-900 w-[380px]'
        disabled={false}
      />
    </div>
  );
}
