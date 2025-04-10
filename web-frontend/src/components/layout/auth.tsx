import { Link } from '@tanstack/react-router';
import { Sidebar } from '@/components/shared';
import Menu from '@/components/shared/menu';

export default async function AuthLayout(props: { children: React.ReactNode }) {
  const { children } = props;

  return (
    <div className='grid grid-cols-[260px_1fr] h-screen overflow-x-hidden'>
      <Sidebar />
      <div className='flex flex-col h-screen overflow-hidden'>
        <header className='bg-white shadow'>
          <Menu />
        </header>
        <main className='bg-gray-200 p-4 flex-grow overflow-y-auto'>
          {children}
        </main>
      </div>
    </div>
  );
}
