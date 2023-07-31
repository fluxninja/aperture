import { SyntheticEvent, useEffect, useMemo, useRef, useState } from 'react'
import { useSearchParams } from 'react-router-dom'

export const useTabChange = (tabs: string[]) => {
  const [value, setValue] = useState(0)
  const [search, setSearchParams] = useSearchParams()
  const queryParams = useMemo(() => new URLSearchParams(search), [search])
  const tab = useMemo(() => queryParams.get('tab'), [queryParams])
  const tabsRef = useRef(tabs)

  useEffect(() => {
    if (!tab) {
      setSearchParams({
        tab: tabsRef.current[0],
      })
      return
    }
    const index = tabsRef.current.indexOf(tab)
    const currentTab = index !== -1 ? index : 0
    setSearchParams({ tab: tabsRef.current[currentTab] })
    setValue(currentTab)
  }, [tab, setSearchParams, setValue])

  const handleChange = (_: SyntheticEvent, newValue: number) => {
    setValue(newValue)
    setSearchParams({ tab: tabs[newValue] })
  }

  return {
    handleChange,
    value: value.toString(),
  }
}
