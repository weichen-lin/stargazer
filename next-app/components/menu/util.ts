import { usePathname } from 'next/navigation'

const Bars = [
  { name: 'Dashboard', path: 'dashboard' },
  { name: 'My Stars', path: 'stars' },
  { name: 'Start Chat', path: 'chats' },
]

function getFirstPath(pathname: string): string | null {
  const regex = /^\/([^/]+)/
  const match = regex.exec(pathname)
  if (match) {
    return match[1]
  }
  return null
}

export const useMenuName = () => {
  const pathname = usePathname()

  const getName = () => {
    const path = getFirstPath(pathname) ?? ''
    const bar = Bars.find(bar => bar.path === path)
    return bar?.name ?? ''
  }

  return { menuName: getName() }
}
