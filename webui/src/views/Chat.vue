<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,

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
		},
        async fetchMessage(){
            


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