<script setup>
import { ref, onMounted } from 'vue';
import api from '../services/api';
import { useRouter } from 'vue-router';
import AddTodo from './AddTodo.vue';
import dayjs from 'dayjs';

const notes = ref([]);
const loading = ref(true);
const error = ref(null);
const editingNote = ref(null);
const editText = ref("");
const router = useRouter();

// Fetch notes from API
async function fetchNotes() {
    loading.value = true;
    error.value = null;

    try {
        const response = await api.get("/api/notes/");
        if (response.data.error === false) {
            console.log(response.data.data);
            notes.value = response.data.data; // Extract notes from `todos`
        } else {
            error.value = response.data.message || "Failed to load notes!";
        }
    } catch (err) {
        error.value = "Error fetching notes!";
        console.error(err);
    } finally {
        loading.value = false;
    }
}
// Format date using 
function formatDateDirectly(isoString) {
    return isoString.replace("T", " ").split(".")[0];
}

// Mark note as completed
async function toggleCompleted(note) {
    try {
        const response = await api.patch(`/api/notes/${note.notes_id}`, { completed: !note.completed });
        if (!response.data.error) {
            note.completed = response.data.todo.completed;
        }
    } catch (err) {
        alert("Error updating note!");
        console.error(err);
    }
}

// Delete a note
async function deleteNote(noteId) {
    if (!confirm("Are you sure you want to delete this note?")) return;

    try {
        await api.delete(`/api/notes/${noteId}`);
        notes.value = notes.value.filter(note => note.notes_id !== noteId);
    } catch (err) {
        alert("Error deleting note!");
        console.error(err);
    }
}

// Start editing a note
function startEditing(note) {
    editingNote.value = note.notes_id;
    editText.value = note.notes_body;
}

async function saveEdit(note) {
    try {
        const response = await api.patch(`/api/notes/${note.notes_id}`, {
            body: editText.value // Use the correct column name (e.g., "body" instead of "notes_body")
        });
        if (!response.data.error) {
            note.notes_body = response.data.todo.body; // Update the note's body in the UI
            editingNote.value = null; // Exit edit mode
            await fetchNotes(); // Refresh notes
        }
    } catch (err) {
        alert("Error updating note!");
        console.error(err);
    }
}

// Logout
async function logout() {
    try {
        await api.post("/api/notes/logout");
        localStorage.removeItem("access_token");
        localStorage.removeItem("refresh_token");
        localStorage.removeItem("todo");
        localStorage.removeItem("todos");
        router.push('/login');
    } catch (err) {
        alert("Error logging out!");
        console.error(err);
    }
}

onMounted(fetchNotes);
</script>

<template>
    <div class="container">
        <h1>Your Notes</h1>
        <AddTodo @todoAdded="fetchNotes" />
        <button class="logout-btn" @click="logout">Logout</button>

        <div v-if="loading" class="loading">Loading notes...</div>
        <div v-if="error" class="error">{{ error }}</div>

        <ul v-if="!loading && !error">
            <li v-for="note in notes" :key="note.notes_id" class="note-card">
                <div v-if="editingNote === note.notes_id">
                    <textarea v-model="editText"></textarea>
                    <button class="save-btn" @click="saveEdit(note)">Save</button>
                    <button class="cancel-btn" @click="editingNote = null">Cancel</button>
                </div>
                <div v-else>
                    <p class="note-body">{{ note.notes_body }}</p>
                    <small class="note-date">
                        Created: {{ formatDateDirectly(note.CreatedAt) }}
                    </small>
                    <small v-if="note.UpdatedAt && note.UpdatedAt !== note.CreatedAt" class="note-date">
                        | Updated: {{ formatDateDirectly(note.UpdatedAt) }}
                    </small>
                    <div class="actions">
                        <button class="complete-btn" :class="{ completed: note.completed }"
                            @click="toggleCompleted(note)">
                            {{ note.completed ? "✔ Completed" : "Mark as Completed" }}
                        </button>
                        <button class="edit-btn" @click="startEditing(note)">✏ Edit</button>
                        <button class="delete-btn" @click="deleteNote(note.notes_id)">🗑 Delete</button>
                    </div>
                </div>
            </li>
        </ul>
    </div>
</template>

<style scoped>
.container {
    max-width: 600px;
    margin: 40px auto;
    padding: 20px;
    background: #fff;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    border-radius: 8px;
}

h1 {
    text-align: center;
    color: #333;
}

.loading,
.error {
    text-align: center;
    font-size: 1.2rem;
    color: red;
}

.note-card {
    list-style: none;
    padding: 15px;
    border: 1px solid #ddd;
    border-radius: 8px;
    margin-bottom: 10px;
    background: #f9f9f9;
    transition: 0.3s;
}

.note-card:hover {
    background: #f1f1f1;
}

.note-body {
    font-size: 1.1rem;
    color: #333;
}

.note-date {
    font-size: 0.8rem;
    color: #777;
    display: block;
    margin-top: 5px;
}

.actions {
    margin-top: 10px;
    display: flex;
    justify-content: space-between;
}

button {
    padding: 8px 12px;
    border: none;
    cursor: pointer;
    border-radius: 5px;
    font-size: 0.9rem;
}

.complete-btn {
    background: #007bff;
    color: white;
}

.complete-btn.completed {
    background: #28a745;
}

.edit-btn {
    background: #ffc107;
    color: black;
}

.delete-btn {
    background: #dc3545;
    color: white;
}

.save-btn {
    background: #28a745;
    color: white;
}

.cancel-btn {
    background: #777;
    color: white;
}

textarea {
    width: 100%;
    height: 80px;
    margin-bottom: 10px;
    padding: 8px;
    font-size: 1rem;
    border: 1px solid #ddd;
    border-radius: 5px;
}

.logout-btn {
    display: block;
    margin: 10px auto;
    padding: 8px 12px;
    background: #ff4500;
    color: white;
    text-align: center;
    font-weight: bold;
}
</style>
