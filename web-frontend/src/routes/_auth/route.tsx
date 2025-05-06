import { createFileRoute, Outlet } from '@tanstack/react-router';
import { Sidebar } from '@/components/shared';
import Menu from '@/components/shared/menu';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Toaster } from 'sonner';

const PUBLISHABLE_KEY = import.meta.env.VITE_CLERK_PUBLISHABLE_KEY;

if (!PUBLISHABLE_KEY) {
  throw new Error('Missing Publishable Key');
}

export const Route = createFileRoute('/_auth')({
  component: RouteComponent,
});

const queryClient = new QueryClient();

function RouteComponent() {
  return (
    <QueryClientProvider client={queryClient}>
      <Toaster richColors />
      <div className='grid grid-cols-[260px_1fr] h-screen overflow-x-hidden'>
        <Sidebar />
        <div className='flex flex-col h-screen overflow-hidden'>
          <header className='bg-white shadow'>
            <Menu />
          </header>
          <main className='bg-gray-200 p-4 flex-grow overflow-y-auto'>
            <Outlet />
          </main>
        </div>
      </div>
    </QueryClientProvider>
  );
}
