<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,

			messaggi:[],
            photo: [],
            messaggio: '',
			id: null,
            chatId: null,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
            this.id = sessionStorage.getItem("cs")
            this.chatId = parseInt(this.$route.params.id, 10);
			this.fetchMessage()

		},
        async fetchMessage(){
            let response = await this.$axios.get(`/conversation/${this.chatId}`, {headers:{cs:this.id}});
			this.messaggi = response.data

        },
		async send() {
            let response = await this.$axios.post("/conversation/message", {id_chat:this.chatId, testo:this.message, photo:this.photo}, {headers: {cs:this.id}});

            this.fetchMessage()
		},
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
    
	<div v-for="mes in messaggi">
		<div class="row g-0 border rounded my-2 p-2">
			<div class="d-flex justify-content-between">
				<p v-if="mes.gruppo" class="text-wrap">
					<span class="text-capitalize fw-bold">{{ mes.nome }}</span>
					{{ mes.testo }}
				</p>
				
				<p v-else> 
					<span>{{ mes.testo  }}</span>
				</p>
			</div>
		</div>
	</div>

<div class="min-vh-100">
</div>
<form  @submit.prevent="send" class="sticky-bottom mb-2">
<div >
    <div class="input-group mb-3">
    <input type="text" class="form-control rounded" placeholder="write message" aria-label="message" aria-describedby="basic-addon1" v-model="message">
    <button class="btn btn-secondary mt-2 rounded ms-2"  >Send</button>
    <br>
    </div>
</div>
</form>
  
    
</template>