import React, { FC, useCallback, useEffect, useState } from 'react'
import {
  MonitorRequest,
  MonitorRequestProps,
  RequestRecord,
} from '../components/monitor-request'
import { Box, Typography, styled } from '@mui/material'
import {
  GracefulError,
  GracefulErrorProps,
  useGracefulRequest,
} from '@fluxninja-tools/graceful-js'
import { api } from '../api'

export const HomePage: FC = () => {
  const {
    refetch,
    isError,
    requestRecord: crawlerRequestRecord,
    isLoading: isLoadingCrawler,
  } = useGracefulRequestForRateLimit('Crawler')
  const {
    refetch: refetchSubscriber,
    isError: isErrorSubscriber,
    requestRecord: subscriberRequestRecord,
    isLoading: isLoadingSubscriber,
  } = useGracefulRequestForRateLimit('Subscriber')

  return (
    <>
      <RequestMonitorPanel
        monitorRequestProps={{
          requestRecord: crawlerRequestRecord,
          userType: 'Crawler',
          refetch,
        }}
        isErrored={isError}
        isLoading={isLoadingCrawler}
        errorComponentProps={{
          url: '/api/rate-limit',
          requestBody: {},
        }}
      />
      <RequestMonitorPanel
        monitorRequestProps={{
          requestRecord: subscriberRequestRecord,
          refetch: refetchSubscriber,
          userType: 'Subscriber',
        }}
        isErrored={isErrorSubscriber}
        isLoading={isLoadingSubscriber}
        errorComponentProps={{
          url: '/api/rate-limit',
          requestBody: {},
        }}
      />
    </>
  )
}

export const useGracefulRequestForRateLimit = (userID: string) => {
  const [requestRecord, setRequestRecord] = useState<RequestRecord[]>([]) // record state for each request
  const [requestCount, setRequestCount] = useState(0) // number of request count state
  const [intervalId, setIntervalId] = useState<NodeJS.Timeout | null>(null) // interval id state, used to clear interval

  const { isError, refetch, error, isRetry, isLoading } =
    useGracefulRequest<'Axios'>({
      typeOfRequest: 'Axios',
      requestFnc: () =>
        api.post(`/rate-limit`, {}, { headers: { 'user-id': userID } }),
      options: {
        disabled: true,
      },
    })

  // update record state if request counter is not 0
  const updateRecord = useCallback(() => {
    if (!requestCount) {
      return
    }
    setRequestRecord((prevErrors) => [
      ...prevErrors,
      { isError, rateLimitInfo: error?.rateLimitInfo || null, isRetry },
    ])
  }, [isError, requestCount, error?.rateLimitInfo, isRetry])

  // start making request after 800ms
  const startFetch = useCallback(() => {
    const intervalId = setInterval(() => {
      setRequestCount((prevCount) => prevCount + 1)
      refetch()
    }, 400)

    setIntervalId(intervalId)

    return () => {
      clearInterval(intervalId)
    }
  }, [refetch])

  // stop making request if isError is true or requestCount is greater than 50
  useEffect(() => {
    updateRecord()
    if (!intervalId) {
      return
    }
    if (isError) {
      clearInterval(intervalId)
      return
    }
    if (requestCount >= 60) {
      clearInterval(intervalId)
      return
    }
  }, [requestCount, intervalId, isError])

  return { isError, refetch: startFetch, requestRecord, isLoading }
}

export interface RequestMonitorPanelProps {
  monitorRequestProps: MonitorRequestProps
  isErrored: boolean
  isLoading: boolean
  errorComponentProps: GracefulErrorProps
}

export const RequestMonitorPanel: FC<RequestMonitorPanelProps> = ({
  monitorRequestProps,
  isErrored,
  isLoading,
  errorComponentProps,
}) => (
  <HomePageWrapper>
    <HomePageColumnBox>
      <MonitorRequest {...monitorRequestProps} />
    </HomePageColumnBox>
    <HomePageColumnBox>
      {isErrored && !isLoading ? (
        <GracefulError {...errorComponentProps} /> // TODO: not rendering right error component. Only default error component is rendering
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
