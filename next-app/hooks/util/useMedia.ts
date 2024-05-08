import { useMediaQuery } from 'react-responsive'

const useMedia = () => {
  const isDesktopOrLaptop = useMediaQuery({
    query: '(min-width: 1280px)',
  })
  const isBigScreen = useMediaQuery({ query: '(min-width: 1824px)' })
  const isTabletOrMobile = useMediaQuery({ query: '(max-width: 1224px)' })
  const isPortrait = useMediaQuery({ query: '(orientation: portrait)' })
  const isRetina = useMediaQuery({ query: '(min-width: 1280px) and (max-width: 1480px)' })
  const isMobile = useMediaQuery({ query: '(max-width: 400px)' })

  return {
    isDesktopOrLaptop,
    isBigScreen,
    isTabletOrMobile,
    isPortrait,
    isRetina,
    isMobile,
  }
}

export default useMedia
