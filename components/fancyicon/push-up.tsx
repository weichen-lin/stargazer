import { useEffect, useRef } from 'react'
import { Player } from '@lordicon/react'
import PushUpJson from './wired-outline-1764-pushups.json'

export default function PushUp() {
  const playerRef = useRef<Player>(null)

  useEffect(() => {
    playerRef.current?.play()
  }, [])

  return (
    <Player size={48} ref={playerRef} icon={PushUpJson} onComplete={() => playerRef.current?.playFromBeginning()} />
  )
}
