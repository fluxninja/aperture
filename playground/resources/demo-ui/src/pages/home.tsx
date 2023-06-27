import React, { FC, useCallback, useState } from 'react'
import {
  MonitorRequest,
  MonitorRequestProps,
  RequestRecord,
} from '../components/monitor-request'
import { Box, Typography, styled } from '@mui/material'
import { GracefulError, useGracefulRequest } from '@fluxninja-tools/graceful-js'
import { api } from '../api'

export const HomePage: FC = () => {
  const {
    refetch,
    isError,
    requestRecord: crawlerRequestRecord,
  } = useGracefulRequestForRateLimit('Crawler')
  const {
    refetch: refetchSubscriber,
    isError: isErrorSubscriber,
    requestRecord: subscriberRequestRecord,
  } = useGracefulRequestForRateLimit('Subscriber')

  return (
    <>
      <RequestMonitorPanel
        monitorRequestProps={{
          requestRecord: crawlerRequestRecord,
          userType: 'Crawler',
          refetch,
        }}
        gracefulError={
          <GracefulError
            {...{
              url: `${window.location.origin}/api/rate-limit`,
            }}
          />
        }
        isErrored={isError}
      />
      <RequestMonitorPanel
        monitorRequestProps={{
          requestRecord: subscriberRequestRecord,
          userType: 'Subscriber',
          refetch: refetchSubscriber,
        }}
        gracefulError={
          <GracefulError
            {...{
              url: `${window.location.origin}/api/rate-limit`,
            }}
          />
        }
        isErrored={isErrorSubscriber}
      />
    </>
  )
}

export const useGracefulRequestForRateLimit = (userID: string) => {
  const [requestRecord, setRequestRecord] = useState<RequestRecord[]>([])

  const { isError, refetch, error } = useGracefulRequest<'Axios'>({
    typeOfRequest: 'Axios',
    requestFnc: () =>
      api.post(`/rate-limit`, {
        headers: { 'Content-Type': 'application/json', 'user-id': userID },
        data: {},
      }),
    options: {
      disabled: true,
    },
  })

  const updateRecord = useCallback(() => {
    setRequestRecord((prevErrors) => [
      ...prevErrors,
      { isError, rateLimitInfo: error?.rateLimitInfo || null },
    ])
  }, [isError, error])

  const startFetch = useCallback(async () => {
    loop: for (let i = 0; i < 60; i++) {
      await new Promise((resolve) => setTimeout(resolve, 800))

      refetch()
      updateRecord()
      // TODO: fix loop to break on error and request per second
      if (isError) {
        break loop
      }
    }
  }, [refetch, isError, updateRecord])

  return { isError, refetch: startFetch, requestRecord }
}

export interface RequestMonitorPanelProps {
  monitorRequestProps: MonitorRequestProps
  isErrored: boolean
  gracefulError: JSX.Element
}

export const RequestMonitorPanel: FC<RequestMonitorPanelProps> = ({
  monitorRequestProps,
  isErrored,
  gracefulError,
}) => (
  <HomePageWrapper>
    <HomePageColumnBox>
      <MonitorRequest {...monitorRequestProps} />
    </HomePageColumnBox>
    <HomePageColumnBox>
      {isErrored ? (
        gracefulError
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
