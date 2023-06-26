import axios from 'axios'

export const API_URL = '/api' //'http://localhost:9099'

export const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
    'User-Id': 'foo',
    Cookie:
      'session=eyJ1c2VyIjoia2Vub2JpIn0.YbsY4Q.kTaKRTyOIfVlIbNB48d9YH6Q0wo',
  },
})
