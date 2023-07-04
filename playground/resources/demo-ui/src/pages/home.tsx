import React, { FC } from 'react'
import {
  MonitorRequest,
  MonitorRequestProps,
} from '../components/monitor-request'
import { Box, styled, Tab, CircularProgress } from '@mui/material'

import { TabContext, TabList, TabPanel } from '@mui/lab'
import {
  GracefulErrorByStatus,
  RateLimitInfo,
} from '@fluxninja-tools/graceful-js'
import { RequestSpec } from '../api'
import { SuccessIcon } from './success-icon'
import { useRequestToEndpoint, useTabChange } from '../hooks'
import { AxiosError } from 'axios'

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
    error: rateLimitRequestError,
  } = useRequestToEndpoint(RATE_LIMIT_REQUEST)

  const {
    refetch: refetchUser,
    isError: isErrorUser,
    requestRecord: userRequestRecord,
    isLoading: isLoadingUser,
    error: workloadPrioritySubscriberError,
  } = useRequestToEndpoint(WORKLOAD_PRIORITIZATION_SUBSCRIBER_REQUEST)

  const {
    refetch: refetchGuest,
    isError: isErrorGuest,
    requestRecord: guestRequestRecord,
    isLoading: isLoadingGuest,
    error: workloadPriorityGuestError,
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
          error={rateLimitRequestError}
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
          error={workloadPrioritySubscriberError}
        />
        <RequestMonitorPanel
          monitorRequestProps={{
            requestRecord: guestRequestRecord,
            refetch: refetchGuest,
            userType: 'Guest',
          }}
          isErrored={isErrorGuest}
          isLoading={isLoadingGuest}
          error={workloadPriorityGuestError}
        />
      </TabPanel>
    </TabContext>
  )
}

export interface RequestMonitorPanelProps {
  monitorRequestProps: MonitorRequestProps
  isErrored: boolean
  isLoading: boolean
  error: AxiosError<unknown, unknown> & {
    rateLimitInfo?: RateLimitInfo
  }
}

export const RequestMonitorPanel: FC<RequestMonitorPanelProps> = ({
  monitorRequestProps,
  isErrored,
  isLoading,
  error,
}) => {
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
          <GracefulErrorByStatus status={error.response?.status} />
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
