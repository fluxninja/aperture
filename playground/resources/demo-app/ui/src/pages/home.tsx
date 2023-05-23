import React, { FC } from 'react'
import { MonitorRequest } from '../components/monitor-request'
import { Box, Typography, styled } from '@mui/material'

// add a graceful-js error component to show error or success and add api calls to aperture configured endpoint

export const HomePage: FC = () => (
  <HomePageWrapper>
    <HomePageColumnBox>
      <MonitorRequest
        userType="Subscriber"
        requestRecord={[
          {
            isError: true,
          },
          {
            isError: false,
          },
          {
            isError: true,
          },
          {
            isError: false,
          },
          {
            isError: true,
          },
        ]}
      />
    </HomePageColumnBox>
    <HomePageColumnBox>
      <Typography
        variant="h5"
        sx={(theme) => ({ color: theme.palette.grey[400] })}
      >
        {/* add a graceful-js error component to show error or success */}
        200
      </Typography>
    </HomePageColumnBox>
  </HomePageWrapper>
)

export const HomePageWrapper = styled(Box)(({ theme }) => ({
  display: 'grid',
  gridTemplateColumns: '1fr 1fr',
  gridTemplateRows: 'auto',
  gap: theme.spacing(1),
  minHeight: 500,
  margin: '0px auto',
  width: '70%',
  padding: theme.spacing(2),
  [theme.breakpoints.down('sm')]: {
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    width: '100%',
  },
}))

export const HomePageColumnBox = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  justifyContent: 'center',
  gap: theme.spacing(2),
  minHeight: 500,
}))
