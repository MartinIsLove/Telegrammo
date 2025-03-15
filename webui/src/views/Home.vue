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
		<div v-for="chat in chats" >
			
			<ol class="list-group list-group mt-1">
				<button class="btn btn-secondary w-100 " @click="openChat(chat)">
				<li class="list-group-item d-flex justify-content-between align-items-start">
					
					<div class="me-auto row custom-box">
						
						<div class="col-4">
							<img v-if="chat.propic" :src="'data:image/png;base64,'+ chat.propic" class="img-fluid rounded-circle" width="140" height="140" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="Placeholder" preserveAspectRatio="xMidYMid slice" focusable="false"><title>Placeholder</title><rect width="100%" height="100%" fill="var(--bs-secondary-color)"></rect>
						</div>
						<div class="col-auto">
							<span class="fw-bold">
								{{chat.nome}}
							</span>
							
							<br>

							<span v-if="chat.gruppo">
								<span class="text-capitalize fw-bold">{{ chat.username }}</span>:{{ chat.lastmsg }}  
							</span>

							<span v-else>
								<span>{{ chat.lastmsg  }}</span>
								
							</span>
						</div>
					<!-- </button> -->
					</div>
					<!-- </button> -->
					<span class="badge text-bg-secondary rounded-pill">14</span>
				</li>
				</button>
			</ol>
			<!-- </button> -->
		</div>
	
</template>

<style>
.custom-box {
    min-height: 60px;  /* Imposta un'altezza minima di 300px */
}

</style>