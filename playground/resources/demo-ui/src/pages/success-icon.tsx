import { SvgIcon, SvgIconProps } from '@mui/material'
import React, { FC, PropsWithRef } from 'react'

export const SuccessIcon: FC<PropsWithRef<SvgIconProps>> = React.forwardRef(
  (props, ref) => (
    <SvgIcon width="100" height="100" viewBox="0 0 49 44" {...props} ref={ref}>
      <path
        className="cls-1"
        d="M0,68.07s80.09,33.28,200.09-24.59C219.37,33.36,249.65,7.35,281.21,0"
        fill="#F8773D"
      />
      <path
        className="cls-2"
        d="M200.09,43.48s-37.74,21.09-79.45-.14c0,0-36.5-17.24-79.4-.22,0,0-32.08,16.88-41.24,25.95"
        fill="#56AE89"
      />
      <circle className="cls-3" cx="100.19" cy="26.68" r="23.66" fill="#EFEEED"/>
      <circle className="cls-3" cx="172.45" cy="26.68" r="23.66" fill="#EFEEED"/>
      <path
        className="cls-2"
        d="M90.92,29.59a2.77,2.77,0,1,1-2.77-2.77A2.77,2.77,0,0,1,90.92,29.59Z"
        fill="#56AE89"
      />
      <path
        className="cls-2"
        d="M163.18,29.59a2.77,2.77,0,1,1-2.77-2.77A2.77,2.77,0,0,1,163.18,29.59Z"
        fill="#56AE89"
      />
      <circle className="cls-4" cx="136.73" cy="73.73" r="28.34" fill="#F8773D"/>
      <path
        className="cls-5"
        d="M141.61,72.81l-8.4,8.22-3.77-3.68"
        fill="#EFEEED"
      />
    </SvgIcon>
  )
)

SuccessIcon.displayName = 'SuccessIcon'

