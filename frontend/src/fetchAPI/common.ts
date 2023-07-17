import { getToken, resetCookie } from '../utils/cookie'
import router from 'next/router'

export const createUrl = (endpoint: string, pathParameters: Array<string | number>): string => {
  let url = 'http://localhost:8080/api' + endpoint
  if (pathParameters.length > 0) {
    url += '/' + pathParameters.map((p) => p.toString()).join('/')
  }
  return url
}

export const getFetcher = async (
  url: string,
  headerMap: Map<string, string> = new Map(),
): Promise<any> => {
  return await fetcher(url, 'GET', headerMap, new Map())
}

export const postFetcher = async (
  url: string,
  bodyMap: Map<string, number | string | boolean | number[] | string[]> = new Map(),
  headerMap: Map<string, string> = new Map(),
): Promise<any> => {
  return await fetcher(url, 'POST', headerMap, bodyMap)
}

export const patchFetcher = async (
  url: string,
  bodyMap: Map<string, number | string | boolean | number[] | string[]> = new Map(),
  headerMap: Map<string, string> = new Map(),
): Promise<any> => {
  return await fetcher(url, 'PATCH', headerMap, bodyMap)
}

export const deleteFetcher = async (
  url: string,
  bodyMap: Map<string, number | string | boolean | number[] | string[]> = new Map(),
  headerMap: Map<string, string> = new Map(),
): Promise<any> => {
  return await fetcher(url, 'DELETE', headerMap, bodyMap)
}

export const fetcher = async (
  url: string,
  method: string,
  headerMap: Map<string, string>,
  bodyMap: Map<string, number | string | boolean | number[] | string[]>,
): Promise<any> => {
  if (!(method == 'GET' || method == 'POST' || method == 'PATCH' || method == 'DELETE')) {
    throw new Error('fetcher error: method setting')
  }
  // Authorizationがない場合は設定する(AuthorizationはheadersMapに含めなくても自動的に設定される)
  if (!headerMap.has('Authorization')) {
    const jwtToken: string = getToken()
    if (jwtToken) {
      headerMap.set('Authorization', jwtToken)
    }
  }
  // 'Content-Type': 'application/json'も無ければ追加する
  if (!headerMap.has('Content-Type')) {
    headerMap.set('Content-Type', 'application/json')
  }

  let response: Response
  if (method == 'GET') {
    response = await fetch(url, {
      method: method,
      headers: Object.fromEntries(headerMap),
    })
  } else {
    response = await fetch(url, {
      method: method,
      headers: Object.fromEntries(headerMap),
      body: JSON.stringify(Object.fromEntries(bodyMap)),
    })
  }
  const res = await response.json()
  if (!response.ok) {
    if (response.status === 401) {
      resetCookie()
      router.replace('/')
      return
    }
    throw new Error(res)
  }
  return res
}
