import React, { FC, useMemo } from 'react'
import {
  MonitorRequest,
  MonitorRequestProps,
} from '../components/monitor-request'
import { Box, styled, Tab, CircularProgress } from '@mui/material'

import { TabContext, TabList, TabPanel } from '@mui/lab'
import { RequestSpec } from '../api'
import { SuccessIcon } from './success-icon'
import { useRequestToEndpoint, useTabChange } from '../hooks'

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

const WAITING_ROOM_REQUEST: RequestSpec = {
  url: '/request',
  method: 'POST',
  data: {
    'User-Id': 'DemoUI',
    'User-Type': 'guest',
  },
}

export const HomePage: FC = () => {
  const { handleChange, value } = useTabChange([
    'Waiting Room',
    'Rate Limit',
    'Workload Prioritization',
  ])

  const {
    refetch,
    requestRecord,
    isLoading: isLoadingRequest,
  } = useRequestToEndpoint(RATE_LIMIT_REQUEST)

  const {
    refetch: refetchUser,

    requestRecord: userRequestRecord,
    isLoading: isLoadingUser,
  } = useRequestToEndpoint(WORKLOAD_PRIORITIZATION_SUBSCRIBER_REQUEST)

  const {
    refetch: refetchGuest,
    requestRecord: guestRequestRecord,
    isLoading: isLoadingGuest,
  } = useRequestToEndpoint(WORKLOAD_PRIORITIZATION_GUEST_REQUEST)

  const {
    refetch: refetchWaitingRoom,
    requestRecord: waitingRoomRequestRecord,
    isLoading: isLoadingWaitingRoom,
  } = useRequestToEndpoint(WAITING_ROOM_REQUEST)

  return (
    <TabContext value={value}>
      <TabList
        onChange={handleChange}
        aria-label="FluxNinja Scenarios"
        variant="fullWidth"
      >
        <Tab label="Waiting Room" value="0" />
        <Tab label="Rate Limit" value="1" />
        <Tab label="Workload Prioritization" value="2" />
      </TabList>

      <TabPanel value="0">
        <RequestMonitorPanel
          monitorRequestProps={{
            requestRecord: waitingRoomRequestRecord,
            userType: 'Waiting Room',
            refetch: refetchWaitingRoom,
          }}
          isLoading={isLoadingWaitingRoom}
        />
      </TabPanel>

      <TabPanel value="1">
        <RequestMonitorPanel
          monitorRequestProps={{
            requestRecord,
            userType: 'Rate Limit',
            refetch,
          }}
          isLoading={isLoadingRequest}
        />
      </TabPanel>
      <TabPanel value="2">
        <RequestMonitorPanel
          monitorRequestProps={{
            requestRecord: userRequestRecord,
            refetch: refetchUser,
            userType: 'Subscriber',
          }}
          isLoading={isLoadingUser}
        />
        <RequestMonitorPanel
          monitorRequestProps={{
            requestRecord: guestRequestRecord,
            refetch: refetchGuest,
            userType: 'Guest',
          }}
          isLoading={isLoadingGuest}
        />
      </TabPanel>
    </TabContext>
  )
}

export interface RequestMonitorPanelProps {
  monitorRequestProps: MonitorRequestProps
  isLoading: boolean
}

export const RequestMonitorPanel: FC<RequestMonitorPanelProps> = ({
  monitorRequestProps,
  isLoading,
}) => {
  const { errorComponent, isError } = useMemo(
    () =>
      monitorRequestProps.requestRecord?.[
        monitorRequestProps.requestRecord?.length - 1
      ] || { errorComponent: null, isError: false },
    [monitorRequestProps.requestRecord]
  )

  return (
    <HomePageWrapper>
      <HomePageColumnBox>
        <MonitorRequest {...monitorRequestProps} />
      </HomePageColumnBox>
      <HomePageColumnBox>
        {isLoading && !errorComponent && (
          <FlexBox>
            <CircularProgress />
          </FlexBox>
        )}
        {errorComponent}
        {!isLoading && !isError && (
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
