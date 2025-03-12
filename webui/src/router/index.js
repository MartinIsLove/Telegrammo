import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import Login from '../views/Login.vue'
import Home from '../views/Home.vue'
import Profile from '../views/Profile.vue'
import Chat from '../views/CreateChat.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: Login},
		{path: '/home', component: Home},
		{path: '/profile', component: Profile},
		{path: '/chat', component: Chat},

	]
})

export default router
