<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,

			messaggi:[],
            messaggio: '',
			id: null,
            chatId: null,
			photo: null,
			message: "",
			reply: false,
			messageReplied: -1,

			emoji: ["👠", "💅", "👍🏻", "👌🏻","❤️"],
			intervalId: null,
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
		async fetchMessageNoBottom(){
            let response = await this.$axios.get(`/conversation/${this.chatId}`, {headers:{cs:this.id}});
			this.messaggi = response.data
			this.messaggi.message.forEach(mes => {
				mes.showEmoji = false; 
			});
        },
        async fetchMessage(){
            let response = await this.$axios.get(`/conversation/${this.chatId}`, {headers:{cs:this.id}});
			this.messaggi = response.data
			this.messaggi.message.forEach(mes => {
				mes.showEmoji = false; 
			});
			this.$nextTick(() => {
                this.scrollToBottom();
            });
        },
		async send() {
			if (this.photo === null && this.message === ""){
				return 
			}	
			
			let formData = new FormData();
			formData.append('photo', this.photo)
			formData.append('id_chat', this.chatId)
			formData.append('testo', this.message)
			if (this.messageReplied === -1){
				formData.append('id_reply', this.messageReplied)
				await this.$axios.post("/message", formData, {headers: {'Content-Type': 'multipart/form-data', cs: this.id}});

			}else{
				formData.append('id_reply', this.messageReplied.id_forward)
				await this.$axios.post("/message", formData, {headers: {'Content-Type': 'multipart/form-data', cs: this.id}});
			}
			this.message = ''
			this.photo = null 
            this.$refs.photoInput.value = '' 
			this.reply = false;
			this.messageReplied = -1;
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
		formatTime(date) {
            if (!date) return '';
            const d = new Date(date);
            return d.toLocaleTimeString('it-IT', {
                hour: '2-digit',
                minute: '2-digit'
            });
        },
		scrollToBottom() {
            const container = this.$refs.messageContainer;
            if (container) {
                container.scrollTop = container.scrollHeight;
            }
        },
		optionsButtonHandler() {
			this.$router.push("/options/group/"+ this.chatId);
		},
		handleFileUpload(event) {
			const file = event.target.files[0];
			if (file) {
				this.photo = file;
			}
		},
		removefile(){
			this.photo = null;
		},
		async handleCheckboxChange(mes, emojiItem, event) {
			if (event) event.preventDefault(); // Prevent default behavior
			try {
				if (emojiItem !== "") {
					await this.$axios.post(`/message/comment`, {
						id_chat: this.chatId, 
						emoji: emojiItem, 
						id_mes: mes.id_mess
					}, {
						headers: {cs: this.id}
					});
				}
				this.fetchMessageNoBottom();
			} catch(e) {
				console.error(e);
			}
		},
		getEmojiCount(emojiItem){
			if (emojiItem === "👠"){
				return 0
			}if (emojiItem === "❤️"){
				return 1
			}if (emojiItem === "👍🏻"){
				return 2
			}if (emojiItem === "👌🏻"){
				return 3
			}if (emojiItem === "💅"){
				return 4
			}

		},
		async handleCheckboxChangeUncomment(mes, event) {
			if (event) event.preventDefault(); // Prevent default behavior
			try {
				await this.$axios.post(`/message/uncomment`, {
					id_chat: this.chatId,
					id_mes: mes.id_mess
				}, {
					headers: {cs: this.id}
				});
				this.fetchMessageNoBottom();
			} catch(e) {
				console.error(e);
			}
		},
		async forwardMessage(id_mes, id_utente){
			this.$router.push("/chat/"+this.chatId+"/forward/"+id_mes+"/"+id_utente );

		},
		async deleteMessage(id_mes, id_forw){
			
			try{
				// await this.$axios.delete(`/message/delete/${id_mes}`, {id_forward:id_forw, id_chat:this.chatId},{headers:{cs:this.id}});
				await this.$axios({
					method: 'delete',
					url: `/message/delete/${id_mes}`,
					data: {
						id_forward: id_forw,
						id_chat: this.chatId,
					},
					headers: {
						cs: this.id, 
					},
        		});
				this.fetchMessageNoBottom()
			}
			catch(e){
				console.error(e);
			}
		},
		async replyMessage(mes){

			console.log("dio")
			this.reply=true
			this.messageReplied=mes
			

			},
		closeReply() {
			this.reply = false;
			this.messageReplied = -1;
		},
		
	},
	mounted() {
		this.refresh()
		this.intervalId=setInterval(() => {
            this.fetchMessageNoBottom();
        }, 2000); 

    
	},
	beforeRouteLeave(to, from, next) {
        clearInterval(this.intervalId); 
        next();
    },
}
</script>

<template>
  	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bootstrap-icons/1.8.1/font/bootstrap-icons.min.css">
	
	<div class="bg-dark" style="width: 100%; position: sticky; top: 0;  ">
		<div class="d-flex justify-content-between align-items-center w-100">
			<div v-if="messaggi.gruppo == true" class="d-flex justify-content-start align-items-center">
				
				<button class="btn btn-secondary mt-2 rounded text-tart" @click="optionsButtonHandler">
					options
				</button>
				
			</div>
			<h2 class="mt-2 position-relative fw-bold flex-grow-1 text-center " style="transform: translateX(-40px);" >{{ messaggi.nome }}</h2>
		</div>
		<hr style="width: 100%;">
    </div>


	<div ref="messageContainer"  style="overflow-y: auto; max-height: calc(100vh - 150px); min-height: 80vh;">
		
		<div v-for="mes in messaggi.message" :key="mes.id_mess">
			
			<div v-if="mes.testo === '' && mes.photo === null" class="d-flex justify-content-center my-2">
				<div class="date-container">
					<p class="text-wrap text-center mb-0">
						{{ formatDate(mes.data) }}
					</p>
				</div>
			</div>
			
			<!-- ------------------------------------------------------------------------------------------------------------------------------------------ -->
			<div v-else-if="(id==mes.id_mitt && !(mes.forward_id!=-1 && mes.forward_id_mit!=id)) || (mes.forward_id!=-1 && mes.forward_id_mit==id) " class="d-flex justify-content-end my-2">

				<div class=" message-container row g-0 border rounded p-2 ">
					
					<div class="row pe-0 mb-2" style="width: 100%;">
						<div v-if="mes.forward_id!=-1 " class="col-auto ps-0">
							<span v-if="mes.gruppo" class="text-capitalize fw-bold d-flex justify-content-start">{{mes.forward_user_mit }}</span>
						</div>
						<div v-else class="col-auto ps-0">
							<span v-if="mes.gruppo" class="text-capitalize fw-bold d-flex justify-content-start">{{ mes.nome }}</span>
						</div>

						<div class="col-auto ms-auto pe-0">
							<button class="forward-button btn btn-secondary " @click="replyMessage(mes)">
								<i class="bi bi-arrow-90deg-left"></i>
							</button>

							<button class="forward-button btn btn-secondary mx-2" @click="forwardMessage(mes.id_mess, mes.id_mitt)">
								<i class="bi bi-arrow-90deg-right"></i>
							</button>
						
							<button class="trash-button btn btn-secondary" @click="deleteMessage(mes.id_mess, mes.id_forward)">
								<i class="bi bi-trash"></i>
							</button>
						</div>
						<div v-if="mes.id_reply!==-1" class="reply-container mt-2">
							
							<span v-if="mes.gruppo" class="text-capitalize fw-bold d-flex justify-content-start ms-2">{{ mes.mit_reply }}</span>
							<div class="mb-2 ms-2" v-if="mes.photo_reply !== null">
								<img :src="'data:image/png;base64,' + mes.photo_reply" class="img-fluid rounded" />
							</div >
							<p class="ms-2">
								{{ mes.testo_reply }}
							</p>
						</div>
					
					</div>

					<div class="mb-2" v-if="mes.photo !== null">
						<img :src="'data:image/png;base64,' + mes.photo" class="img-fluid rounded" />
					</div>

					<p>{{ mes.testo }}</p>
						
					<div class="row pe-0">

						<div class="emoji-checkbox-container col-auto ps-0">
							<div v-for="emojiItem in emoji" :key="emojiItem" :value="emojiItem" class="emoji-item"> 
								<div v-if="emojiItem === mes.myEmoji">
									<input  type="checkbox" class="btn-check" :id="'btn-check-' + mes.id_mess + '-' + emojiItem" autocomplete="off" @change="handleCheckboxChangeUncomment(mes, $event)">
									<label class="btn btn-primary mb-2" :for="'btn-check-' + mes.id_mess + '-' + emojiItem">{{ emojiItem }} {{mes.emoji[getEmojiCount(emojiItem)]}}</label>	
								</div>
								<div v-else>
									<input  type="checkbox" class="btn-check" :id="'btn-check-' + mes.id_mess + '-' + emojiItem" autocomplete="off" @change="handleCheckboxChange(mes, emojiItem, $event)">
									<label class="btn btn-secondary mb-2" :for="'btn-check-' + mes.id_mess + '-' + emojiItem">{{ emojiItem }} {{mes.emoji[getEmojiCount(emojiItem)]}}</label>	
								</div>
							</div>
						</div>
						<div class="col-auto ms-auto pe-0 pt-3">
							<small v-if=" mes.forward_id === -1" class="text-end"> {{ formatTime(mes.data) }} </small>
							<small v-else class="text-end"> {{ formatTime(mes.forward_date) }} </small>
							<span v-if="!mes.visual">
								<i class="bi bi-check"></i>
								
							</span>
							<span v-else>
								<i class="bi bi-check-all"></i>
							</span>
						</div>
						<div v-if="mes.forward_id !== -1">
							<p class="forwarded-by">
								forwarded by: {{ mes.forward_username }}
							</p>
						</div>

                    </div>
				</div>
			</div>

			<!-- ------------------------------------------------------------------------------------------------------------------------------------------  -->
			<div v-else-if="id!==mes.id_mit || (mes.forward_id!=-1 && mes.forward_id_mit!=id)" class="d-flex justify-content-start my-2">

				<div class=" message-container row g-0 border rounded p-2 ">
					
					<div class="row pe-0 mb-2" style="width: 100%;">
						<div v-if="mes.forward_id!=-1 " class="col-auto ps-0">
							<span v-if="mes.gruppo" class="text-capitalize fw-bold d-flex justify-content-start">{{mes.forward_user_mit }}</span>
						</div>
						<div v-else class="col-auto ps-0">
							<span v-if="mes.gruppo" class="text-capitalize fw-bold d-flex justify-content-start">{{ mes.nome }}</span>
						</div>

						<div class="col-auto ms-auto pe-0">
							<button class="forward-button btn btn-secondary " @click="replyMessage(mes)">
								<i class="bi bi-arrow-90deg-left"></i>
							</button>

							<button class="forward-button btn btn-secondary mx-2" @click="forwardMessage(mes.id_mess, mes.id_mitt)">
								<i class="bi bi-arrow-90deg-right"></i>
							</button>
						
							<button class="trash-button btn btn-secondary" @click="deleteMessage(mes.id_mess, mes.id_forward)">
								<i class="bi bi-trash"></i>
							</button>
						</div>
						<div v-if="mes.id_reply!==-1" class="reply-container mt-2">
							
							<span v-if="mes.gruppo" class="text-capitalize fw-bold d-flex justify-content-start ms-2">{{ mes.mit_reply }}</span>
							<div class="mb-2 ms-2" v-if="mes.photo_reply !== null">
								<img :src="'data:image/png;base64,' + mes.photo_reply" class="img-fluid rounded" />
							</div >
							<p class="ms-2">
								{{ mes.testo_reply }}
							</p>
						</div>
					
					</div>

					<div class="mb-2" v-if="mes.photo !== null">
						<img :src="'data:image/png;base64,' + mes.photo" class="img-fluid rounded" />
					</div>

					<p>{{ mes.testo }}</p>

						
					<div class="row pe-0">

						<div class="emoji-checkbox-container col-auto ps-0">
							<div v-for="emojiItem in emoji" :key="emojiItem" :value="emojiItem" class="emoji-item"> 
								<div v-if="emojiItem === mes.myEmoji">
									<input  type="checkbox" class="btn-check" :id="'btn-check-' + mes.id_mess + '-' + emojiItem" autocomplete="off" @change="handleCheckboxChangeUncomment(mes, $event)">
									<label class="btn btn-primary mb-2" :for="'btn-check-' + mes.id_mess + '-' + emojiItem">{{ emojiItem }} {{mes.emoji[getEmojiCount(emojiItem)]}}</label>	
								</div>
								<div v-else>
									<input  type="checkbox" class="btn-check" :id="'btn-check-' + mes.id_mess + '-' + emojiItem" autocomplete="off" @change="handleCheckboxChange(mes, emojiItem, $event)">
									<label class="btn btn-secondary mb-2" :for="'btn-check-' + mes.id_mess + '-' + emojiItem">{{ emojiItem }} {{mes.emoji[getEmojiCount(emojiItem)]}}</label>	
								</div>
							</div>
						</div>
						<div class="col-auto ms-auto pe-0 pt-3">

							<small v-if=" mes.forward_id === -1" class="text-end"> {{ formatTime(mes.data) }} </small>
							<small v-else class="text-end"> {{ formatTime(mes.forward_date) }} </small>
							
						</div>
						<div v-if="mes.forward_id !== -1">
							<p class="forwarded-by">
								forwarded by: {{ mes.forward_username }}
							</p>
							
							
						</div>
					</div>

				</div>
			</div>
		</div>
	</div> 



	<div class="min-vh-100">
	</div>
	<form @submit.prevent="send" class="sticky-bottom mb-2">
		<div>
			<div v-if="reply" class="reply-container">
				<button class="reply-close-button" @click="closeReply">
					<i class="bi bi-x"></i>
				</button>
				<div class="row pe-0 mb-2">
					<div v-if="messageReplied.forward_id!=-1" class="col-auto ps-0">
						<span v-if="messageReplied.gruppo" class="text-capitalize fw-bold ms-2">{{messageReplied.forward_user_mit }}</span>
					</div>
					<div v-else class="col-auto ps-0">
						<span v-if="messageReplied.gruppo" class="text-capitalize fw-bold ms-2">{{ messageReplied.nome }}</span>
					</div>
				</div>
				
				<div v-if="messageReplied.photo !== null" class="mb-2">
					<img :src="'data:image/png;base64,' + messageReplied.photo" class="img-fluid rounded" style="max-height: 100px;" />
				</div>

				<p class="mb-0">{{ messageReplied.testo }}</p>
				<div v-if="messageReplied.forward_id !== -1">
							<p class="forwarded-by">
								forwarded by: {{ messageReplied.forward_username }}
							</p>
				</div>
			</div>
			
			<div class="input-group mb-3">
				<input type="text" class="form-control rounded" placeholder="write message" aria-label="message" aria-describedby="basic-addon1" v-model="message">
				<button class="btn btn-secondary mt-2 rounded ms-2"  >Send</button>
				<div v-if="!photo">
					<label for="photo" class="btn btn-secondary mt-2 rounded ms-2 btn-paperclip">
						<i v-if="!photo" class="bi bi-paperclip"></i>
						<!-- <i v-else class="bi bi-image"></i> -->
						
						<input  type="file" class="d-none" id="photo" @change="handleFileUpload"  ref="photoInput" v-show="true">
					
					</label>
					<br>
				</div>
				<div v-else>
					<button @click="removefile" class="btn btn-secondary mt-2 rounded ms-2 btn-paperclip">
						<i  class="bi bi-image"></i>
					</button>
				</div>
			</div>
		</div>
	</form>
</template>

<style>
.forward-button {
    /* position: absolute;
    top: 10px;
    right: 10px; */
	flex-shrink: 0;
}
.bottoneEmoji {
  display: block; 
  margin-bottom: 5px; 
}
.bottoneEmoji-container {
  position: absolute; 
  top: 40px;    
  right: 0;    
  z-index: 999; 
  background-color: #333; 
  padding: 5px;
  border-radius: 5px;
}
.message-button{
	height: 40px;
    width: 70px;
    -webkit-appearance: none; 
    -moz-appearance: none; 
    appearance: none; 
    background: transparent; 
    border: none; 
    padding: 0; 
    margin: 0;
    cursor: pointer; 
    outline: none; 
}

.message-time {
    position: absolute; 
    bottom: 5px;
    right: 10px;
    font-size: 0.8em; 
    color: #888; 
}

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

    contain: content; /* non rimuovere assolutamente */
}
.message-content {
    word-break: break-word; 
    white-space: pre-wrap; 
}

body, html {
    overflow: hidden; 
    height: 100%; 
}
.btn-paperclip {
    width: 50px; 
    display: flex;
    justify-content: center;
    align-items: center;
}
.emoji-checkbox-container {
    display: flex; 
    flex-wrap: wrap; 
	
}
.emoji-item {
    margin-right: 5px; 
}
.trash-button {
    /* order: 2; */
	flex-shrink: 0;
}

.reply-container {
    background-color: #2b2b2b;
    border-radius: 5px;
    margin-bottom: 10px;
    padding: 10px;
    position: relative;
    border-left: 4px solid #0d6efd;
}

.reply-close-button {
    position: absolute;
    top: 5px;
    right: 5px;
    background: none;
    border: none;
    color: #fff;
    cursor: pointer;
}

.forwarded-by {
    font-family: 'Segoe UI', 'Roboto', sans-serif;
    font-style: italic;
    color: #8ab4f8;
    font-size: 0.9em;
    font-weight: 500;
    letter-spacing: 0.5px;
    padding: 4px 8px;
    background-color: rgba(138, 180, 248, 0.1);
    border-radius: 4px;
    display: inline-block;
    margin-top: 5px;
}
</style>