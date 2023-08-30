import { RateLimitInfo } from '@fluxninja-tools/graceful-js'
import {
  Box,
  BoxProps,
  Button,
  Divider,
  Paper,
  Tooltip,
  Typography,
  styled,
} from '@mui/material'
import { AxiosError } from 'axios'
import React, { useEffect, useState } from 'react'
import { FC } from 'react'

export declare type RequestRecord = {
  isError: boolean
  errorComponent: JSX.Element | null
  error: AxiosError & {
    rateLimitInfo?: RateLimitInfo
  }
  rateLimitInfo?: RateLimitInfo
  isRetry?: boolean
}

export interface MonitorRequestProps {
  requestRecord: RequestRecord[]
  userType: string
  refetch: () => void
}

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
    <RequestMessagingInfo />
    <Typography>Request Info:</Typography>
    <Box display="flex" gap={0.3}>
      {requestRecord.map((record, index) => (
        <Tooltip
          key={index}
          title={record.error?.response?.status || 200}
          placement="top"
          arrow
        >
          <MonitorRequestItem {...record} />
        </Tooltip>
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
            disabled={requestRecord.length < 60 && requestRecord.length > 0}
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

export const RequestMessagingInfo: FC = () => (
  <RequestMessagingInfoWrapper>
    <RowBox>
      <RequestMessageInfoIndicator indicator="success" />
      <Typography>Success</Typography>
    </RowBox>
    <RowBox>
      <RequestMessageInfoIndicator indicator="error" />
      <Typography>Error</Typography>
    </RowBox>
    <RowBox>
      <RequestMessageInfoIndicator indicator="warning" />
      <Typography>Retry</Typography>
    </RowBox>
  </RequestMessagingInfoWrapper>
)

const boxFlex: BoxProps = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: 100,
}

export const RequestMessagingInfoWrapper = styled(Box)(({ theme }) => ({
  display: 'grid',
  gridTemplateColumns: '1fr 1fr 1fr',
  gap: theme.spacing(1),
  [theme.breakpoints.down('md')]: {
    display: 'flex',
    flexDirection: 'column',
    gap: theme.spacing(1),
  },
}))

export const RequestMessageInfoIndicator = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'indicator',
})<{
  indicator: 'success' | 'error' | 'warning'
}>(({ theme, indicator }) => ({
  backgroundColor: theme.palette[indicator].main,
  borderRadius: theme.spacing(1),
  height: 5,
  width: 15,
}))

export const RowBox = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: theme.spacing(1),
}))

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
  cursor: 'pointer',
  backgroundColor: isError
    ? theme.palette.error.main
    : isRetry
    ? theme.palette.warning.main
    : theme.palette.success.main,
  borderRadius: theme.spacing(1),
}))
