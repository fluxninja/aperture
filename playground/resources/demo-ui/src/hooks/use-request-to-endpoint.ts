import {useGracefulRequest} from '@fluxninja-tools/graceful-js'
import {useCallback, useEffect, useState} from 'react'
import {api, RequestSpec} from '../api'
import {RequestRecord} from '../components/monitor-request'

export const useRequestToEndpoint = (reqSpec: RequestSpec) => {
  const [requestRecord, setRequestRecord] = useState<RequestRecord[]>([]) // record state for each request
  const [requestCount, setRequestCount] = useState(0) // number of request count state
  const [intervalId, setIntervalId] = useState<NodeJS.Timeout | null>(null) // interval id state, used to clear interval

  const { isError, refetch, error, data, isRetry, isLoading } =
    useGracefulRequest<'Axios'>({
      typeOfRequest: 'Axios',
      requestFnc: () => api(reqSpec),
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
    }, 500)

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
    if (requestCount >= 150) {
      clearInterval(intervalId)
      return
    }
  }, [requestCount, intervalId, isError])

  return { isError, refetch: startFetch, requestRecord, isLoading, data, error }
}
