<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { api, setToken } from '../lib/api'

const router = useRouter()
const loading = ref(false)
const err = ref('')
const form = ref({ username: '', password: '' })

async function submit() {
  err.value = ''
  loading.value = true
  try {
    const data = await api.login(form.value)
    setToken(data.token)
    router.push('/app')
  } catch (e) {
    err.value = e.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <section class="auth-page">
    <div class="auth-card">
      <h1>FinalAI</h1>
      <p>登录你的 AI 工作台</p>

      <form @submit.prevent="submit" class="auth-form">
        <label>
          用户名
          <input v-model="form.username" required />
        </label>
        <label>
          密码
          <input v-model="form.password" type="password" required />
        </label>

        <button :disabled="loading">{{ loading ? '登录中...' : '登录' }}</button>
        <p class="error" v-if="err">{{ err }}</p>
      </form>

      <router-link to="/register" class="switch-link">没有账号？去注册</router-link>
    </div>
  </section>
</template>
