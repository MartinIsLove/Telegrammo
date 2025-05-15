<script>
export default {
    data: function() {
        return {
            errormsg: null,
            loading: false,
            some_data: null,

            id: null,
            chatId: null,

            chUsername: '',
            users: [],
            selectedUser: [], 
            if_selection: false,
            group_id: null,
            group_name: '',
            group_photo: null,

            error: '',
            message: '',
        }
    },
    watch: {
        selectedUser(newVal) {
            this.if_selection = newVal.length > 0;
        }
    },
    methods: {
        async refresh() {
            this.loading = true;
            this.errormsg = null;
            this.id = sessionStorage.getItem("cs")
            this.chatId = parseInt(this.$route.params.id, 10);
        },
        async getUsers() {
            try{
                this.message=''
                this.error=''
                this.users=[]
                this.id = sessionStorage.getItem("cs")
                if (this.chUsername === ''){
                    let tmp = '$'
                    await this.$axios.get(`/search/users/${tmp}`, {headers: {cs:this.id}});
                }else{
                let response = await this.$axios.get(`/search/users/${this.chUsername}`, {headers: {cs:this.id}});
                this.searchCompleted=true
                this.users = response.data
                if (response.status === 200){
                        this.message = 'usernames find';
                    }
                }
            }
            catch(error){
                this.error=error
                console.error('Error fetching user data:', error);
            }
        },
        handleFileUpload(event){
            this.group_photo = event.target.files[0];
        },
        async handleSelectedUser(){
            try{
                let formData = new FormData();
                formData.append('propic', this.group_photo);
                formData.append('nome_chat', this.group_name);
                this.selectedUser.forEach(user => {
                    formData.append('membri', user.id.toString());
                });

                this.id = sessionStorage.getItem("cs")
                let response = await this.$axios.post(`/conversation/group`,formData ,{headers: {'Content-Type': 'multipart/form-data', cs:this.id}});
                if (response.status === 200){
                        this.message = 'chat created';
                        this.group_id = response.data
                        this.$router.push('/chat/'+ this.group_id.id_group )
                }
            }
            catch(error){
                this.error=error
                console.error('error create group user data: ', error)
            }
        },
    },
    mounted() {
        this.refresh()
    }
}
</script>

<template>
    <div class="d-flex justify-content-center align-items-center vh-100">
        <div class="col-6">
            <form @submit.prevent="handleSelectedUser">
                <div class="mb-3">
                    <label for="group_name" class="form-label">Group Name</label>
                    <input type="text" class="form-control" id="group_name" aria-describedby="inserisci username" v-model="group_name">
                </div>
                <button type="submit" class="btn btn-primary">Submit</button>
            </form>
            <label for="photo" class="form-label mt-3">Change Photo</label>
            <input type="file" class="form-control mb-2" id="photo" @change="handleFileUpload">
            <div class="mt-4">
                <h2 class="fw-normal">Usernames</h2>
                <form class="mt-3" @submit.prevent="getUsers">
                    <label for="username" class="form-label">Search User</label>
                    <input type="text" class="form-control mb-2" id="username" placeholder="Insert Username" v-model="chUsername" @input="getUsers">
                </form>
                <div v-for="user in users" :key="user.id" class="form-check">
                    <input class="form-check-input" type="checkbox" name="selectedUser" :value="user" v-model="selectedUser">
                    <label class="form-check-label" :for="'user-'+user.id">
                        {{ user.username }}
                    </label>
                </div>
                <div class="mt-3" v-if="if_selection">
                    <h5>Selected Users:</h5>
                    <div v-for="user2 in selectedUser" :key="user2.id" class="selected-user">
                        -{{ user2.username }}
                    </div>
                </div>
                <div>
                    <p v-if="error">
                        {{ error.response.data }}
                    </p>
                    <p v-if="message">
                        {{ message }}
                    </p>
                </div>
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