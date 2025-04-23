<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../services/api';

const router = useRouter();
const title = ref("Signup");
const username = ref("");
const password = ref("");
const repeatPassword = ref("");
const error = ref("");

async function signup(event) {
  event.preventDefault(); // Prevent page refresh

  if (!username.value || !password.value || !repeatPassword.value) {
    error.value = "All fields are required.";
    return;
  }

  if (password.value !== repeatPassword.value) {
    error.value = "Passwords do not match!";
    return;
  }

  try {
    const response = await api.post("/api/users/signup", {
      email: username.value,
      password: password.value,
    });

    if (response.data.error === false) {
      alert("Signup successful! Please login.");
      router.push("/login"); // Redirect to login page
    } else {
      error.value = response.data.message || "Signup failed!";
    }
  } catch (err) {
    error.value = err.response?.data?.message || "Signup failed!";
    console.error("Signup error:", err);
  }
}
</script>

<template>
  <div class="container">
    <h1>{{ title }}</h1>
    <p v-if="error" class="error">{{ error }}</p>

    <form @submit="signup">
      <label for="username">Email</label>
      <input v-model="username" id="username" placeholder="Enter your email" type="email" required>

      <label for="password">Password</label>
      <input v-model="password" id="password" placeholder="Enter password" type="password" required>

      <label for="repeat-password">Repeat Password</label>
      <input v-model="repeatPassword" id="repeat-password" placeholder="Repeat password" type="password" required>

      <button class="button" type="submit">Signup</button>
    </form>

    <p class="toggle-text">
      Already have an account? <RouterLink to="/login">Login here</RouterLink>.
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