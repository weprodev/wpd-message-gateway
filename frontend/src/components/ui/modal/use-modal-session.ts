import { useEffect, useRef, useState } from "react"

/**
 * Bumps a session counter whenever the modal opens and runs onOpen.
 * Use the returned ref to ignore stale async results after close/reopen.
 */
export function useModalSession(isOpen: boolean, onOpen: () => void) {
  const [prevIsOpen, setPrevIsOpen] = useState(isOpen)
  const [session, setSession] = useState(0)
  const sessionRef = useRef(session)

  useEffect(() => {
    sessionRef.current = session
  }, [session])

  if (isOpen !== prevIsOpen) {
    setPrevIsOpen(isOpen)
    if (isOpen) {
      setSession((current) => current + 1)
      onOpen()
    }
  }

  return sessionRef
}
