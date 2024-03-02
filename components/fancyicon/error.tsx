import { useEffect, useRef } from 'react'
import { Player } from '@lordicon/react'
import ErrorJson from './wired-outline-1764-pushups.json'

export default function Error() {
  const playerRef = useRef<Player>(null)

  useEffect(() => {
    playerRef.current?.play()
  }, [])

  return <Player size={48} ref={playerRef} icon={ErrorJson} onComplete={() => playerRef.current?.playFromBeginning()} />
}
