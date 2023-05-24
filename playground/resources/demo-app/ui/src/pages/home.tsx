import React, { FC, useEffect, useState } from 'react'
import {
  MonitorRequest,
  MonitorRequestProps,
  RequestRecord,
} from '../components/monitor-request'
import { Box, Typography, styled } from '@mui/material'
import { GracefulError, gracefulRequest } from '@fluxninja-tools/graceful-js'
import { api } from '../api'

// add a graceful-js error component to show error or success and add api calls to aperture configured endpoint

export const HomePage: FC = () => {
  const collection = useGracefulRequest()
  return (
    <RequestMonitorPanel
      monitorRequestProps={{
        requestRecord: collection,
        userType: 'Subscriber',
      }}
      isErrored={collection[collection.length - 1]?.isError}
    />
  )
}

export const useGracefulRequest = () => {
  const [requestCollection, setRequestCollection] = useState<RequestRecord[]>(
    []
  )

  useEffect(() => {
    gracefulRequest(
      'Axios',
      () => api.get('/rate-limit'),
      (err, res) => {
        if (err) {
          setRequestCollection((prev) => [
            ...prev,
            { isError: true, rateLimitInfo: err.rateLimitInfo },
          ])
          return
        }
        setRequestCollection((prev) => [
          ...prev,
          { isError: false, rateLimitInfo: res.rateLimitInfo },
        ])
      }
    )
      .then(() => {
        setRequestCollection((prev) => [...prev, { isError: false }])
      })
      .catch(() => {
        setRequestCollection((prev) => [...prev, { isError: true }])
      })
  }, [])

  return requestCollection
}

export interface RequestMonitorPanelProps {
  monitorRequestProps: MonitorRequestProps
  isErrored: boolean
}

export const RequestMonitorPanel: FC<RequestMonitorPanelProps> = ({
  monitorRequestProps,
  isErrored,
}) => (
  <HomePageWrapper>
    <HomePageColumnBox>
      <MonitorRequest {...monitorRequestProps} />
    </HomePageColumnBox>
    <HomePageColumnBox>
      {isErrored ? (
        <GracefulError />
      ) : (
        <Typography
          variant="h5"
          sx={(theme) => ({ color: theme.palette.grey[400] })}
          display="flex"
          justifyContent="center"
          alignItems="center"
        >
          200
        </Typography>
      )}
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
