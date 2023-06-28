import axios from 'axios'

export const API_URL = '/api'

export const api = axios.create({
  baseURL: API_URL,
})

export interface RequestSpec {
  method: 'GET' | 'POST' | 'PUT' | 'DELETE'
  endpoint: string
  userType: 'subscriber' | 'guest' | 'crawler' | 'executive'
  userId: string
  body?: any,
}
