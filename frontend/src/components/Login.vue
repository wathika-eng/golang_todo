<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../services/api';

const router = useRouter();
const title = ref("Login");
const username = ref("qowalu@gmail.com");
const password = ref("Pa$$w0rd!");
const error = ref("");

async function login(event) {
  event.preventDefault();

  if (!username.value || !password.value) {
    error.value = "Please fill in all fields.";
    return;
  }

  try {
    const response = await api.post("/api/users/login", {
      email: username.value,
      password: password.value,
    });

    // Extract tokens
    const accessToken = response.data.access_token;
    const refreshToken = response.data.refresh_token;

    // Store tokens in localStorage
    localStorage.setItem("access_token", accessToken);
    localStorage.setItem("refresh_token", refreshToken);

    alert(`Welcome back, ${username.value}!`);
    router.push("/notes"); // Redirect to notes page
  } catch (err) {
    error.value = err.response?.data?.message || "Login failed!";
    console.error("Login error:", err);
  }
}
</script>

<template>
  <div class="container">
    <h1>{{ title }}</h1>
    <p v-if="error" class="error">{{ error }}</p>

    <form @submit="login">
      <label for="username">Email</label>
      <input v-model="username" id="username" placeholder="Enter your email" type="email" required>

      <label for="password">Password</label>
      <input v-model="password" id="password" placeholder="Enter password" type="password" required>

      <button class="button" type="submit">Login</button>
    </form>

    <p class="toggle-text">
      Don't have an account? <RouterLink to="/signup">Signup here</RouterLink>.
    </p>
  </div>
</template>

<style scoped>
.container {
  max-width: 400px;
  margin: auto;
  padding: 20px;
  background: #f9f9f9;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  text-align: center;
}

h1 {
  margin-bottom: 20px;
}

form {
  display: flex;
  flex-direction: column;
}

label {
  font-size: 14px;
  margin-top: 10px;
  text-align: left;
}

input {
  padding: 10px;
  margin-top: 5px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 16px;
}

.button {
  margin-top: 20px;
  padding: 10px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
}

.button:hover {
  background: #0056b3;
}

.error {
  color: red;
  font-size: 14px;
  margin-bottom: 10px;
}

.toggle-text {
  margin-top: 20px;
  font-size: 14px;
}

.toggle-text a {
  color: #007bff;
  text-decoration: none;
}

.toggle-text a:hover {
  text-decoration: underline;
}
</style>