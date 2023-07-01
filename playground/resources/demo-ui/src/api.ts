import axios, { AxiosRequestConfig } from 'axios'

export const API_URL = '/aperture'

export const api = axios.create({
  baseURL: API_URL,
})

export declare type RequestSpec = AxiosRequestConfig
