import React, { FC, useCallback, useEffect, useState } from 'react'
import {
  MonitorRequest,
  MonitorRequestProps,
  RequestRecord,
} from '../components/monitor-request'
import { Box, styled, Tabs, Tab, CircularProgress } from '@mui/material'

import { TabContext, TabList, TabPanel } from '@mui/lab'
import {
  GracefulError,
  GracefulErrorProps,
  useGraceful,
  useGracefulRequest,
} from '@fluxninja-tools/graceful-js'
import { api, RequestSpec } from '../api'
import { SuccessIcon } from './success-icon'
import { useTabChange } from '../hooks'

const RATE_LIMIT_REQUEST: RequestSpec = {
  method: 'GET',
  url: '/api/rate-limit',
  headers: {
    'User-Id': 'DemoUI',
    'User-Type': 'user',
  },
}

const WORKLOAD_PRIORITIZATION_SUBSCRIBER_REQUEST: RequestSpec = {
  url: '/request',
  method: 'POST',
  data: {},
  headers: {
    'User-Id': 'DemoUI',
    'User-Type': 'subscriber',
  },
}

const WORKLOAD_PRIORITIZATION_GUEST_REQUEST: RequestSpec = {
  url: '/request',
  method: 'POST',
  data: {},
  headers: {
    'User-Id': 'DemoUI',
    'User-Type': 'guest',
  },
}

export const HomePage: FC = () => {
  const { handleChange, value } = useTabChange([
    'Rate Limit',
    'Workload Prioritization',
  ])

  const {
    refetch,
    isError,
    requestRecord,
    isLoading: isLoadingRequest,
  } = useRequestToEndpoint(RATE_LIMIT_REQUEST)

  const {
    refetch: refetchUser,
    isError: isErrorUser,
    requestRecord: userRequestRecord,
    isLoading: isLoadingUser,
  } = useRequestToEndpoint(WORKLOAD_PRIORITIZATION_SUBSCRIBER_REQUEST)

  const {
    refetch: refetchGuest,
    isError: isErrorGuest,
    requestRecord: guestRequestRecord,
    isLoading: isLoadingGuest,
  } = useRequestToEndpoint(WORKLOAD_PRIORITIZATION_GUEST_REQUEST)

  return (
    <TabContext value={value}>
      <TabList
        onChange={handleChange}
        aria-label="FluxNinja Scenarios"
        variant="fullWidth"
      >
        <Tab label="Rate Limit" value="0" />
        <Tab label="Workload Prioritization" value="1" />
      </TabList>
      <TabPanel value="0">
        <RequestMonitorPanel
          monitorRequestProps={{
            requestRecord: requestRecord,
            userType: 'Rate Limit',
            refetch,
          }}
          isErrored={isError}
          isLoading={isLoadingRequest}
          errorComponentProps={{
            url: '/aperture/api/rate-limit',
          }}
        />
      </TabPanel>
      <TabPanel value="1">
        <RequestMonitorPanel
          monitorRequestProps={{
            requestRecord: userRequestRecord,
            refetch: refetchUser,
            userType: 'Subscriber',
          }}
          isErrored={isErrorUser}
          isLoading={isLoadingUser}
          errorComponentProps={{
            url: '/aperture/request',
            requestBody: {},
          }}
        />
        <RequestMonitorPanel
          monitorRequestProps={{
            requestRecord: guestRequestRecord,
            refetch: refetchGuest,
            userType: 'Guest',
          }}
          isErrored={isErrorGuest}
          isLoading={isLoadingGuest}
          errorComponentProps={{
            url: '/aperture/request',
            requestBody: {},
          }}
        />
      </TabPanel>
    </TabContext>
  )
}

export const useRequestToEndpoint = (reqSpec: RequestSpec) => {
  const [requestRecord, setRequestRecord] = useState<RequestRecord[]>([]) // record state for each request
  const [requestCount, setRequestCount] = useState(0) // number of request count state
  const [intervalId, setIntervalId] = useState<NodeJS.Timeout | null>(null) // interval id state, used to clear interval

  const { isError, refetch, error, data, isRetry, isLoading } =
    useGracefulRequest<'Axios'>({
      typeOfRequest: 'Axios',
      requestFnc: () => api(reqSpec),
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
    }, 800)

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

  return { isError, refetch: startFetch, requestRecord, isLoading, data }
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
}) => {
  const { errorInfo } = useGraceful()
  return (
    <HomePageWrapper>
      <HomePageColumnBox>
        <MonitorRequest {...monitorRequestProps} />
      </HomePageColumnBox>
      <HomePageColumnBox>
        {isLoading ? (
          <FlexBox>
            <CircularProgress />
          </FlexBox>
        ) : isErrored && !isLoading ? (
          errorInfo.get(JSON.stringify(errorComponentProps))
            ?.errorComponent || <GracefulError {...errorComponentProps} /> // TODO: Fix error component in graceful-js
        ) : (
          <FlexBox>
            <SuccessIcon style={{ width: '15rem', height: '15rem' }} />
          </FlexBox>
        )}
      </HomePageColumnBox>
    </HomePageWrapper>
  )
}

const FlexBox = styled(Box)(() => ({
  display: 'flex',
  flexDirection: 'column',
  justifyContent: 'center',
  alignItems: 'center',
  minHeight: 500,
}))

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
    height: '100%',
  },
}))

export const HomePageColumnBox = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  justifyContent: 'center',
  gap: theme.spacing(2),
  minHeight: 500,
}))

export const HomePageTabs = styled(Tabs)(({ theme }) => ({
  width: '100%',
  minHeight: 500,
  display: 'flex',
  alignItems: 'center',
  flexDirection: 'column',
  justifyContent: 'center',
  gap: theme.spacing(2),
}))
