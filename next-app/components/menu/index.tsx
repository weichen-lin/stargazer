'use client'

import MobileBar from './mobile'
import DesktopBar from './desktop'

export { MobileBar, DesktopBar }

export default function Menu() {
  return (
    <>
      <MobileBar />
      <DesktopBar />
    </>
  )
}
