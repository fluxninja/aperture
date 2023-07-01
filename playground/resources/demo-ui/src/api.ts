import axios from 'axios'

export const API_URL = '/api'

export const REQUEST_URL = '/request'

export const api = axios.create({
  baseURL: API_URL,
})

export const req = axios.create({
  baseURL: REQUEST_URL,
})

export interface RequestSpec {
  method: 'GET' | 'POST' | 'PUT' | 'DELETE'
  endpoint: string
  userType: 'subscriber' | 'guest' | 'crawler' | 'user'
  userId: string
  body?: any
}
