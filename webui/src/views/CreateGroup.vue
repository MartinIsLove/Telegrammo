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
            searchCompleted: false,
            if_selection: false,

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
                this.cs = sessionStorage.getItem("cs")
                if (this.chUsername === ''){
                    let tmp = '$'
                    await this.$axios.get(`/search/users/${tmp}`, {headers: {cs:this.cs}});
                }else{
                let response = await this.$axios.get(`/search/users/${this.chUsername}`, {headers: {cs:this.cs}});
                this.searchCompleted=true
                this.users = response.data
                if (response.status === 200){
                        this.message = 'usernames find';
                    }
                }
                
                
                

                // if (this.selectedUser.length > 0){
                //     this.if_selection = true
                // }else {
                //     this.if_selection = false
                // }
            }
            catch(error){
                this.error=error
                console.error('Error fetching user data:', error);
            }
                


        },
	},
	mounted() {
		this.refresh()
	}
}
</script>
<template>
    
    <div class="row">
        <div class="col-4"></div>
            <div class="col-2 mt-2">
                <h2 class="fw-normal">usernames</h2>
                <form class="mt-5" @submit.prevent="getUsers">
                    <label for="username" class="form-label ">search user</label>
                    <input type="text" class="form-control mb-2" id="username" placeholder="insert username" v-model="chUsername" @input="getUsers">
                </form>
            <div class="col-6"></div>
            <div v-for="user in users" :key="id" >
                <div class="form-check">
                <input class="form-check-input" type="checkbox" name="selectedUser" :value="user" v-model="selectedUser">
                
                    <label class="form-check-label" :for="'user-'+user.id">
                        
                        {{ user.username }}
                    </label>
                </div>
            </div>
            
            <div v-if="searchCompleted">
                <button class="btn btn-primary mt-2" @click="handleSelectedUser">Select</button>
            </div>
            <div class="mt-3"  v-if="if_selection">
                <h5>utenti selezionati:</h5>
                <div v-for="user2 in selectedUser" :key="user2.id" class="selected-user">
                    -{{ user2.username }}
                </div>
            </div>
            <div>
                <p v-if=error>
                    {{ error }}
                </p>
                <p v-if=message>
                    {{ message }}
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