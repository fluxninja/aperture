import React, { FC, useCallback, useEffect, useState } from 'react'
import {
  MonitorRequest,
  MonitorRequestProps,
  RequestRecord,
} from '../components/monitor-request'
import { Box, Typography, styled, Tabs, Tab } from '@mui/material'
import { TabContext, TabList, TabPanel } from '@mui/lab'
import {
  GracefulError,
  GracefulErrorProps,
  useGracefulRequest,
} from '@fluxninja-tools/graceful-js'
import { api, RequestSpec } from '../api'

export const HomePage: FC = () => {
  const [value, setValue] = useState('1')

  const handleChange = (event: React.SyntheticEvent, newValue: string) => {
    setValue(newValue)
  }

  const reqSpec: RequestSpec = {
    method: 'POST',
    endpoint: '/rate-limit',
    userType: 'Guest',
    userId: 'Vu',
  }

  // Request to rate-limit endpoint with Guest user
  const {
    refetch,
    isError,
    requestRecord: requestRecord,
    isLoading: isLoadingRequest,
  } = makeRequestToEndpoint(reqSpec)

  // Request to workload-prioritization endpoint with Guest user
  const reqSpec2: RequestSpec = {
    method: 'POST',
    endpoint: '/workload-prioritization',
    userType: 'Guest',
    userId: 'Vu',
  }
  const {
    refetch: refetchSubscriber,
    isError: isErrorSubscriber,
    requestRecord: subscriberRequestRecord,
    isLoading: isLoadingSubscriber,
  } = makeRequestToEndpoint(reqSpec2)

  // Request to workload-prioritization endpoint with Subscriber user
  reqSpec2.userType = 'Subscriber'
  const {
    refetch: refetchSubscriber2,
    isError: isErrorSubscriber2,
    requestRecord: subscriberRequestRecord2,
    isLoading: isLoadingSubscriber2,
  } = makeRequestToEndpoint(reqSpec2)

  return (
    <TabContext value={value}>
      <TabList onChange={handleChange} aria-label="GracefulJS Tabs">
        <Tab label="Rate Limit" value="1" />
        <Tab label="Workload Prioritization" value="2" />
        <Tab label="TODO" value="3" />
      </TabList>
      <TabPanel value="1">
        <RequestMonitorPanel
          monitorRequestProps={{
            requestRecord: requestRecord,
            userType: 'Guest',
            refetch,
          }}
          isErrored={isError}
          isLoading={isLoadingRequest}
          errorComponentProps={{
            url: '/api/rate-limit',
            requestBody: {},
          }}
        />
      </TabPanel>
      <TabPanel value="2">
        <RequestMonitorPanel
          monitorRequestProps={{
            requestRecord: subscriberRequestRecord,
            refetch: refetchSubscriber,
            userType: 'Guest',
          }}
          isErrored={isErrorSubscriber}
          isLoading={isLoadingSubscriber}
          errorComponentProps={{
            url: '/api/workload-prioritization',
            requestBody: {},
          }}
        />
        <RequestMonitorPanel
          monitorRequestProps={{
            requestRecord: subscriberRequestRecord2,
            refetch: refetchSubscriber2,
            userType: 'Subscriber',
          }}
          isErrored={isErrorSubscriber2}
          isLoading={isLoadingSubscriber2}
          errorComponentProps={{
            url: '/api/workload-prioritization',
            requestBody: {},
          }}
        />
      </TabPanel>
      <TabPanel value="3">
        <RequestMonitorPanel
          monitorRequestProps={{
            requestRecord: subscriberRequestRecord2,
            refetch: refetchSubscriber2,
            userType: 'Subscriber',
          }}
          isErrored={isErrorSubscriber2}
          isLoading={isLoadingSubscriber2}
          errorComponentProps={{
            url: '/api/workload-prioritization',
            requestBody: {},
          }}
        />
      </TabPanel>
    </TabContext>
  )
}

export const makeRequestToEndpoint = (reqSpec: RequestSpec) => {
  const [requestRecord, setRequestRecord] = useState<RequestRecord[]>([]) // record state for each request
  const [requestCount, setRequestCount] = useState(0) // number of request count state
  const [intervalId, setIntervalId] = useState<NodeJS.Timeout | null>(null) // interval id state, used to clear interval

  const { isError, refetch, error, isRetry, isLoading } =
    useGracefulRequest<'Axios'>({
      typeOfRequest: 'Axios',
      requestFnc: () =>
        api.post(
          reqSpec.endpoint,
          {},
          {
            headers: {
              'User-Id': reqSpec.userId,
              'User-Type': reqSpec.userType,
            },
          }
        ),
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

export const HomePageTabs = styled(Tabs)(({ theme }) => ({
  width: '100%',
  minHeight: 500,
  display: 'flex',
  flexDirection: 'column',
  justifyContent: 'center',
  gap: theme.spacing(2),
}))
