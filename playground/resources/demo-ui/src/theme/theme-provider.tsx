import {
  CssBaseline,
  GlobalStyles,
  GlobalStylesProps,
  ThemeProvider,
  createTheme,
} from '@mui/material'
import React, { FC, PropsWithChildren } from 'react'

const theme = createTheme({
  typography: {
    fontFamily: "'Nunito', sans-serif",
  },
  palette: {
    mode: 'dark',
    primary: {
      main: '#8EC4AD',
    },
    secondary: {
      main: '#F27A40',
    },
  },
})

const globalStyles: GlobalStylesProps['styles'] = {
  html: {
    fontFamily: "'Nunito', sans-serif",
  },
  body: {
    boxSizing: 'border-box',
    margin: 0,
    padding: 0,
  },
}

export const DemoAppThemeProvider: FC<PropsWithChildren> = ({ children }) => (
  <ThemeProvider theme={theme}>
    <GlobalStyles styles={globalStyles} />
    <CssBaseline />
    {children}
  </ThemeProvider>
)
