<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,

			id: null,
			chats: [],
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				this.id = sessionStorage.getItem("cs")
				let response = await this.$axios.get("/conversations", {headers: {cs:this.id}});
				this.chats = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		createChatButtonHandler() {
			this.$router.push("/createChat");
		},
		openChat(chat){
			console.log("Chat aperta:", chat);
		}
	},
	mounted() {
		this.refresh()
	}
}
</script>
<template>
    
	<div class="position-absolute top-1 start-1">
		<button class="btn btn-secondary ms-2 mb-2" @click="createChatButtonHandler">
			Create Chat
		</button>
	</div>
	<div class="mt-5">

	</div>
	<!-- <div v-for="chat in chats" >
		<ol class="list-group list-group mt-1">
			<li class="list-group-item d-flex justify-content-between align-items-start">
				
				<div class="me-auto row custom-box">
					
					<div class="col-2 ms-0 me-2">
						<img v-if="chat.propic" :src="'data:image/png;base64,'+ chat.propic" class=" rounded-circle me-2" width="70" height="70"  role="img" focusable="false" style="object-fit: cover;">
					</div>
					<div class="col-10 d-flex">
	
						<p class="fw-bold">
							
							{{chat.nome}}
						</p>
													
						<p v-if="chat.gruppo" class="text-wrap">
							<span class="text-capitalize fw-bold">{{ chat.username }}</span>:{{ chat.lastmsg }}  
						</p>
						
						<p v-else>
							<span>{{ chat.lastmsg  }}</span>
							
						</p>
				
					</div>
				</div>
				<span class="badge text-bg-secondary rounded-pill">14</span>
			</li>
		</ol>
	</div> -->

	<div v-for="chat in chats" >
		<RouterLink :to="'/chat/' + chat.id_chat" class="nav-link">
		<div class="row g-0 border rounded my-2 p-2">


				
				<div class="col-md-1 col-2">
					<img v-if="chat.propic" :src="'data:image/png;base64,'+ chat.propic" class=" rounded-circle me-2" width="70" height="70"  role="img" focusable="false" style="object-fit: cover;">
				</div>
				
				<div class="col-md-11 col-10">
					<div class="card-body">
						
						<h5>{{chat.nome}}</h5>
						
						<div class="d-flex justify-content-between">
							<p v-if="chat.gruppo" class="text-wrap">
								<span class="text-capitalize fw-bold">{{ chat.username }}</span>:{{ chat.lastmsg }}  
							</p>
							
							<p v-else> 
								<span>{{ chat.lastmsg  }}</span>
							</p>
						</div>
						
					</div>
				</div>
				
				
			</div>
		</RouterLink>
		
	</div>

</template>

<style>
.custom-box {
    min-height: 60px;  /* Imposta un'altezza minima di 300px */
}

</style>