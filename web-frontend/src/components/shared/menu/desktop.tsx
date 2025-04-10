// import { ModeToggle } from '@/components/provider';
import { motion } from 'framer-motion';
import { CommandSearch } from '@/components/shared/menu/search';
import UserInfo from '@/components/shared/menu/user-info';
import { useMenuName } from './util';
import { cn } from '@/lib/utils';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';

// import { Detail } from '@/components/dialog';
// import { useRepoDetail } from '@/hooks/util';

const DesktopBar = () => {
  const { menuName } = useMenuName();
  // const { open } = useRepoDetail();

  return (
    <motion.div
      initial={{ x: 80 }}
      animate={{ x: 0 }}
      className={cn(
        'flex-col items-center justify-between hidden lg:flex h-[57px] bg-background',
        'gap-y-6 backdrop-blur-md border-b-[1px] border-secondary px-4',
        'dark:bg-black dark:border-slate-800 dark:text-white w-full'
      )}
    >
      <div className='flex justify-end w-full h-full items-center gap-x-4'>
        <CommandSearch />
        <UserInfo />
      </div>
    </motion.div>
  );
};

export default DesktopBar;
