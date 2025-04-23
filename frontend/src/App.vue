<script setup>
import { shallowRef, ref, onMounted } from 'vue';
import Login from './components/Login.vue';
import Signup from './components/Signup.vue';
import Notes from './components/Notes.vue';
import api from './services/api';

const status = shallowRef(null);
const isLoggedIn = ref(false);
const apiStatus = ref('Checking API status...');
const apiStatusColor = ref('gray');

async function checkApiStatus() {
  try {
    const response = await api('/api/users/test');
    if (response.status === 200) {
      apiStatus.value = 'API is live';
      apiStatusColor.value = 'green';
    } else {
      apiStatus.value = 'API is down';
      apiStatusColor.value = 'red';
    }
  } catch (error) {
    apiStatus.value = 'API is down';
    apiStatusColor.value = 'red';
  }
}

function checkLoginStatus() {
  const accessToken = localStorage.getItem("access_token");
  isLoggedIn.value = !!accessToken;
}

function showSignup() {
  status.value = Signup;
}

function showLogin() {
  status.value = Login;
}

function logout() {
  localStorage.removeItem("access_token");
  localStorage.removeItem("refresh_token");
  isLoggedIn.value = false;
  status.value = "";
}

onMounted(() => {
  checkLoginStatus();
  checkApiStatus();
});
</script>

<template>
  <div class="api-banner" :style="{ backgroundColor: apiStatusColor }">
    {{ apiStatus }}
  </div>

  <div class="auth-buttons" v-if="!isLoggedIn">
    <button @click="showSignup">Signup</button>
    <button @click="showLogin">Login</button>
  </div>

  <div v-if="isLoggedIn">
    <button @click="logout">Logout</button>
  </div>

  <main>
    <component :is="status" v-if="!isLoggedIn" />
    <Notes v-if="isLoggedIn" />
  </main>
</template>

<style scoped>
.api-banner {
  text-align: center;
  padding: 10px;
  color: white;
  font-weight: bold;
  margin-bottom: 20px;
}

.auth-buttons {
  text-align: center;
  margin-bottom: 20px;
}

button {
  margin: 5px;
  padding: 10px 15px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

button:hover {
  background: #0056b3;
}

main {
  margin-top: 20px;
}
</style>
