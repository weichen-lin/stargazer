import { Outlet, createRootRoute } from '@tanstack/react-router';
// import { TanStackRouterDevtools } from '@tanstack/react-router-devtools';

export const Route = createRootRoute({
  head: () => ({
    title: 'StarGazer',
    meta: [
      {
        name: 'description',
        content:
          'Organize and manage your GitHub starred repositories with tags and groups. Share your favorite collections with the world easily.',
      },
    ],
    links: [{ rel: 'icon', href: '/favicon.ico' }],
  }),
  component: () => <Outlet />,
});
