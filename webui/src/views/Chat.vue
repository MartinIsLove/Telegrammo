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

			this.$nextTick(() => {
                this.scrollToBottom();
            });
        },
		async send() {
            await this.$axios.post("/conversation/message", {id_chat:this.chatId, testo:this.message, photo:this.photo}, {headers: {cs:this.id}});

            this.fetchMessage()
		},
		formatDate(date) {
            if (!date) return '';
            const d = new Date(date);
            return d.toLocaleDateString('it-IT', {
                day: '2-digit',
                month: '2-digit',
                year: 'numeric'
            });
        },
		scrollToBottom() {
            const container = this.$refs.messageContainer;
            if (container) {
                container.scrollTop = container.scrollHeight;
            }
        }
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	
	<div class="bg-dark" style="width: 100%; position: sticky; top: 0; display: flex; flex-direction: column; justify-content: center; align-items: center; ">
        <h2 class="mt-2 fw-bold" style="text-align: center;">{{ messaggi.nome }}</h2>
        <hr style="width: 100%;">
    </div>
	<!-- <div style="height: 8%;"></div> -->
	<div ref="messageContainer"  style="overflow-y: auto; max-height: calc(100vh - 150px);">
		<div v-for="mes in messaggi.message">
			<div v-if="mes.testo === '' && mes.photo === null" class="d-flex justify-content-center my-2">
				<div class="date-container">
					<p class="text-wrap text-center mb-0">
						{{ formatDate(mes.data) }}
					</p>
				</div>
			</div>
			<div v-else-if="id == mes.id_mitt" class="d-flex justify-content-end my-2 ">
				<div class="message-container row g-0 border rounded p-2">
					<div class="d-flex justify-content-between w-100">
						<p v-if="mes.gruppo" class="text-wrap message-content">
							<span class="text-capitalize fw-bold">{{ mes.nome }}</span><br>
							{{ mes.testo }}
						</p>
						<p v-else class="message-content">
							<span>{{ mes.testo }}</span>
						</p>
					</div>
				</div>
			</div>
			<div v-else class="row g-0 border rounded my-2 p-2 message-container">
				<div class="d-flex justify-content-between w-100">
					<p v-if="mes.gruppo" class="text-wrap message-content">
						<span class="text-capitalize fw-bold">{{ mes.nome }}</span><br>
						{{ mes.testo }}
					</p>
					<p v-else class="message-content">
						<span>{{ mes.testo }}</span>
					</p>
				</div>
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
<style>

.text-center {
    text-align: center;
}
.date-container {
    display: inline-block;
    padding: 10px; 
    background-color: #1b1616;
    border-radius: 5px; 
}
.message-container {
    width: 40%; 
    margin: 0 auto; 
	
}
.message-content {
    word-break: break-word; 
    white-space: pre-wrap; 
}
body, html {
    overflow: hidden; 
    height: 100%; 
}
</style>