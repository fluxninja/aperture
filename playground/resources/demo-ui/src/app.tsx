import React, { FC } from 'react'
import { DemoAppThemeProvider } from './theme'
import { HomePage } from './pages'
import {
  GracefulProvider,
  Config as GracefulJsConfig,
} from '@fluxninja-tools/graceful-js'
import { BrowserRouter } from 'react-router-dom'
import { api } from './api'

const gracefulJsConfig: GracefulJsConfig = {
  axios: api,
  maxBackOffTime: 25,
  maxRequestResolveTime: 10,
}

export const App: FC = () => {
  return (
    <DemoAppThemeProvider>
      <GracefulProvider config={gracefulJsConfig}>
        <BrowserRouter>
          <HomePage />
        </BrowserRouter>
      </GracefulProvider>
    </DemoAppThemeProvider>
  )
}
