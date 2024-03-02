import { useEffect, useRef } from 'react'
import { Player } from '@lordicon/react'
import CheckupJson from './wired-outline-37-approve-checked-simple.json'

export default function Check() {
  const playerRef = useRef<Player>(null)

  useEffect(() => {
    playerRef.current?.play()
  }, [])

  return <Player size={32} ref={playerRef} icon={CheckupJson} />
}
