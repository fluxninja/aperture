import React, { FC, useEffect, useRef, useState } from 'react'
import {
  MonitorRequest,
  MonitorRequestProps,
  RequestRecord,
} from '../components/monitor-request'
import { Button, Box, Typography, styled } from '@mui/material'
import {
  GracefulError,
  useGraceful,
  useGracefulRequest,
} from '@fluxninja-tools/graceful-js'
import { API_URL } from '../api'
import axios from 'axios'
// import { kApi } from '../k-api'

export const HomePage: FC = () => {
  const { refetch, isError, isErroredType } =
    useGracefulRequestForRateLimit('Crawler')
  const {
    refetch: refetchSubscriber,
    isError: isErrorSubscriber,
    isErroredType: isErroredTypeSubscriber,
  } = useGracefulRequestForRateLimit('Subscriber')

  return (
    <>
      <RequestMonitorPanel
        monitorRequestProps={{
          requestRecord: isErroredType,
          userType: 'Crawler',
        }}
        isErrored={isError}
        refetch={refetch}
      />
      <RequestMonitorPanel
        monitorRequestProps={{
          requestRecord: isErroredTypeSubscriber,
          userType: 'Subscriber',
        }}
        isErrored={isErrorSubscriber}
        refetch={refetchSubscriber}
      />
    </>
  )
}

export const useGracefulRequestForRateLimit = (userID: string) => {
  const [isErroredType, setIsErroredType] = useState<RequestRecord[]>([])
  const isErrorRef = useRef(false)

  const { isError, refetch } = useGracefulRequest<'Axios'>({
    typeOfRequest: 'Axios',
    requestFnc: () =>
      axios.post(`${API_URL}/rate-limit`, {
        headers: { 'Content-Type': 'application/json', 'user-id': userID },
        data: {},
      }),
    options: {
      disabled: true,
    },
  })

  isErrorRef.current = isError

  useEffect(() => {
    const requests = async () => {
      for (let i = 0; i < 100; i++) {
        await new Promise((resolve) => setTimeout(resolve, 800))

        refetch()

        if (isErrorRef.current) {
          setIsErroredType((prevErrors) => [
            ...prevErrors,
            { isError: true, rateLimitInfo: null },
          ])
          break
        }
      }
    }

    requests()
  }, [refetch])

  return { isError, refetch, isErroredType }
}

export interface RequestMonitorPanelProps {
  monitorRequestProps: MonitorRequestProps
  isErrored: boolean
  refetch?: () => void
}

export const RequestMonitorPanel: FC<RequestMonitorPanelProps> = ({
  monitorRequestProps,
  isErrored,
  refetch,
}) => (
  <HomePageWrapper>
    <HomePageColumnBox>
      <MonitorRequest {...monitorRequestProps} />
      <Button variant="contained" color="primary" onClick={refetch}>
        Start Refetch
      </Button>
    </HomePageColumnBox>
    <HomePageColumnBox>
      {isErrored ? (
        <GracefulError url="http://localhost:8099/api/rate-limit" />
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
