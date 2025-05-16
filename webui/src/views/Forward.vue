<script>
    export default {
    data: function() {
        return {
            errormsg: null,
            loading: false,
            some_data: null,

            id: null,
            chatId: null,

            users: [],
            group_id: null,
            group_name: '',
            group_photo: null,
            
            id_chat:null,
            id_utente:null,
            id_mes: null,
            searchCompleted: false,
            if_selection: false,
            selectedChat: [], 
            selectedUser: [],


            allselectedchat: [],
            chats: [],
            chChat: '',

            error: '',
            message: '',
        }
    },
    watch: {
        selectedChat(newVal) {
            this.if_selection = newVal.length > 0 || this.selectedUser.length > 0;
        },
        selectedUser(newVal) {
            this.if_selection = newVal.length > 0 || this.selectedChat.length > 0;
        },
        chChat(newVal) {
            if (newVal === '') {
                this.searchCompleted = false;
                this.refresh(); 
            } else {
                this.searchCompleted = true;
            }
        }
    },
    mounted() {
        this.id = sessionStorage.getItem("cs")
		if (this.id == null){
			this.$router.push("/");
		}
        this.refresh()
    },
    methods: {
        async refresh() {
            this.loading = true;
            this.errormsg = null;
            this.id = sessionStorage.getItem("cs")
            this.id_chat = parseInt(this.$route.params.idchat, 10);
            this.id_mes = parseInt(this.$route.params.idmes, 10);
            this.id_utente= parseInt(this.$route.params.idutente, 10);
            try {
                let response = await this.$axios.get("/conversations/forward", { headers: { cs: this.id } });
                this.chats = response.data;
            } catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
        },
        

        async handleSelectedChats() {
            try {
                
                this.id = sessionStorage.getItem("cs")
                
                const selectedChatIds = this.selectedChat.map(chat => chat.id_chat);
                const selectedUserIds = this.selectedUser.map(user => user.id_utenti);
                let response = await this.$axios.post(`/message/forward`,{id_mes:this.id_mes, id_chat:selectedChatIds, id_utenti: selectedUserIds, id_for:this.id_utente} , { headers: { cs: this.id } });
                
                if (response.status === 204) {
                    this.message = 'chat created';
                    this.group_id = response.data
                    this.$router.push('/chat/'+this.id_chat)
                }
            } catch (error) {
                this.error = error
                console.error('error create group user data: ', error)
            }
        },
        isSelected(chat) {
            // Se è un gruppo/chat
            if (chat.id_chat) {
                return this.selectedChat.some(selected => selected.id_chat === chat.id_chat);
            }
            // Se è un utente
            if (chat.id_utenti) {
                return this.selectedUser.some(selected => selected.id_utenti === chat.id_utenti);
            }
            return false;
        },
        remove(chat) {
            if (chat.id_chat) {
                const index = this.selectedChat.findIndex(selected => selected.id_chat === chat.id_chat);
                if (index !== -1) {
                    this.selectedChat.splice(index, 1);
                }
            } else if (chat.id_utenti) {
                const index = this.selectedUser.findIndex(selected => selected.id_utenti === chat.id_utenti);
                if (index !== -1) {
                    this.selectedUser.splice(index, 1);
                }
            }
        },
    },
    
}
</script>

<template>
    <div style="position: sticky; top: 0; z-index: 1000; padding: 10px; border-bottom: 1px solid #ccc;">
        <div class="d-flex">
            <form class="mt-3 flex-grow-1" @submit.prevent="getChats">
                <label for="username" class="form-label">Search Chat</label>
                <input type="text" class="form-control mb-2" id="username" placeholder="Insert Username" v-model="chChat" >
            </form>
            <button class="btn btn-secondary ms-1 mt-5 ml-auto h-100" @click="handleSelectedChats">
                Submit
            </button>
        </div>
    </div>
    <div class="mt-3" v-if="if_selection">
        <h5>Selected Chats:</h5>
        <div v-for="chat2 in [...selectedChat, ...selectedUser]" :key="chat2.id_chat || chat2.id_utenti" class="selected-user">
            <div class="row g-0 border rounded my-2 p-2">
                <div class="col-md-1 col-2">
                    <img v-if="chat2.propic" :src="'data:image/png;base64,'+ chat2.propic" class="rounded-circle me-2" width="70" height="70" role="img" focusable="false" style="object-fit: cover;">
                </div>
                <div class="col-md-11 col-10">
                    <div class="card-body">
                        <h5>{{chat2.nome || chat2.username}}</h5>
                    </div>
                    <div class="col-md-1 col-2">
                        <button class="btn btn-danger" @click="remove(chat2)">Remove</button>
                    </div>
                    
                </div>
                
            </div>
        </div>
    </div>
    <div v-if="searchCompleted">
        <div class="custom-scroll">
            <div v-for="chat in chats.filter(chat => !isSelected(chat))" :key="chat.id_chat" class="form-check">
                <input class="form-check-input" type="checkbox" name="selectedUser" :value="chat" v-model="selectedChat">
                <div class="row g-0 border rounded my-2 p-2">
                    <div class="col-md-1 col-2">
                        <img v-if="chat.propic" :src="'data:image/png;base64,'+ chat.propic" class="rounded-circle me-2" width="70" height="70" role="img" focusable="false" style="object-fit: cover;">
                    </div>
                    <div class="col-md-11 col-10">
                        <div class="card-body">
                            <h5>{{chat.nome || chat.username}}</h5>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <div v-else>
        <div class="custom-scroll">
            <div v-for="chat in chats.filter(chat => !isSelected(chat))" :key="chat.id_chat" class="form-check">
                <input v-if="chat.id_chat != 0" class="form-check-input" type="checkbox" name="selectedUser" :value="chat" v-model="selectedChat">
                <input v-if="chat.id_utenti != 0" class="form-check-input" type="checkbox" name="selectedUser" :value="chat" v-model="selectedUser">
                <div  class="row g-0 border rounded my-2 p-2">
                    <div class="col-md-1 col-2">
                        <img v-if="chat.propic" :src="'data:image/png;base64,'+ chat.propic" class="rounded-circle me-2" width="70" height="70" role="img" focusable="false" style="object-fit: cover;">
                    </div>
                    <div class="col-md-11 col-10">
                        <div class="card-body">
                            <h5>{{chat.nome || chat.username}}</h5>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style>
.custom-box {
    min-height: 60px;
}
.custom-scroll {
    max-height: 95vh;
    overflow-y: auto;
}
</style>