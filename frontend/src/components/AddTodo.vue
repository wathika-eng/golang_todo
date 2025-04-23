<script setup>
import { ref } from 'vue';
import api from '../services/api';

const newTodo = ref('');
const loading = ref(false);
const error = ref(null);
const emit = defineEmits(['todoAdded']);

// Function to submit a new todo
async function addTodo() {
    if (!newTodo.value.trim()) {
        error.value = "Todo cannot be empty!";
        return;
    }

    loading.value = true;
    error.value = null;

    try {
        const response = await api.post('/api/notes/create', { notes_body: newTodo.value });

        if (!response.data.error) {
            emit('todoAdded', response.data); // Notify parent component
            newTodo.value = ''; // Clear input
        } else {
            error.value = "Failed to add todo!";
        }
    } catch (err) {
        error.value = "Error adding todo!";
        console.error(err);
    } finally {
        loading.value = false;
    }
}
</script>

<template>
    <div class="add-todo-container">
        <h2>Add New Todo</h2>
        <input v-model="newTodo" placeholder="Enter a new todo..." />
        <button @click="addTodo" :disabled="loading">
            {{ loading ? "Adding..." : "Add Todo" }}
        </button>
        <p v-if="error" class="error">{{ error }}</p>
    </div>
</template>

<style scoped>
.add-todo-container {
    max-width: 400px;
    margin: 20px auto;
    text-align: center;
}

input {
    width: 80%;
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 5px;
    font-size: 1rem;
}

button {
    margin-left: 10px;
    padding: 10px 15px;
    border: none;
    background: #28a745;
    color: white;
    cursor: pointer;
    font-size: 1rem;
    border-radius: 5px;
}

button:disabled {
    background: gray;
    cursor: not-allowed;
}

.error {
    color: red;
    margin-top: 10px;
}
</style>
