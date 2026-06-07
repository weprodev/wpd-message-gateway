import { RouterProvider } from "react-router-dom"

import { router } from "@/core/router"
import { ThemeProvider } from "@/shared/context/theme-context"

function App() {
  return (
    <ThemeProvider>
      <RouterProvider router={router} />
    </ThemeProvider>
  )
}

export default App
