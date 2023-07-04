import { RateLimitInfo } from '@fluxninja-tools/graceful-js'
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
  rateLimitInfo?: RateLimitInfo
  isRetry?: boolean
}

export interface MonitorRequestProps {
  requestRecord: RequestRecord[]
  userType: 'Subscriber' | 'Guest' | 'Crawler' | 'User' | 'Rate Limit'
  refetch: () => void
}

// fix where is showing NaN% when there is no request
export const MonitorRequest: FC<MonitorRequestProps> = ({
  requestRecord,
  userType,
  refetch,
}) => (
  <MonitorRequestWrapper component={Paper}>
    <Typography variant="h6" textAlign="center">
      {userType}
    </Typography>
    <Divider />
    <Typography>Request Info:</Typography>
    <Box display="flex" gap={0.3}>
      {requestRecord.map((record, index) => (
        <MonitorRequestItem key={index} {...record} />
      ))}
    </Box>
    <Box display="grid" gridTemplateColumns="1fr 1fr" gap={1}>
      <ColumnBoxStyled component={Paper}>
        <InfoHeading>Success</InfoHeading>
        <Box {...boxFlex}>
          <SuccessRate requestRecord={requestRecord} />
        </Box>
      </ColumnBoxStyled>
      <ColumnBoxStyled component={Paper}>
        <InfoHeading>Error</InfoHeading>
        <Box {...boxFlex}>
          <ErrorRate requestRecord={requestRecord} />
        </Box>
      </ColumnBoxStyled>
      <ColumnBoxStyled component={Paper}>
        <InfoHeading>Start Request</InfoHeading>
        <Box {...boxFlex}>
          <Button
            variant="contained"
            color="primary"
            sx={{ alignSelf: 'center' }}
            onClick={refetch}
          >
            Start
          </Button>
        </Box>
      </ColumnBoxStyled>
      <ColumnBoxStyled component={Paper}>
        <InfoHeading>Reset</InfoHeading>
        <Box {...boxFlex}>
          <Button
            variant="contained"
            color="secondary"
            sx={{ alignSelf: 'center' }}
            onClick={() => window.location.reload()}
          >
            Reset
          </Button>
        </Box>
      </ColumnBoxStyled>
    </Box>
  </MonitorRequestWrapper>
)

export type SuccessRateProps = Pick<MonitorRequestProps, 'requestRecord'>

export const findPercentage = (share: number, total: number) =>
  Math.round(((share / total) * 100 + Number.EPSILON) * 100) / 100

export const useSuccessErrorRatePercent = (requestRecord: RequestRecord[]) => {
  const [successRate, setSuccessRate] = useState(0)
  const [errorRate, setErrorRate] = useState(0)

  useEffect(() => {
    const successCount = requestRecord.filter((record) => !record.isError)
    const percent = findPercentage(successCount.length, requestRecord.length)
    setSuccessRate(percent)
  }, [requestRecord, setSuccessRate])

  useEffect(() => {
    const errorCount = requestRecord.filter((record) => record.isError)
    const percent = findPercentage(errorCount.length, requestRecord.length)
    setErrorRate(percent)
  }, [requestRecord, setErrorRate])

  return { successRate, errorRate }
}

export const SuccessRate: FC<SuccessRateProps> = ({ requestRecord }) => {
  const { successRate } = useSuccessErrorRatePercent(requestRecord)
  return (
    <SuccessErrorRateStyled isError={false}>
      {successRate || 0}%
    </SuccessErrorRateStyled>
  )
}

export const ErrorRate: FC<SuccessRateProps> = ({ requestRecord }) => {
  const { errorRate } = useSuccessErrorRatePercent(requestRecord)
  return (
    <SuccessErrorRateStyled isError={true}>
      {errorRate || 0}%
    </SuccessErrorRateStyled>
  )
}

const boxFlex: BoxProps = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: 100,
}

export const ColumnBoxStyled = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  gap: 2,
  border: `0.5px solid ${theme.palette.grey[300]}`,
  paddingBottom: theme.spacing(1),
  minHeight: 100,
}))

export const InfoHeading = styled(Typography)(({ theme }) => ({
  borderBottom: `0.5px solid ${theme.palette.grey[300]}`,
  width: '100%',
  textAlign: 'center',
}))

export const MonitorRequestWrapper = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(2),
  padding: theme.spacing(2),
  maxWidth: 500,
}))

export const SuccessErrorRateStyled = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'isError',
})<{
  isError: boolean
}>(({ theme, isError }) => ({
  background: isError ? theme.palette.error.main : theme.palette.success.main,
  color: theme.palette.common.white,
  fontSize: 18,
  fontWeight: '600',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  borderRadius: '50%',
  height: 75,
  width: 75,
  alignSelf: 'center',
}))

export const MonitorRequestItem = styled(Box, {
  shouldForwardProp: (prop) =>
    prop !== 'isError' && prop !== 'rateLimitInfo' && prop !== 'isRetry',
})<RequestRecord>(({ theme, isError, isRetry }) => ({
  height: 50,
  width: 5,
  backgroundColor: isError
    ? theme.palette.error.main
    : isRetry
    ? theme.palette.warning.main
    : theme.palette.success.main,
  borderRadius: theme.spacing(1),
}))
