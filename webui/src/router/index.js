import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import Login from '../views/Login.vue'
import Home from '../views/Home.vue'
import Profile from '../views/Profile.vue'
import CreateChat from '../views/CreateChat.vue'
import Chat from '../views/Chat.vue'
import Options from '../views/GroupOptions.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: Login},
		{path: '/home', component: Home},
		{path: '/profile', component: Profile},
		{path: '/createChat', component: CreateChat},
		{path: '/chat/:id', component: Chat},
		{path: '/options/group/:id', component: Options},
	]
})

export default router
