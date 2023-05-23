import React, { FC } from 'react'
import { DemoAppThemeProvider } from './theme'
import { HomePage } from './pages'

export const App: FC = () => {
  return (
    <DemoAppThemeProvider>
      <HomePage />
    </DemoAppThemeProvider>
  )
}
