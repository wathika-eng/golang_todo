import { createMemoryHistory, createRouter, RouteRecordRaw } from 'vue-router';

import Login from '../components/Login.vue';
import Notes from '../components/Notes.vue';
import Signup from '../components/Signup.vue';

const routes: Array<RouteRecordRaw> = [
	{ path: '/login', component: Login },
	{ path: '/signup', component: Signup },
	{ path: '/notes', component: Notes },
	{ path: '/', redirect: '/login' }, // Redirect root path to login
];

const router = createRouter({
	history: createMemoryHistory(),
	routes,
});

export default router;
