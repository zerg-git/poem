<template>
  <div class="profile-container">
    <div class="profile-card">
      <div class="profile-header">
        <div class="avatar">
          <img v-if="user.avatar_url" :src="user.avatar_url" :alt="user.nickname" />
          <div v-else class="avatar-placeholder">{{ user.nickname?.charAt(0) || '?' }}</div>
        </div>
        <h2>{{ user.nickname || user.username }}</h2>
        <p class="username">@{{ user.username }}</p>
      </div>

      <div class="profile-stats">
        <div class="stat-item">
          <span class="stat-label">等级</span>
          <span class="stat-value">{{ user.level }}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">经验</span>
          <span class="stat-value">{{ user.experience }}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">金币</span>
          <span class="stat-value">{{ user.coins }}</span>
        </div>
      </div>

      <!-- 可编辑的信息区域 -->
      <div class="profile-info">
        <div class="info-item editable" @click="editEmail">
          <label>邮箱</label>
          <div class="info-value">
            <span>{{ user.email || '未设置' }}</span>
            <span class="edit-icon">✏️</span>
          </div>
        </div>
        <div class="info-item editable" @click="editPhone">
          <label>手机</label>
          <div class="info-value">
            <span>{{ user.phone || '未设置' }}</span>
            <span class="edit-icon">✏️</span>
          </div>
        </div>
        <div class="info-item">
          <label>注册时间</label>
          <span>{{ user.created_at }}</span>
        </div>
      </div>

      <button @click="handleLogout" class="logout-btn">退出登录</button>
    </div>

    <!-- 编辑邮箱弹窗 -->
    <div v-if="showEmailModal" class="modal-overlay" @click.self="closeModals">
      <div class="modal">
        <div class="modal-header">
          <h3>编辑邮箱</h3>
          <button @click="closeModals" class="close-btn">×</button>
        </div>
        <form @submit.prevent="saveEmail" class="modal-body">
          <div class="form-group">
            <label>新邮箱</label>
            <input
              v-model="emailForm.email"
              type="email"
              placeholder="请输入邮箱"
              required
            />
          </div>
          <div v-if="emailError" class="error-message">{{ emailError }}</div>
          <div class="modal-footer">
            <button type="button" @click="closeModals" class="cancel-btn">取消</button>
            <button type="submit" class="save-btn" :disabled="saving">
              {{ saving ? '保存中...' : '保存' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 编辑手机弹窗 -->
    <div v-if="showPhoneModal" class="modal-overlay" @click.self="closeModals">
      <div class="modal">
        <div class="modal-header">
          <h3>编辑手机号</h3>
          <button @click="closeModals" class="close-btn">×</button>
        </div>
        <form @submit.prevent="savePhone" class="modal-body">
          <div class="form-group">
            <label>新手机号</label>
            <input
              v-model="phoneForm.phone"
              type="tel"
              placeholder="请输入11位手机号"
              maxlength="11"
              pattern="[0-9]{11}"
              required
            />
          </div>
          <div v-if="phoneError" class="error-message">{{ phoneError }}</div>
          <div class="modal-footer">
            <button type="button" @click="closeModals" class="cancel-btn">取消</button>
            <button type="submit" class="save-btn" :disabled="saving">
              {{ saving ? '保存中...' : '保存' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { userApi } from '@/api/user-api'

const router = useRouter()
const userStore = useUserStore()

const user = ref({
  username: '',
  nickname: '',
  avatar_url: '',
  email: '',
  phone: '',
  level: 1,
  experience: 0,
  coins: 0,
  created_at: ''
})

const showEmailModal = ref(false)
const showPhoneModal = ref(false)
const saving = ref(false)
const emailError = ref('')
const phoneError = ref('')

const emailForm = ref({
  email: ''
})

const phoneForm = ref({
  phone: ''
})

onMounted(async () => {
  await loadProfile()
})

const loadProfile = async () => {
  try {
    const response = await userApi.getProfile()
    if (response.success) {
      user.value = response.data
    }
  } catch (error) {
    console.error('加载用户信息失败:', error)
  }
}

const editEmail = () => {
  emailForm.value.email = user.value.email || ''
  emailError.value = ''
  showEmailModal.value = true
}

const editPhone = () => {
  phoneForm.value.phone = user.value.phone || ''
  phoneError.value = ''
  showPhoneModal.value = true
}

const closeModals = () => {
  showEmailModal.value = false
  showPhoneModal.value = false
  emailError.value = ''
  phoneError.value = ''
}

const saveEmail = async () => {
  if (!emailForm.value.email) {
    emailError.value = '请输入邮箱'
    return
  }

  saving.value = true
  emailError.value = ''

  try {
    const response = await userApi.updateProfile({
      email: emailForm.value.email
    })

    if (response.success) {
      user.value.email = response.data.email
      userStore.updateUser(response.data)
      closeModals()
    }
  } catch (error) {
    if (error.error === '邮箱已被使用') {
      emailError.value = '该邮箱已被其他用户使用'
    } else {
      emailError.value = error.error || '保存失败，请重试'
    }
  } finally {
    saving.value = false
  }
}

const savePhone = async () => {
  if (!phoneForm.value.phone) {
    phoneError.value = '请输入手机号'
    return
  }

  if (!/^[0-9]{11}$/.test(phoneForm.value.phone)) {
    phoneError.value = '请输入正确的11位手机号'
    return
  }

  saving.value = true
  phoneError.value = ''

  try {
    const response = await userApi.updateProfile({
      phone: phoneForm.value.phone
    })

    if (response.success) {
      user.value.phone = response.data.phone
      userStore.updateUser(response.data)
      closeModals()
    }
  } catch (error) {
    if (error.error === '手机号已被使用') {
      phoneError.value = '该手机号已被其他用户使用'
    } else {
      phoneError.value = error.error || '保存失败，请重试'
    }
  } finally {
    saving.value = false
  }
}

const handleLogout = () => {
  userApi.logout()
  userStore.clearAuth()
  router.push('/login')
}
</script>

<style scoped>
.profile-container {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.profile-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 500px;
  padding: 30px;
}

.profile-header {
  text-align: center;
  margin-bottom: 30px;
}

.avatar {
  width: 80px;
  height: 80px;
  margin: 0 auto 16px;
}

.avatar img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

.avatar-placeholder {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 36px;
  color: white;
  font-weight: bold;
}

.profile-header h2 {
  margin: 0 0 4px 0;
  font-size: 24px;
  color: #333;
}

.username {
  color: #666;
  font-size: 14px;
  margin: 0;
}

.profile-stats {
  display: flex;
  justify-content: space-around;
  padding: 20px 0;
  border-top: 1px solid #eee;
  border-bottom: 1px solid #eee;
  margin-bottom: 20px;
}

.stat-item {
  text-align: center;
}

.stat-label {
  display: block;
  font-size: 12px;
  color: #999;
  margin-bottom: 4px;
}

.stat-value {
  display: block;
  font-size: 20px;
  font-weight: bold;
  color: #333;
}

.profile-info {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 24px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-item.editable {
  cursor: pointer;
  padding: 8px;
  border-radius: 8px;
  transition: background 0.2s;
}

.info-item.editable:hover {
  background: #f9f9f9;
}

.info-item label {
  font-size: 14px;
  color: #666;
}

.info-item span {
  font-size: 14px;
  color: #333;
}

.info-value {
  display: flex;
  align-items: center;
  gap: 8px;
}

.edit-icon {
  font-size: 14px;
  opacity: 0;
  transition: opacity 0.2s;
}

.editable:hover .edit-icon {
  opacity: 1;
}

.logout-btn {
  width: 100%;
  padding: 12px;
  background: #f44336;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: opacity 0.2s;
}

.logout-btn:hover {
  opacity: 0.9;
}

/* 弹窗样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 400px;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #eee;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  color: #999;
  cursor: pointer;
  padding: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  color: #333;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.form-group input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  box-sizing: border-box;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
}

.error-message {
  background: #fee;
  color: #c33;
  padding: 10px;
  border-radius: 6px;
  font-size: 14px;
  margin-bottom: 16px;
}

.modal-footer {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.cancel-btn,
.save-btn {
  padding: 10px 20px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: opacity 0.2s;
}

.cancel-btn {
  background: #f5f5f5;
  color: #666;
}

.cancel-btn:hover {
  background: #eee;
}

.save-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.save-btn:hover:not(:disabled) {
  opacity: 0.9;
}

.save-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
