import { SyntheticEvent, useEffect, useMemo, useState } from 'react'
import { useSearchParams } from 'react-router-dom'

export const useTabChange = (tabs: string[]) => {
  const [value, setValue] = useState(0)
  const [search, setSearchParams] = useSearchParams()
  const queryParams = new URLSearchParams(search)
  const tab = useMemo(() => queryParams.get('tab'), [queryParams])

  useEffect(() => {
    if (!tab) {
      setSearchParams({
        tab: tabs[0],
      })
      return
    }
    const index = tabs.indexOf(tab)
    const currentTab = index !== -1 ? index : 0
    setSearchParams({ tab: tabs[currentTab] })
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
