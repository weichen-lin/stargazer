import { cn } from '@/lib/utils';
import { createFileRoute } from '@tanstack/react-router';
import TypewriterEffectSmooth from '@/components/ui/typewriter-effect';
import { Button } from '@/components/ui/button';
import { HomeIcon } from 'lucide-react';
import { GitHubLogoIcon } from '@radix-ui/react-icons';
import { SignInButton, useSession, useUser } from '@clerk/clerk-react';
import { Link } from '@tanstack/react-router';

export const Route = createFileRoute('/')({
  component: Home,
});

function Home() {
  const { isLoaded, isSignedIn } = useSession();

  const words = [
    {
      text: 'Start',
      className: 'text-xl lg:text-3xl',
    },
    {
      text: 'managing',
      className: 'text-xl lg:text-3xl',
    },
    {
      text: 'your',
      className: 'text-xl lg:text-3xl',
    },
    {
      text: 'stars',
      className: 'text-xl lg:text-3xl',
    },
    {
      text: 'with',
      className: 'text-xl lg:text-3xl',
    },
    {
      text: 'StarGazer.',
      className: 'text-xl text-blue-500 dark:text-blue-500 lg:text-3xl',
    },
  ];

  return (
    <main className='h-screen p-6 overflow-hidden'>
      <div
        className={cn(
          'w-full max-w-[1024px] mx-auto h-full',
          'flex flex-col justify-between'
        )}
      >
        <div className='flex items-center justify-between'>
          <div className='flex gap-x-8 items-center'>
            <img
              src='/icon.jpeg'
              width={60}
              height={60}
              className='rounded-full'
              alt='stargazer logo'
            />
            <div className='text-4xl'>StarGazer</div>
          </div>
          {/* <ModeToggle /> */}
        </div>
        <div className='flex justify-between items-start md:items-center flex-col lg:flex-row gap-6'>
          <div className='flex flex-col justify-start'>
            <TypewriterEffectSmooth
              words={words}
              cursorClassName='h-5'
              className='flex items-center'
            />
            <p className='leading-7'>
              Discover the future of star management beyond GitHub Stars. Unlock
              the full potential of innovative, data-driven solutions that
              elevate how you manage and track your favorite repositories – this
              is just the beginning!
            </p>
          </div>
          <img
            src='/home.png'
            width={400}
            height={400}
            className='mx-auto'
            alt='home pic'
          />
        </div>
        {isLoaded && isSignedIn ? (
          <Link to='/dashboard'>
            <Button
              className='flex gap-x-4 max-w-[320px] mx-auto'
              // onClick={handleSignIn}
            >
              <GitHubLogoIcon />
              Dashboard
            </Button>
          </Link>
        ) : (
          <SignInButton
            mode='modal'
            forceRedirectUrl='/dashboard'
            signUpForceRedirectUrl='/dashboard'
            signUpFallbackRedirectUrl='/dashboard'
          >
            <Button
              className='flex gap-x-4 max-w-[320px] mx-auto'
              // onClick={handleSignIn}
            >
              <GitHubLogoIcon />
              Sign in with Github
            </Button>
          </SignInButton>
        )}
        <div className='text-center'>
          © WeiChen Lin {new Date().getFullYear()}
        </div>
      </div>
    </main>
  );
}
