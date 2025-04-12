import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { ChevronUp } from 'lucide-react';
import {
  Menubar,
  MenubarCheckboxItem,
  MenubarContent,
  MenubarItem,
  MenubarMenu,
  MenubarRadioGroup,
  MenubarRadioItem,
  MenubarSeparator,
  MenubarShortcut,
  MenubarSub,
  MenubarSubContent,
  MenubarSubTrigger,
  MenubarTrigger,
} from '@/components/ui/menubar';
import { useUser } from '@clerk/clerk-react';

export default function UserInfo() {
  const { user } = useUser();

  return (
    <Menubar className=''>
      <MenubarMenu>
        <MenubarTrigger className='cursor-pointer group'>
          <Button variant='outline' className='bg-white group'>
            <div className='flex items-center gap-x-2'>
              <Avatar>
                <AvatarImage src={user?.imageUrl} />
                <AvatarFallback>WL</AvatarFallback>
              </Avatar>
              <span className='text-sm font-medium text-gray-900 dark:text-white'>
                {user?.username}
              </span>
              <ChevronUp className='ml-2 h-4 w-4 text-gray-500 transition-transform duration-200 group-[data=open]:rotate-0 group-data-[state=open]:rotate-180' />
            </div>
          </Button>
        </MenubarTrigger>
        <MenubarContent>
          <MenubarItem>
            New Tab <MenubarShortcut>⌘T</MenubarShortcut>
          </MenubarItem>
          <MenubarItem>
            New Window <MenubarShortcut>⌘N</MenubarShortcut>
          </MenubarItem>
          <MenubarItem disabled>New Incognito Window</MenubarItem>
          <MenubarSeparator />
          <MenubarSub>
            <MenubarSubTrigger>Share</MenubarSubTrigger>
            <MenubarSubContent>
              <MenubarItem>Email link</MenubarItem>
              <MenubarItem>Messages</MenubarItem>
              <MenubarItem>Notes</MenubarItem>
            </MenubarSubContent>
          </MenubarSub>
          <MenubarSeparator />
          <MenubarItem>
            Print... <MenubarShortcut>⌘P</MenubarShortcut>
          </MenubarItem>
        </MenubarContent>
      </MenubarMenu>
    </Menubar>
  );
}
