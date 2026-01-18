import { useState, useEffect } from 'react'

export function useDarkMode(initialValue = true) {
  const [darkMode, setDarkMode] = useState(initialValue)

  useEffect(() => {
    document.documentElement.classList.toggle('dark', darkMode)
  }, [darkMode])

  const toggle = () => setDarkMode(prev => !prev)

  return { darkMode, toggle }
}
