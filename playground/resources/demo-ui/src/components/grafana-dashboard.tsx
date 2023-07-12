import React, { FC } from 'react'

export interface GrafanaDashboardProps {
  src: string
}

export const GrafanaDashboard: FC<GrafanaDashboardProps> = ({ src }) => {
  return (
    <iframe style={{ width: '100%', height: '600px', border: 0 }} src={src} />
  )
}
