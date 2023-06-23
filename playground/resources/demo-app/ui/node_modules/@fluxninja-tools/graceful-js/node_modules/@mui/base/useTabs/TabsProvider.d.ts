import * as React from 'react';
import { TabsContextValue } from '../Tabs/TabsContext';
import { CompoundComponentContextValue } from '../utils/useCompound';
export type TabsProviderValue = CompoundComponentContextValue<string | number, string> & TabsContextValue;
export interface TabsProviderProps {
    value: TabsProviderValue;
    children: React.ReactNode;
}
/**
 * Sets up the contexts for the underlying Tab and TabPanel components.
 *
 * @ignore - do not document.
 */
export default function TabsProvider(props: TabsProviderProps): JSX.Element;
