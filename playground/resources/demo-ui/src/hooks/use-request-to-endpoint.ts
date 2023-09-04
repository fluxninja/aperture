import { useCallback, useEffect, useRef, useState } from 'react'
import { RequestSpec, api } from '../api'
import { RequestRecord } from '../components/monitor-request'
import { useGracefulRequest } from '@fluxninja-tools/graceful-js'

export const useRequestToEndpoint = (reqSpec: RequestSpec) => {
  const [requestRecord, setRequestRecord] = useState<RequestRecord[]>([]) // record state for each request
  const [requestCount, setRequestCount] = useState(0) // number of request count state
  const [intervalId, setIntervalId] = useState<NodeJS.Timeout | null>(null) // interval id state, used to clear interval

  const startFetchRef = useRef(false)
  const lastCounter = useRef(0)

  const { isError, refetch, error, data, isRetry, isLoading, errorComponent } =
    useGracefulRequest<'Axios'>({
      typeOfRequest: 'Axios',
      requestFnc: () => api(reqSpec),
      options: {
        disabled: true,
      },
    })

  /**
   * Push to the array of request record after every request we make by interval
   * Or if a retry happens after an error
   */
  useEffect(() => {
    if (
      !requestCount ||
      isLoading ||
      (lastCounter.current === requestCount && !isRetry)
    ) {
      return
    }

    lastCounter.current = requestCount
    setRequestRecord((prevErrors) => [
      ...prevErrors,
      {
        isError,
        rateLimitInfo: error?.rateLimitInfo || null,
        isRetry,
        errorComponent,
        error,
      },
    ])
  }, [error, errorComponent, isError, isLoading, isRetry, requestCount])

  /**
   * Start making request after every 800ms. We are using ref to make it sure that
   * we don't start multiple intervals.
   */
  const startFetch = useCallback(() => {
    if (startFetchRef.current) {
      return
    }

    const intervalId = setInterval(() => {
      setRequestCount((prevCount) => prevCount + 1)
      refetch()
    }, 800)

    setIntervalId(intervalId)

    startFetchRef.current = true

    return () => {
      clearInterval(intervalId)
    }
  }, [refetch])

  /**
   * We stop making request if we get error or request count is greater than 60
   */
  useEffect(() => {
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
  }, [intervalId, isError, requestCount])

  /**
   * After error we will stop the request, but if request get resolved after the graceful retry
   * We start requesting again until we get error or request count is greater than 60
   */
  useEffect(() => {
    if (isRetry && !isError && requestCount < 60 && startFetchRef.current) {
      clearInterval(intervalId)
      startFetchRef.current = false
      startFetch()
      return
    }
  }, [isError, isRetry])

  return { isError, refetch: startFetch, requestRecord, isLoading, data }
}
