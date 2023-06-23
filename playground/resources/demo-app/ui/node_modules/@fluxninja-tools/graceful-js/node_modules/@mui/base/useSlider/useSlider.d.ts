import { UseSliderParameters, UseSliderReturnValue } from './useSlider.types';
export declare function valueToPercent(value: number, min: number, max: number): number;
export declare const Identity: (x: any) => any;
/**
 *
 * Demos:
 *
 * - [Slider](https://mui.com/base/react-slider/#hook)
 *
 * API:
 *
 * - [useSlider API](https://mui.com/base/react-slider/hooks-api/#use-slider)
 */
export default function useSlider(parameters: UseSliderParameters): UseSliderReturnValue;
