import React, { FC, useEffect, useMemo, useState } from 'react'
import {
  MonitorRequest,
  MonitorRequestProps,
  RequestRecord,
} from '../components/monitor-request'
import { Box, Typography, styled } from '@mui/material'
import { GracefulError, useGraceful } from '@fluxninja-tools/graceful-js'
import { API_URL } from '../api'
import axios from 'axios'
// import { kApi } from '../k-api'

export const HomePage: FC = () => {
  const { collection, gracefulError } = useGracefulRequestForRateLimit()
  return (
    <RequestMonitorPanel
      monitorRequestProps={{
        requestRecord: collection,
        userType: 'Subscriber',
      }}
      isErrored={collection[collection.length - 1]?.isError && gracefulError}
    />
  )
}

export const useGracefulRequestForRateLimit = () => {
  const [requestCollection, setRequestCollection] = useState<RequestRecord[]>(
    []
  )

  const { ctx } = useGraceful()

  const gracefulError = useMemo(
    () => ctx.url === `${API_URL}/request` && ctx.isError,
    [ctx]
  )

  useEffect(() => {
    const requests = async () => {
      for (let i = 0; i < 100; i++) {
        await new Promise((resolve) => setTimeout(resolve, 800))
        try {
          // await kApi()
          const res = await axios.post(`${API_URL}/rate-limit`, {
            headers: { 'Content-Type': 'application/json', 'user-id': 'foo' },
            data: {},
          })

          if (res?.status === 200) {
            setRequestCollection((prev) => [
              ...prev,
              { isError: false, rateLimitInfo: null },
            ])
          }
        } catch (err) {
          setRequestCollection((prev) => [
            ...prev,
            { isError: true, rateLimitInfo: null },
          ])
          break
        }
      }
    }

    requests()
  }, [setRequestCollection])

  return { collection: requestCollection, gracefulError }
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
        <GracefulError url="http://localhost:8099/request" />
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
