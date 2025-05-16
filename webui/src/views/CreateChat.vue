<script>
export default {
    data: function() {
        return {
            errormsg: null,
            loading: false,
            some_data: null,

            cs: '',
            chUsername: '',
            users: [],
            selectedUser: null, 
            searchCompleted: false,
            chat_id: null,
            
            error: '',
            message: '',
        }
    },
    mounted() {
        this.cs = sessionStorage.getItem("cs")
		if (this.cs == null){
			this.$router.push("/");
		}
    },
    methods: {
        async getUsers() {
            try{
                this.message=''
                this.error=''
                this.users=[]
                let response = await this.$axios.get(`/search/users/${this.chUsername}`, {headers: {cs:this.cs}});
                if (response.status === 200){
                        this.message = 'usernames find';
                    }
                
                this.searchCompleted=true
                this.users = response.data
            }
            catch(error){
                
                this.message = 'no users found';
                
                this.error=error
                console.error('Error fetching user data:', error);
            }
        },
        async handleSelectedUser(){
            try{
                this.cs = sessionStorage.getItem("cs")
                let response = await this.$axios.post(`/conversation`,{id: this.selectedUser.id} ,{headers: {cs:this.cs}});
                if (response.status === 200){
                        this.message = 'chat created';
                        this.chat_id=response.data
                        this.$router.push('/chat/'+ this.chat_id.id_chat)
                    }
                
                
            }
            catch(error){
                this.error=error
                
                console.error('error create chat user data: ', error)
            }

        }
    },
    
}
</script>

<template>
    <div class="d-flex justify-content-center align-items-center vh-100">
        <div class="col-6">
            <h2 class="fw-normal">Create Chat</h2>
            <form class="mt-5" @submit.prevent="getUsers">
                <label for="username" class="form-label">Search User</label>
                <input type="text" class="form-control mb-2" id="username" placeholder="Insert Username" v-model="chUsername" @input="getUsers">
            </form>
            <div v-for="user in users" :key="user.id" class="form-check">
                <input class="form-check-input" type="radio" name="selectedUser" :value="user" v-model="selectedUser">
                <label class="form-check-label" :for="'user-'+user.id">
                    {{ user.username }}
                </label>
            </div>
            <div v-if="searchCompleted">
                <button class="btn btn-primary mt-2" @click="handleSelectedUser">Select</button>
            </div>
            <div>
                <p v-if="message">
                    {{ message }}
                </p>
                <p v-if="error">
                    {{ error.response.data }}
                </p>
            </div>
        </div>
    </div>
</template>

<style scoped>
.selected-user {
    margin: 5px 0;
    padding: 5px;
    border: 1px solid #ccc;
    border-radius: 5px;
}
</style>