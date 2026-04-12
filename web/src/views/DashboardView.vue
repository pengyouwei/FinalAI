<script setup>
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { api, clearToken, getToken, streamPost } from '../lib/api'

const router = useRouter()
const profile = ref(null)
const error = ref('')

const modelOptions = [
  { value: '1', label: 'OpenAI' },
  { value: '2', label: 'Ollama' },
  { value: '3', label: 'RAG 知识库' },
]

const sessions = ref([])
const currentSessionId = ref('')
const modelType = ref('1')
const question = ref('')
const messages = ref([])

const imageFile = ref(null)
const imageResult = ref('')
const ragFile = ref(null)
const ragUploadResult = ref(null)
const streamBusy = ref(false)
const submitBusy = ref(false)
const uploadBusy = ref(false)
const deletingSessionId = ref('')
const createMode = ref(true)
const streamMode = ref(true)
const messageListRef = ref(null)
const textareaRef = ref(null)

const hasToken = computed(() => !!getToken())

function logout() {
  clearToken()
  router.push('/login')
}

function appendUserMessage(text) {
  messages.value.push({ role: 'user', content: text })
}

function appendAssistantMessage(text) {
  messages.value.push({ role: 'assistant', content: text })
}

function startNewChat() {
  currentSessionId.value = ''
  messages.value = []
  createMode.value = true
}

function adjustTextareaHeight() {
  const el = textareaRef.value
  if (el) {
    el.style.height = 'auto'
    el.style.height = `${el.scrollHeight}px`
  }
}

async function loadProfile() {
  profile.value = await api.profile()
}

async function loadSessions() {
  const data = await api.sessions()
  sessions.value = data.sessions || []
}

async function loadHistory() {
  if (!currentSessionId.value) {
    return
  }
  const data = await api.history({ sessionId: currentSessionId.value })
  messages.value = (data.history || []).map((item) => ({
    role: item.is_user ? 'user' : 'assistant',
    content: item.content,
  }))
}

async function createAndSend() {
  if (!question.value.trim()) {
    return
  }
  error.value = ''
  submitBusy.value = true
  try {
    const ask = question.value
    appendUserMessage(ask)
    question.value = ''
    adjustTextareaHeight()
    const data = await api.createChat({ question: ask, modelType: modelType.value })
    currentSessionId.value = data.sessionId
    appendAssistantMessage(data.information)
    await loadSessions()
  } catch (e) {
    error.value = e.message
  } finally {
    submitBusy.value = false
  }
}

async function sendInCurrent() {
  if (!question.value.trim() || !currentSessionId.value) {
    return
  }
  error.value = ''
  submitBusy.value = true
  try {
    const ask = question.value
    appendUserMessage(ask)
    question.value = ''
    adjustTextareaHeight()
    const data = await api.sendChat({
      sessionId: currentSessionId.value,
      question: ask,
      modelType: modelType.value,
    })
    appendAssistantMessage(data.information)
  } catch (e) {
    error.value = e.message
  } finally {
    submitBusy.value = false
  }
}

async function createAndStream() {
  if (!question.value.trim()) {
    return
  }
  error.value = ''
  streamBusy.value = true
  let currentChunk = ''

  try {
    const ask = question.value
    appendUserMessage(ask)
    question.value = ''
    adjustTextareaHeight()
    appendAssistantMessage('')

    await streamPost('/chat/create/stream', { question: ask, modelType: modelType.value }, (data) => {
      if (data === '[DONE]') {
        return
      }

      if (data.startsWith('{') && data.includes('sessionId')) {
        try {
          const parsed = JSON.parse(data)
          currentSessionId.value = parsed.sessionId || currentSessionId.value
        } catch {
          // ignore parse error
        }
        return
      }

      currentChunk += data
      messages.value[messages.value.length - 1].content = currentChunk
    })

    await loadSessions()
  } catch (e) {
    error.value = e.message
  } finally {
    streamBusy.value = false
  }
}

async function sendAndStream() {
  if (!question.value.trim() || !currentSessionId.value) {
    return
  }
  error.value = ''
  streamBusy.value = true
  let currentChunk = ''

  try {
    const ask = question.value
    appendUserMessage(ask)
    question.value = ''
    adjustTextareaHeight()
    appendAssistantMessage('')

    await streamPost('/chat/send/stream', {
      sessionId: currentSessionId.value,
      question: ask,
      modelType: modelType.value,
    }, (data) => {
      if (data === '[DONE]') {
        return
      }
      currentChunk += data
      messages.value[messages.value.length - 1].content = currentChunk
    })

  } catch (e) {
    error.value = e.message
  } finally {
    streamBusy.value = false
  }
}

async function submitQuestion() {
  if (!question.value.trim()) {
    return
  }

  if (!createMode.value && !currentSessionId.value) {
    error.value = '请先在左侧选择一个会话，或点击“+ 新建聊天”'
    return
  }

  if (createMode.value) {
    if (streamMode.value) {
      await createAndStream()
      return
    }
    await createAndSend()
    return
  }

  if (streamMode.value) {
    await sendAndStream()
    return
  }

  await sendInCurrent()
}

async function doImageRecognize() {
  if (!imageFile.value) {
    return
  }
  error.value = ''
  submitBusy.value = true
  try {
    const data = await api.recognizeImage(imageFile.value)
    imageResult.value = data.class_name
  } catch (e) {
    error.value = e.message
  } finally {
    submitBusy.value = false
  }
}

function onRagFileChange(e) {
  ragFile.value = e.target.files?.[0] || null
}

async function doRagUpload() {
  if (!ragFile.value) {
    error.value = '请先选择要上传的 RAG 文件'
    return
  }

  error.value = ''
  uploadBusy.value = true
  ragUploadResult.value = null
  try {
    const data = await api.uploadRagFile(ragFile.value)
    ragUploadResult.value = data
    modelType.value = '3'
  } catch (e) {
    error.value = e.message
  } finally {
    uploadBusy.value = false
  }
}

function chooseSession(id) {
  currentSessionId.value = id
  createMode.value = false
  loadHistory()
}

async function deleteSession(id) {
  if (!id || deletingSessionId.value) {
    return
  }

  const ok = window.confirm('确认删除该会话及其历史消息吗？该操作不可恢复。')
  if (!ok) {
    return
  }

  error.value = ''
  deletingSessionId.value = id
  try {
    await api.deleteSession({ sessionId: id })

    if (currentSessionId.value === id) {
      currentSessionId.value = ''
      messages.value = []
      createMode.value = true
    }

    await loadSessions()
  } catch (e) {
    error.value = e.message
  } finally {
    deletingSessionId.value = ''
  }
}

watch(messages, async () => {
  await nextTick()
  const el = messageListRef.value
  if (el) {
    el.scrollTop = el.scrollHeight
  }
}, { deep: true })

watch(question, () => {
  adjustTextareaHeight()
})

onMounted(async () => {
  if (!hasToken.value) {
    router.push('/login')
    return
  }

  try {
    await Promise.all([loadProfile(), loadSessions()])
  } catch (e) {
    error.value = e.message
  }
})
</script>

<template>
  <section class="gpt-layout">
    <aside class="gpt-sidebar">
      <div class="sidebar-header">
        <button class="sidebar-main-btn" @click="startNewChat">+ 新建聊天</button>
      </div>

      <div class="session-list-wrap">
        <div class="sidebar-title-row">
          <h3>会话历史</h3>
          <button class="mini-btn" @click="loadSessions" title="刷新会话列表">
            <svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" viewBox="0 0 24 24"><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 11A8.1 8.1 0 0 0 4.5 9M4 5v4h4m-4 4a8.1 8.1 0 0 0 15.5 2m.5 4v-4h-4"></path></svg>
          </button>
        </div>
        <ul class="session-list">
          <li
            v-for="s in sessions"
            :key="s.session_id"
            :class="{ active: s.session_id === currentSessionId && !createMode }"
            @click="chooseSession(s.session_id)"
          >
            <div class="session-item-head">
              <p>{{ s.title }}</p>
              <button
                class="session-delete-btn"
                :disabled="deletingSessionId === s.session_id"
                @click.stop="deleteSession(s.session_id)"
                title="删除该会话"
              >
                {{ deletingSessionId === s.session_id ? '...' : '删' }}
              </button>
            </div>
            <small>{{ s.session_id.slice(0, 10) }}...</small>
          </li>
        </ul>
      </div>
    </aside>

    <main class="gpt-main">
      <header class="gpt-header">
        <div class="header-left">
          <h1>FinalAI Chat</h1>
          <p class="session-info" v-if="!createMode && currentSessionId">
            当前会话: {{ currentSessionId.slice(0, 10) }}...
          </p>
          <p class="session-info" v-else>新会话</p>
        </div>
        <div class="mode-pills">
          <button :class="['pill', { active: createMode }]" @click="createMode = true">新会话</button>
          <button :class="['pill', { active: !createMode }]" @click="createMode = false" :disabled="!currentSessionId">当前会话</button>
          <span class="pill-divider"></span>
          <button :class="['pill', { active: streamMode }]" @click="streamMode = true">流式</button>
          <button :class="['pill', { active: !streamMode }]" @click="streamMode = false">普通</button>
          <select v-model="modelType" class="model-select">
            <option v-for="item in modelOptions" :key="item.value" :value="item.value">{{ item.label }}</option>
          </select>
        </div>
      </header>

      <section class="gpt-messages" ref="messageListRef">
        <div class="empty-hint" v-if="messages.length === 0">
          <div class="empty-logo">🤖</div>
          <p>今天我能帮你做些什么？</p>
        </div>
        <article v-for="(m, i) in messages" :key="i" :class="['message-row', m.role]">
          <div class="avatar">
            <span v-if="m.role === 'user'">U</span>
            <svg v-else xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" viewBox="0 0 256 256"><path fill="currentColor" d="M224 128a96 96 0 1 1-96-96a96 96 0 0 1 96 96m-32-48a12 12 0 1 0-12-12a12 12 0 0 0 12 12m-80 0a12 12 0 1 0-12-12a12 12 0 0 0 12 12m-2.2 62.2a12 12 0 0 0-19.6-12.4a48 48 0 0 1-72.4 0a12 12 0 1 0-19.6 12.4a72.06 72.06 0 0 0 111.6 0"></path></svg>
          </div>
          <div class="message-bubble">
            <p v-if="m.content">{{ m.content }}</p>
            <div class="loading-dots" v-else>
              <span></span>
              <span></span>
              <span></span>
            </div>
          </div>
        </article>
      </section>

      <section class="composer-wrap">
        <div class="composer">
          <textarea
            v-model="question"
            rows="1"
            placeholder="给 FinalAI 发送消息..."
            @keydown.enter.exact.prevent="submitQuestion"
            @input="adjustTextareaHeight"
            ref="textareaRef"
          />
          <button :disabled="submitBusy || streamBusy" @click="submitQuestion" title="发送消息">
            <svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" viewBox="0 0 24 24"><path fill="currentColor" d="M3 20V4l19 8z"></path></svg>
          </button>
        </div>
        <p class="error" v-if="error">{{ error }}</p>
      </section>
    </main>

    <aside class="gpt-toolside">
      <div class="toolside-scroll">
        <section class="vision-card">
          <h4>图片识别</h4>
          <input type="file" accept="image/*" @change="(e) => (imageFile = e.target.files?.[0] || null)" />
          <button class="mini-btn fill" :disabled="submitBusy" @click="doImageRecognize">识别</button>
          <p v-if="imageResult" class="vision-result">结果: {{ imageResult }}</p>
        </section>

        <section class="vision-card">
          <h4>RAG 文件上传</h4>
          <p class="tool-hint">上传成功后将自动切换到 RAG 知识库模型。</p>
          <input type="file" accept=".txt,.md,.markdown,text/plain,text/markdown" @change="onRagFileChange" />
          <button class="mini-btn fill full" :disabled="uploadBusy || submitBusy || streamBusy" @click="doRagUpload">
            {{ uploadBusy ? '上传中...' : '上传' }}
          </button>
          <p v-if="ragUploadResult?.file_name" class="vision-result">文件: {{ ragUploadResult.file_name }}</p>
          <p v-if="ragUploadResult?.file_path" class="vision-result">路径: {{ ragUploadResult.file_path }}</p>
          <p v-if="ragUploadResult" class="vision-success">知识库已就绪，可直接提问。</p>
        </section>
      </div>

      <footer class="sidebar-footer toolside-footer">
        <span v-if="profile" class="username">@{{ profile.username }}</span>
        <button class="mini-btn danger" @click="logout">退出</button>
      </footer>
    </aside>
  </section>
</template>
