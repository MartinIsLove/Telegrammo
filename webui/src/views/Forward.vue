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
            
            searchCompleted: false,
            if_selection: false,
            selectedChat: [], 
            chats: [],
            chChat: '',

            error: '',
            message: '',
        }
    },
    watch: {
        selectedChat(newVal) {
            this.if_selection = newVal.length > 0;
        },
        chChat(newVal) {
            if (newVal === '') {
                this.searchCompleted = false;
                this.refresh(); 
            } else {
                this.getChats(); 
            }
        }
    },
    methods: {
        async refresh() {
            this.loading = true;
            this.errormsg = null;
            this.id = sessionStorage.getItem("cs")
            try {
                let response = await this.$axios.get("/conversations", { headers: { cs: this.id } });
                this.chats = response.data;
            } catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
        },
        async getChats() {
            try {
                this.message = ''
                this.error = ''
                this.users = []
                this.id = sessionStorage.getItem("cs")
                if (this.chChat === "") {
                    this.searchCompleted = false
                } else {
                    let response = await this.$axios.get(`/search/chat/${this.chChat}`, { headers: { cs: this.id } });
                    this.searchCompleted = true
                    this.chats = response.data
                    if (response.status === 200) {
                        this.message = 'usernames find';
                    }
                }
            } catch (error) {
                this.error = error
                console.error('Error fetching user data:', error);
            }
        },
        async handleSelectedUser() {
            try {
                let formData = new FormData();
                formData.append('propic', this.group_photo);
                formData.append('nome_chat', this.group_name);
                formData.append('membri', JSON.stringify(this.selectedUser.map(user => user.id)));

                this.id = sessionStorage.getItem("cs")
                let response = await this.$axios.post(`/conversation/group`, formData, { headers: { 'Content-Type': 'multipart/form-data', cs: this.id } });
                if (response.status === 200) {
                    this.message = 'chat created';
                    this.group_id = response.data
                    this.$router.push('/chat/' + this.group_id.id_group)
                }
            } catch (error) {
                this.error = error
                console.error('error create group user data: ', error)
            }
        },
        isSelected(chat) {
            return this.selectedChat.some(selected => selected.id_chat === chat.id_chat);
        },
        remove(chat){
            const index = this.selectedChat.findIndex(selected => selected.id_chat === chat.id_chat);
            if (index !== -1) {
                this.selectedChat.splice(index, 1);
            }
        }
    },
    mounted() {
        this.refresh()
    }
}
</script>

<template>
    <div style="position: sticky; top: 0; z-index: 1000; padding: 10px; border-bottom: 1px solid #ccc;">
        <div class="d-flex">
            <form class="mt-3 flex-grow-1" @submit.prevent="getChats">
                <label for="username" class="form-label">Search Chat</label>
                <input type="text" class="form-control mb-2" id="username" placeholder="Insert Username" v-model="chChat" @input="getChats">
            </form>
            <button class="btn btn-secondary ms-1 mt-5 ml-auto h-100">
                Submit
            </button>
        </div>
    </div>
    <div class="mt-3" v-if="if_selection">
        <h5>Selected Chats:</h5>
        <div v-for="chat2 in selectedChat" :key="chat2.id_chat" class="selected-user">
            <div class="row g-0 border rounded my-2 p-2">
                <div class="col-md-1 col-2">
                    <img v-if="chat2.propic" :src="'data:image/png;base64,'+ chat2.propic" class="rounded-circle me-2" width="70" height="70" role="img" focusable="false" style="object-fit: cover;">
                </div>
                <div class="col-md-11 col-10">
                    <div class="card-body">
                        <h5>{{chat2.nome}}</h5>
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
                            <h5>{{chat.nome}}</h5>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <div v-else>
        <div class="custom-scroll">
            <div v-for="chat in chats.filter(chat => !isSelected(chat))" :key="chat.id_chat" class="form-check">
                <input class="form-check-input" type="checkbox" name="selectedUser" :value="chat" v-model="selectedChat">
                <div class="row g-0 border rounded my-2 p-2">
                    <div class="col-md-1 col-2">
                        <img v-if="chat.propic" :src="'data:image/png;base64,'+ chat.propic" class="rounded-circle me-2" width="70" height="70" role="img" focusable="false" style="object-fit: cover;">
                    </div>
                    <div class="col-md-11 col-10">
                        <div class="card-body">
                            <h5>{{chat.nome}}</h5>
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