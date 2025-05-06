import { useRouterState } from '@tanstack/react-router';

const Bars = [
  { name: '', path: 'dashboard' },
  { name: 'My Stars', path: 'stars' },
  { name: 'Collections', path: 'collections' },
];

function getFirstPath(pathname: string): string | null {
  const regex = /^\/([^/]+)/;
  const match = regex.exec(pathname);
  if (match) {
    return match[1];
  }
  return null;
}

export const useMenuName = () => {
  const { location } = useRouterState();

  const getName = () => {
    const path = getFirstPath(location.pathname) ?? '';
    const bar = Bars.find((bar) => bar.path === path);
    return bar?.name ?? '';
  };

  return { menuName: getName() };
};
