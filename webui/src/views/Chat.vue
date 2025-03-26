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
			
			
			await this.$axios.post("/message", formData, {headers: {'Content-Type': 'multipart/form-data', cs: this.id}});
			this.message = ''
			this.photo = null 
            this.$refs.photoInput.value = '' 
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
		handleFileUpload(event){
			this.photo = event.target.files[0];
		},
		async handleCheckboxChange(mes, emojiItem) {
            console.log(mes, emojiItem);
			try{
				if (emojiItem !== ""){
					await this.$axios.post(`/message/comment`, {id_chat:this.chatId, emoji:emojiItem, id_mes:mes.id_mess},{headers:{cs:this.id}});
				}
				this.fetchMessageNoBottom()
				
			}
			catch(e){
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
		async handleCheckboxChangeUncomment(mes){
			try{
				await this.$axios.post(`/message/uncomment`, {id_chat:this.chatId, id_mes:mes.id_mess},{headers:{cs:this.id}});
				this.fetchMessageNoBottom()
			}
			catch(e){
				console.error(e);
			}
		},
		async forwardMessage(id_mes, id_utente){
			this.$router.push("/chat/"+this.chatId+"/forward/"+id_mes+"/"+id_utente );

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
			
			<button class="btn btn-secondary mt-2 rounded text-tart"  @click="optionsButtonHandler">
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
			

			<div v-else-if="id==mes.id_mitt && mes.forward_id==-1" class="d-flex justify-content-end my-2 ">
				<div class=" message-container row g-0 border rounded p-2 position-relative ">
					<div class="d-flex justify-content-between w-100">
						<p v-if="mes.gruppo" class="text-wrap message-content">
							<span class="text-capitalize fw-bold">{{ mes.nome }}</span><br>
							<div v-if="mes.photo !== null">
                                <img :src="'data:image/png;base64,' + mes.photo" class="img-fluid" />
							</div>
							
							{{ mes.testo }}
							<div>
								<button class="forward-button btn btn-secondary" @click="forwardMessage(mes.id_mess, mes.id_mitt)">
									<i class="bi bi-arrow-90deg-right"></i>
								</button>
							</div>
						</p> 
						<p v-else class="message-content">
							<div v-if="mes.photo !== null">
                                <img :src="'data:image/png;base64,' + mes.photo" class="img-fluid" />
							</div>
							<span>{{ mes.testo }}</span>
							<div>
								<button class="forward-button btn btn-secondary" @click="forwardMessage(mes.id_mess, mes.id_mitt)">
									<i class="bi bi-arrow-90deg-right"></i>
								</button>
							</div>
						</p>
						
					</div>
					<div class="emoji-checkbox-container " >
						<div v-for="emojiItem in emoji" :key="emojiItem" :value="emojiItem" class="emoji-item"> 
							
							<div v-if="emojiItem === mes.myEmoji">
								<input type="checkbox" class="btn-check" :id="'btn-check-' + mes.id_mess + '-' + emojiItem" autocomplete="off" @change="handleCheckboxChangeUncomment(mes)">
								<label class="btn btn-primary mb-2" :for="'btn-check-' + mes.id_mess + '-' + emojiItem">{{ emojiItem }} {{mes.emoji[getEmojiCount(emojiItem)]}}</label>	
							</div>		
							<div v-else>
								<input type="checkbox" class="btn-check" :id="'btn-check-' + mes.id_mess + '-' + emojiItem" autocomplete="off" @change="handleCheckboxChange(mes, emojiItem)">
								<label class="btn btn-secondary mb-2" :for="'btn-check-' + mes.id_mess + '-' + emojiItem">{{ emojiItem }} {{mes.emoji[getEmojiCount(emojiItem)]}}</label>	
							</div>		
						</div>
					</div>
					
					<div class="message-time text-end">
                        {{ formatTime(mes.data) }}
                    </div>
				</div>
			</div>
			<div v-else-if="id!==mes.id_mit && mes.forward_id==-1" class="row g-0 border rounded my-2 p-2 message-container position-relative">
				<div class="d-flex justify-content-between w-100">
					<p v-if="mes.gruppo" class="text-wrap message-content">
						<span class="text-capitalize fw-bold">{{ mes.nome }}</span><br>
						<div v-if="mes.photo !== null">
                                <img :src="'data:image/png;base64,' + mes.photo" class="img-fluid" />
						</div>
						{{ mes.testo }}
						<div>
							<button class="forward-button btn btn-secondary" @click="forwardMessage(mes.id_mess, mes.id_mitt)">
								<i class="bi bi-arrow-90deg-right"></i>
							</button>
						</div>
					</p>
					<p v-else class="message-content">
						<div v-if="mes.photo !== null">
                                <img :src="'data:image/png;base64,' + mes.photo" class="img-fluid" />
						</div>
						<span>{{ mes.testo }}</span>
						<div>
							<button class="forward-button btn btn-secondary" @click="forwardMessage(mes.id_mess, mes.id_mitt)">
								<i class="bi bi-arrow-90deg-right"></i>
							</button>
						</div>
					</p>
					
					<div class="message-time">
                        {{ formatTime(mes.data) }}
                    </div>
					
				</div>
				<div class="emoji-checkbox-container " >
					<div v-for="emojiItem in emoji" :key="emojiItem" :value="emojiItem" class="emoji-item"> 
						
						<div v-if="emojiItem === mes.myEmoji">
							<input type="checkbox" class="btn-check" :id="'btn-check-' + mes.id_mess + '-' + emojiItem" autocomplete="off" @change="handleCheckboxChangeUncomment(mes)">
							<label class="btn btn-primary mb-2" :for="'btn-check-' + mes.id_mess + '-' + emojiItem">{{ emojiItem }} {{mes.emoji[getEmojiCount(emojiItem)]}}</label>	
						</div>		
						<div v-else>
							<input type="checkbox" class="btn-check" :id="'btn-check-' + mes.id_mess + '-' + emojiItem" autocomplete="off" @change="handleCheckboxChange(mes, emojiItem)">
							<label class="btn btn-secondary mb-2" :for="'btn-check-' + mes.id_mess + '-' + emojiItem">{{ emojiItem }} {{mes.emoji[getEmojiCount(emojiItem)]}}</label>	
						</div>		
					</div>
					
				</div>
				
			</div>


			<div v-else-if="id!=mes.forward_id && mes.forward_id!=-1 && mes.forward_id_mit==id" class="d-flex justify-content-end my-2 ">
				<div class=" message-container row g-0 border rounded p-2 position-relative ">
					<div class="d-flex justify-content-between w-100">
						<p v-if="mes.gruppo" class="text-wrap message-content">
							<span class="text-capitalize fw-bold">{{ mes.forward_user_mit }}</span><br>
							<div v-if="mes.photo !== null">
                                <img :src="'data:image/png;base64,' + mes.photo" class="img-fluid" />
							</div>
							
							{{ mes.testo }}
							<div>
								<button class="forward-button btn btn-secondary" @click="forwardMessage(mes.id_mess, mes.id_mitt)">
									<i class="bi bi-arrow-90deg-right"></i>
								</button>
							</div>
						</p> 
						<p v-else class="message-content">
							<div v-if="mes.photo !== null">
                                <img :src="'data:image/png;base64,' + mes.photo" class="img-fluid" />
							</div>
							<span>{{ mes.testo }}</span>
							<div>
								<button class="forward-button btn btn-secondary" @click="forwardMessage(mes.id_mess, mes.id_mitt)">
									<i class="bi bi-arrow-90deg-right"></i>
								</button>
							</div>
						</p>
						
					</div>
					<div class="emoji-checkbox-container " >
						<div v-for="emojiItem in emoji" :key="emojiItem" :value="emojiItem" class="emoji-item"> 
							
							<div v-if="emojiItem === mes.myEmoji">
								<input type="checkbox" class="btn-check" :id="'btn-check-' + mes.id_mess + '-' + emojiItem" autocomplete="off" @change="handleCheckboxChangeUncomment(mes)">
								<label class="btn btn-primary mb-2" :for="'btn-check-' + mes.id_mess + '-' + emojiItem">{{ emojiItem }} {{mes.emoji[getEmojiCount(emojiItem)]}}</label>	
							</div>		
							<div v-else>
								<input type="checkbox" class="btn-check" :id="'btn-check-' + mes.id_mess + '-' + emojiItem" autocomplete="off" @change="handleCheckboxChange(mes, emojiItem)">
								<label class="btn btn-secondary mb-2" :for="'btn-check-' + mes.id_mess + '-' + emojiItem">{{ emojiItem }} {{mes.emoji[getEmojiCount(emojiItem)]}}</label>	
							</div>		
						</div>
					</div>
					<div v-if="mes.forward_id !== -1">
						<p>
							forwarded by: {{ mes.forward_username }}
						</p>
					</div>
					<div class="message-time text-end">
                        {{ formatTime(mes.data) }}
                    </div>
				</div>
			</div>

			<div v-else-if="id==mes.forward_id && mes.forward_id!=-1 && mes.forward_id_mit==id" class="d-flex justify-content-end my-2 ">
				<div class=" message-container row g-0 border rounded p-2 position-relative ">
					<div class="d-flex justify-content-between w-100">
						<p v-if="mes.gruppo" class="text-wrap message-content">
							<span class="text-capitalize fw-bold">{{ mes.forward_user_mit }}</span><br>
							<div v-if="mes.photo !== null">
                                <img :src="'data:image/png;base64,' + mes.photo" class="img-fluid" />
							</div>
							
							{{ mes.testo }}
							<div>
								<button class="forward-button btn btn-secondary" @click="forwardMessage(mes.id_mess, mes.id_mitt)">
									<i class="bi bi-arrow-90deg-right"></i>
								</button>
							</div>
						</p> 
						<p v-else class="message-content">
							<div v-if="mes.photo !== null">
                                <img :src="'data:image/png;base64,' + mes.photo" class="img-fluid" />
							</div>
							<span>{{ mes.testo }}</span>
							<div>
								<button class="forward-button btn btn-secondary" @click="forwardMessage(mes.id_mess, mes.id_mitt)">
									<i class="bi bi-arrow-90deg-right"></i>
								</button>
							</div>
						</p>
						
					</div>
					<div class="emoji-checkbox-container " >
						<div v-for="emojiItem in emoji" :key="emojiItem" :value="emojiItem" class="emoji-item"> 
							
							<div v-if="emojiItem === mes.myEmoji">
								<input type="checkbox" class="btn-check" :id="'btn-check-' + mes.id_mess + '-' + emojiItem" autocomplete="off" @change="handleCheckboxChangeUncomment(mes)">
								<label class="btn btn-primary mb-2" :for="'btn-check-' + mes.id_mess + '-' + emojiItem">{{ emojiItem }} {{mes.emoji[getEmojiCount(emojiItem)]}}</label>	
							</div>		
							<div v-else>
								<input type="checkbox" class="btn-check" :id="'btn-check-' + mes.id_mess + '-' + emojiItem" autocomplete="off" @change="handleCheckboxChange(mes, emojiItem)">
								<label class="btn btn-secondary mb-2" :for="'btn-check-' + mes.id_mess + '-' + emojiItem">{{ emojiItem }} {{mes.emoji[getEmojiCount(emojiItem)]}}</label>	
							</div>		
						</div>
					</div>
					<div v-if="mes.forward_id !== -1">
						<p>
							forwarded by: {{ mes.forward_username }}
						</p>
					</div>
					<div class="message-time text-end">
                        {{ formatTime(mes.data) }}
                    </div>
				</div>
			</div>



			<div v-else-if="id==mes.forward_id && mes.forward_id!=-1 && mes.forward_id_mit!=id" class="row g-0 border rounded my-2 p-2 message-container position-relative">
				<div class="d-flex justify-content-between w-100">
					<p v-if="mes.gruppo" class="text-wrap message-content">
						<span class="text-capitalize fw-bold">{{ mes.forward_user_mit }}</span><br>
						<div v-if="mes.photo !== null">
                                <img :src="'data:image/png;base64,' + mes.photo" class="img-fluid" />
						</div>
						{{ mes.testo }}
						<div>
							<button class="forward-button btn btn-secondary" @click="forwardMessage(mes.id_mess, mes.id_mitt)">
								<i class="bi bi-arrow-90deg-right"></i>
							</button>
						</div>
					</p>
					<p v-else class="message-content">
						<div v-if="mes.photo !== null">
                                <img :src="'data:image/png;base64,' + mes.photo" class="img-fluid" />
						</div>
						<span>{{ mes.testo }}</span>
						<div>
							<button class="forward-button btn btn-secondary" @click="forwardMessage(mes.id_mess, mes.id_mitt)">
								<i class="bi bi-arrow-90deg-right"></i>
							</button>
						</div>
					</p>
					
					<div class="message-time">
                        {{ formatTime(mes.data) }}
                    </div>
					
				</div>
				<div class="emoji-checkbox-container " >
					<div v-for="emojiItem in emoji" :key="emojiItem" :value="emojiItem" class="emoji-item"> 
						
						<div v-if="emojiItem === mes.myEmoji">
							<input type="checkbox" class="btn-check" :id="'btn-check-' + mes.id_mess + '-' + emojiItem" autocomplete="off" @change="handleCheckboxChangeUncomment(mes)">
							<label class="btn btn-primary mb-2" :for="'btn-check-' + mes.id_mess + '-' + emojiItem">{{ emojiItem }} {{mes.emoji[getEmojiCount(emojiItem)]}}</label>	
						</div>		
						<div v-else>
							<input type="checkbox" class="btn-check" :id="'btn-check-' + mes.id_mess + '-' + emojiItem" autocomplete="off" @change="handleCheckboxChange(mes, emojiItem)">
							<label class="btn btn-secondary mb-2" :for="'btn-check-' + mes.id_mess + '-' + emojiItem">{{ emojiItem }} {{mes.emoji[getEmojiCount(emojiItem)]}}</label>	
						</div>		
					</div>
					
				</div>
				<div v-if="mes.forward_id !== -1">
					<p>
						forwarded by: {{ mes.forward_username }}
					</p>
				</div>
			</div>

			<div v-else-if="id!=mes.forward_id && mes.forward_id!=-1 && mes.forward_id_mit!=id" class="row g-0 border rounded my-2 p-2 message-container position-relative">
				<div class="d-flex justify-content-between w-100">
					<p v-if="mes.gruppo" class="text-wrap message-content">
						<span class="text-capitalize fw-bold">{{ mes.forward_user_mit}}</span><br>
						<div v-if="mes.photo !== null">
                                <img :src="'data:image/png;base64,' + mes.photo" class="img-fluid" />
						</div>
						{{ mes.testo }}
						<div>
							<button class="forward-button btn btn-secondary" @click="forwardMessage(mes.id_mess, mes.id_mitt)">
								<i class="bi bi-arrow-90deg-right"></i>
							</button>
						</div>
					</p>
					<p v-else class="message-content">
						<div v-if="mes.photo !== null">
                                <img :src="'data:image/png;base64,' + mes.photo" class="img-fluid" />
						</div>
						<span>{{ mes.testo }}</span>
						<div>
							<button class="forward-button btn btn-secondary" @click="forwardMessage(mes.id_mess, mes.id_mitt)">
								<i class="bi bi-arrow-90deg-right"></i>
							</button>
						</div>
					</p>
					
					<div class="message-time">
                        {{ formatTime(mes.data) }}
                    </div>
					
				</div>
				<div class="emoji-checkbox-container " >
					<div v-for="emojiItem in emoji" :key="emojiItem" :value="emojiItem" class="emoji-item"> 
						
						<div v-if="emojiItem === mes.myEmoji">
							<input type="checkbox" class="btn-check" :id="'btn-check-' + mes.id_mess + '-' + emojiItem" autocomplete="off" @change="handleCheckboxChangeUncomment(mes)">
							<label class="btn btn-primary mb-2" :for="'btn-check-' + mes.id_mess + '-' + emojiItem">{{ emojiItem }} {{mes.emoji[getEmojiCount(emojiItem)]}}</label>	
						</div>		
						<div v-else>
							<input type="checkbox" class="btn-check" :id="'btn-check-' + mes.id_mess + '-' + emojiItem" autocomplete="off" @change="handleCheckboxChange(mes, emojiItem)">
							<label class="btn btn-secondary mb-2" :for="'btn-check-' + mes.id_mess + '-' + emojiItem">{{ emojiItem }} {{mes.emoji[getEmojiCount(emojiItem)]}}</label>	
						</div>		
					</div>
					
				</div>
				<div v-if="mes.forward_id !== -1">
					<p>
						forwarded by: {{ mes.forward_username }}
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
		<label for="photo" class="btn btn-secondary mt-2 rounded ms-2 btn-paperclip">
            <i class="bi bi-paperclip"></i>
            <input type="file" class="d-none" id="photo" @change="handleFileUpload" ref="photoInput">
        </label>
    <br>
    </div>
</div>
</form>
</template>
<style>
.forward-button {
    position: absolute;
    top: 10px;
    right: 10px;
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
</style>