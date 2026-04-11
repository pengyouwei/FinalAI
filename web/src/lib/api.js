const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

export function setToken(token) {
  localStorage.setItem('finalai_token', token)
}

export function getToken() {
  return localStorage.getItem('finalai_token') || ''
}

export function clearToken() {
  localStorage.removeItem('finalai_token')
}

function buildHeaders(extra = {}, needAuth = false) {
  const headers = { ...extra }
  if (needAuth) {
    const token = getToken()
    if (token) {
      headers.Authorization = `Bearer ${token}`
    }
  }
  return headers
}

async function request(path, options = {}) {
  const res = await fetch(`${API_BASE}${path}`, options)
  const data = await res.json().catch(() => ({}))
  if (!res.ok || data.code !== 0) {
    const msg = data.msg || `Request failed (${res.status})`
    throw new Error(msg)
  }
  return data.data
}

export const api = {
  register(payload) {
    return request('/user/register', {
      method: 'POST',
      headers: buildHeaders({ 'Content-Type': 'application/json' }),
      body: JSON.stringify(payload),
    })
  },

  login(payload) {
    return request('/user/login', {
      method: 'POST',
      headers: buildHeaders({ 'Content-Type': 'application/json' }),
      body: JSON.stringify(payload),
    })
  },

  profile() {
    return request('/user/profile', {
      method: 'GET',
      headers: buildHeaders({}, true),
    })
  },

  sessions() {
    return request('/chat/sessions', {
      method: 'GET',
      headers: buildHeaders({}, true),
    })
  },

  createChat(payload) {
    return request('/chat/create', {
      method: 'POST',
      headers: buildHeaders({ 'Content-Type': 'application/json' }, true),
      body: JSON.stringify(payload),
    })
  },

  sendChat(payload) {
    return request('/chat/send', {
      method: 'POST',
      headers: buildHeaders({ 'Content-Type': 'application/json' }, true),
      body: JSON.stringify(payload),
    })
  },

  history(payload) {
    return request('/chat/history', {
      method: 'POST',
      headers: buildHeaders({ 'Content-Type': 'application/json' }, true),
      body: JSON.stringify(payload),
    })
  },

  recognizeImage(file) {
    const form = new FormData()
    form.append('image', file)
    return request('/image/recognize', {
      method: 'POST',
      headers: buildHeaders({}, true),
      body: form,
    })
  },
}

export async function streamPost(path, payload, onData) {
  const res = await fetch(`${API_BASE}${path}`, {
    method: 'POST',
    headers: buildHeaders({ 'Content-Type': 'application/json' }, true),
    body: JSON.stringify(payload),
  })

  if (!res.ok || !res.body) {
    throw new Error(`Stream request failed (${res.status})`)
  }

  const reader = res.body.getReader()
  const decoder = new TextDecoder('utf-8')
  let buffer = ''

  while (true) {
    const { done, value } = await reader.read()
    if (done) {
      break
    }

    buffer += decoder.decode(value, { stream: true })
    const parts = buffer.split('\n\n')
    buffer = parts.pop() || ''

    for (const block of parts) {
      const line = block
        .split('\n')
        .find((item) => item.startsWith('data:'))
      if (!line) {
        continue
      }
      const data = line.replace(/^data:\s?/, '')
      onData(data)
    }
  }
}
