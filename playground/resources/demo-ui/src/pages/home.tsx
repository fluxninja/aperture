import React, { FC, useCallback, useEffect, useState } from 'react'
import {
  MonitorRequest,
  MonitorRequestProps,
  RequestRecord,
} from '../components/monitor-request'
import { Box, Typography, styled, Tabs, Tab, Backdrop, Fade } from '@mui/material'
import { Dialog, DialogTitle, DialogContent, DialogContentText, DialogActions, Button } from '@mui/material'

import { TabContext, TabList, TabPanel } from '@mui/lab'
import {
  GracefulError,
  GracefulErrorProps,
  useGracefulRequest,
} from '@fluxninja-tools/graceful-js'
import { api, RequestSpec } from '../api'
import { SuccessIcon } from './success-icon'

export const HomePage: FC = () => {

  const [value, setValue] = useState('1')
  const [userType, setUserType] = useState<'guest' | 'subscriber'>('guest')
  const [open, setOpen] = useState(false);

  const handleChange = (event: React.SyntheticEvent, newValue: string) => {
    setValue(newValue)
    if(newValue === '2') setOpen(true);
  }

  const handleDialogClose = (value: 'guest' | 'subscriber') => {
    setUserType(value);
    setOpen(false);
  };
  
  // Request Spec for rate-limit endpoint with Executive user to retrieve fluxninja founders info
  const reqSpec: RequestSpec = {
    method: 'GET',
    endpoint: '/rate-limit',
    userType: 'user',
    userId: 'DemoUI',
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
    method: 'GET',
    endpoint: '/workload-prioritization',
    userType: 'subscriber',
    userId: 'DemoUI',
  }
  const {
    refetch: refetchUser,
    isError: isErrorUser,
    requestRecord: userRequestRecord,
    isLoading: isLoadingUser,
    data: userRequestResponse,
  } = useRequestToEndpoint(reqSpec2)

  return (
    <TabContext value={value}>
      <TabList
        onChange={handleChange}
        aria-label="FluxNinja Scenarios"
        variant="fullWidth"
      >
        <Tab label="Rate Limit" value="1" />
        <Tab label="Workload Prioritization" value="2" />
      </TabList>
      <TabPanel value="1">
        <RequestMonitorPanel
          monitorRequestProps={{
            requestRecord: requestRecord,
            userType: 'User',
            refetch,
          }}
          isErrored={isError}
          isLoading={isLoadingRequest}
          errorComponentProps={{
            url: `/api${reqSpec.endpoint}`,
            requestBody: {},
          }}
          responseData={requestResponse}
        />
      </TabPanel>
      <TabPanel value="2">
        <Backdrop
          sx={{ color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1 }}
          open={open}
        />
          <Dialog
            open={open}
            TransitionComponent={Fade}
            onClose={() => setOpen(false)}
            aria-labelledby="alert-dialog-title"
            aria-describedby="alert-dialog-description"
          >
          <DialogTitle id="alert-dialog-title">{"User Type"}</DialogTitle>
          <DialogContent>
            <DialogContentText id="alert-dialog-description">
              Please select the user type.
            </DialogContentText>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => handleDialogClose('guest')} color="primary">
              Guest
            </Button>
            <Button onClick={() => handleDialogClose('subscriber')} color="primary" autoFocus>
              Subscriber
            </Button>
          </DialogActions>
        </Dialog>
        <RequestMonitorPanel
          monitorRequestProps={{
            requestRecord: userRequestRecord,
            refetch: refetchUser,
            userType: userType === 'subscriber' ? 'Subscriber' : 'Guest',
          }}
          isErrored={isErrorUser}
          isLoading={isLoadingUser}
          errorComponentProps={{
            url: `/api${reqSpec2.endpoint}`,
            requestBody: {},
          }}
          responseData={userRequestResponse}
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
    }, 100)

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
  }, [requestCount, intervalId, isError, reqSpec.userType])

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
    <SuccessIcon />
    <HomePageColumnBox>
      {isErrored && !isLoading ? (
        <GracefulError {...errorComponentProps} />
      ) : (
        <Box
          sx={(theme) => ({
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'center',
            alignItems: 'center',
            color: theme.palette.grey[400],
          })}
        >
          <Typography variant="h5">{responseData?.data?.message}</Typography>
          {responseData?.status === 200 && <SuccessIcon />}
        </Box>
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
