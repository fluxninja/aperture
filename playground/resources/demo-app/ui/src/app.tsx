import React, { FC } from 'react'
import { DemoAppThemeProvider } from './theme'
import { HomePage } from './pages'
import {
  GracefulProvider,
  Config as GracefulJsConfig,
} from '@fluxninja-tools/graceful-js'
import { api } from './api'

const gracefulJsConfig: GracefulJsConfig = {
  axios: api,
}

export const App: FC = () => {
  return (
    <DemoAppThemeProvider>
      <GracefulProvider config={gracefulJsConfig}>
        <HomePage />
      </GracefulProvider>
    </DemoAppThemeProvider>
  )
}
