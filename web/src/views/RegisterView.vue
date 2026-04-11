<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../lib/api'

const router = useRouter()
const loading = ref(false)
const err = ref('')
const form = ref({ username: '', password: '', confirm_password: '' })

async function submit() {
  err.value = ''
  loading.value = true
  try {
    await api.register(form.value)
    router.push('/login')
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
      <h1>创建账号</h1>
      <p>注册后即可使用聊天与图像识别能力</p>

      <form @submit.prevent="submit" class="auth-form">
        <label>
          用户名
          <input v-model="form.username" required />
        </label>
        <label>
          密码
          <input v-model="form.password" type="password" required />
        </label>
        <label>
          确认密码
          <input v-model="form.confirm_password" type="password" required />
        </label>

        <button :disabled="loading">{{ loading ? '注册中...' : '注册' }}</button>
        <p class="error" v-if="err">{{ err }}</p>
      </form>

      <router-link to="/login" class="switch-link">已有账号？去登录</router-link>
    </div>
  </section>
</template>
