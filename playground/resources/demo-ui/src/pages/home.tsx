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

  // Request Spec for rate-limit endpoint with Executive user to retrieve fluxninja founders info
  const reqSpec: RequestSpec = {
    method: 'GET',
    endpoint: '/rate-limit',
    userType: 'executive',
    userId: 'CEO',
  }
  // Request to rate-limit endpoint with Guest user
  const {
    refetch,
    isError,
    requestRecord: requestRecord,
    isLoading: isLoadingRequest,
    data: requestResponse,
  } = useRequestToEndpoint(reqSpec)

  // Request to workload-prioritization endpoint with Guest user and Subscriber user
  const reqSpec2: RequestSpec = {
    method: 'POST',
    endpoint: '/workload-prioritization',
    userType: 'subscriber',
    userId: 'Test',
  }
  const {
    refetch: refetchSubscriber,
    isError: isErrorSubscriber,
    requestRecord: subscriberRequestRecord,
    isLoading: isLoadingSubscriber,
  } = useRequestToEndpoint(reqSpec2)
  reqSpec2.userType = 'guest'
  const {
    refetch: refetchGuest,
    isError: isErrorGuest,
    requestRecord: GuestRequestRecord,
    isLoading: isLoadingGuest,
  } = useRequestToEndpoint(reqSpec2)

  return (
    <TabContext value={value}>
      <TabList onChange={handleChange} aria-label="FluxNinja Scenarios" variant="fullWidth">
        <Tab label="Rate Limit" value="1"/>
        <Tab label="Workload Prioritization" value="2" />
      </TabList>
      <TabPanel value="1">
        <RequestMonitorPanel
          monitorRequestProps={{
            requestRecord: requestRecord,
            userType: 'FluxNinja Executive',
            refetch,
          }}
          isErrored={isError}
          isLoading={isLoadingRequest}
          errorComponentProps={{
            url: '/api' + reqSpec.endpoint,
            requestBody: {},
          }}
          responseData={requestResponse}
        />
      </TabPanel>
      <TabPanel value="2">
        <RequestMonitorPanel
          monitorRequestProps={{
            requestRecord: subscriberRequestRecord,
            refetch: refetchSubscriber,
            userType: 'Subscriber',
          }}
          isErrored={isErrorSubscriber}
          isLoading={isLoadingSubscriber}
          errorComponentProps={{
            url: '/api/' + reqSpec2.endpoint,
            requestBody: {},
          }}
          responseData={200}
        />
        <RequestMonitorPanel
          monitorRequestProps={{
            requestRecord: GuestRequestRecord,
            refetch: refetchGuest,
            userType: 'Guest',
          }}
          isErrored={isErrorGuest}
          isLoading={isLoadingGuest}
          errorComponentProps={{
            url: '/api' + reqSpec2.endpoint,
            requestBody: {},
          }}
          responseData={200}
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
    requestFnc: () => {
      if (reqSpec.method === 'POST') {
        return api.post(reqSpec.endpoint, reqSpec.body, {
          headers: {
            'User-Id': reqSpec.userId,
            'User-Type': reqSpec.userType,
          },
        })
      } else if (reqSpec.method === 'GET') {
        return api.get(reqSpec.endpoint, {
          headers: {
            'User-Id': reqSpec.userId,
            'User-Type': reqSpec.userType,
          },
        })
      } else {
        throw new Error(`Invalid method: ${reqSpec.method}`)
      }
    },
    options: {
      disabled: true,
    },
  });


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

  return { isError, refetch: startFetch, requestRecord, isLoading, data }
}

export interface RequestMonitorPanelProps {
  monitorRequestProps: MonitorRequestProps
  isErrored: boolean
  isLoading: boolean
  errorComponentProps: GracefulErrorProps
  responseData: any
}

export const RequestMonitorPanel: FC<RequestMonitorPanelProps> = ({
  monitorRequestProps,
  isErrored,
  isLoading,
  errorComponentProps,
  responseData,
}) => (
  <HomePageWrapper>
    <HomePageColumnBox>
      <MonitorRequest {...monitorRequestProps} />
    </HomePageColumnBox>
    <HomePageColumnBox>
      {isErrored && !isLoading ? (
        <GracefulError {...errorComponentProps} />
      ) : (
        <div
          sx={(theme) => ({
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
            alignItems: "center",
            color: theme.palette.grey[400],
          })}
        >
          <Typography variant="h5">{responseData?.data?.role}</Typography>
          <Typography variant="h5">{responseData?.data?.name}</Typography>
          <Typography variant="h5">{responseData?.data?.education}</Typography>
        </div>
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
  alignItems: 'center',
  flexDirection: 'column',
  justifyContent: 'center',
  gap: theme.spacing(2),
}))
