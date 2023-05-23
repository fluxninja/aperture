import {
  Box,
  BoxProps,
  Button,
  Divider,
  Paper,
  Typography,
  styled,
} from '@mui/material'
import React, { useEffect, useState } from 'react'
import { FC } from 'react'

export declare type RequestRecord = {
  isError: boolean
}

export interface MonitorRequestProps {
  requestRecord: RequestRecord[]
  userType: 'Subscriber' | 'Guest'
}

export const MonitorRequest: FC<MonitorRequestProps> = ({
  requestRecord,
  userType,
}) => (
  <MonitorRequestWrapper component={Paper}>
    <Typography variant="h6" textAlign="center">
      {userType}
    </Typography>
    <Divider />
    <Typography>Request made in 60s:</Typography>
    <Box display="flex" flexDirection="row" gap={0.3} justifyContent="center">
      {requestRecord.map((record, index) => (
        <MonitorRequestItem key={index} {...record} />
      ))}
    </Box>
    <Box display="grid" gridTemplateColumns="1fr 1fr" gap={1}>
      <Box {...columnBoxProps}>
        <Typography>Success rate:</Typography>
        <SuccessRate requestRecord={requestRecord} />
      </Box>
      <Box {...columnBoxProps}>
        <Typography>Retry After:</Typography>
        <Button variant="contained" color="secondary">
          Retry
        </Button>
      </Box>
    </Box>
  </MonitorRequestWrapper>
)

export type SuccessRateProps = Pick<MonitorRequestProps, 'requestRecord'>

export const SuccessRate: FC<SuccessRateProps> = ({ requestRecord }) => {
  const [successRate, setSuccessRate] = useState(0)

  useEffect(() => {
    const successCount = requestRecord.filter((record) => !record.isError)
    const percent = (successCount.length / requestRecord.length) * 100
    setSuccessRate(percent)
  }, [requestRecord, setSuccessRate])

  return <SuccessRateStyled>{`${successRate}%`}</SuccessRateStyled>
}

export const columnBoxProps: BoxProps = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  gap: 2,
}

export const MonitorRequestWrapper = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(2),
  padding: theme.spacing(2),
  maxWidth: 500,
}))

export const SuccessRateStyled = styled(Box)(({ theme }) => ({
  background: theme.palette.primary.main,
  color: theme.palette.common.white,
  fontSize: 24,
  fontWeight: '600',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  borderRadius: '50%',
  height: 150,
  width: 150,
}))

export const MonitorRequestItem = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'isError',
})<RequestRecord>(({ theme, isError }) => ({
  height: 50,
  width: 5,
  backgroundColor: isError
    ? theme.palette.error.main
    : theme.palette.success.main,
  borderRadius: theme.spacing(1),
}))
