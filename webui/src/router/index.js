import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import Login from '../views/Login.vue'
import Home from '../views/Home.vue'


const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: Login},
		{path: '/home', component: Home},
		{path: '/link2', component: HomeView},
		{path: '/some/:id/link', component: HomeView},

	]
})

export default router
