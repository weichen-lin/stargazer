import { createFileRoute } from '@tanstack/react-router';
import { Search } from '@/pages/stars';

export const Route = createFileRoute('/_auth/stars')({
  component: Stars,
});

function Stars() {
  return (
    <div className='grid grid-cols-[300px_1fr] gap-4 overflow-hidden h-full'>
      <Search />
      <div className='h-full grid grid-rows-4 grid-cols-3 grid-flow-row gap-4'></div>
    </div>
  );
}
