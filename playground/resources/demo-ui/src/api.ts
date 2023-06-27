import axios from 'axios'

export const API_URL = '/api'

export const api = axios.create({
  baseURL: API_URL,
})

export interface RequestSpec {
  method: 'GET' | 'POST' | 'PUT' | 'DELETE'
  endpoint: string
  userType: 'Crawler' | 'Subscriber' | 'Guest'
  userId: string
}
