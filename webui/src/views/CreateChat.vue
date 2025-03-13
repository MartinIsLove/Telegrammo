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

            error: '',
            message: '',
		}
	},
	methods: {
		async getUsers() {
			try{
                this.message=''
                this.error=''
                this.users=[]
                this.cs = sessionStorage.getItem("cs")
                let response = await this.$axios.get(`/search/users/${this.chUsername}`, {headers: {cs:this.cs}});
                if (response.status === 200){
                        this.message = 'usernames find';
                    }
                this.searchCompleted=true
                this.users = response.data
            }
            catch(error){
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
                        this.$router.push("/profile")
                    }
            }
            catch(error){
                this.error=error
                console.error('error create chat user data: ', error)
            }
        }
	},
	mounted() {
		
	}
}


</script>
<template>
    <div class="row">
        <div class="col-4"></div>
            <div class="col-2 mt-2">
                <h2 class="fw-normal">username</h2>
                <form class="mt-5" @submit.prevent="getUsers">
                    <label for="username" class="form-label ">search user</label>
                    <input type="text" class="form-control mb-2" id="username" placeholder="insert username" v-model="chUsername">
                    <button class="btn btn-secondary mt-2"  >find</button>
                </form>
            <div class="col-6"></div>
            <div v-for="user in users" :key="id" >
                <div class="form-check">
                <input class="form-check-input" type="radio" name="selectedUser" :value="user" v-model="selectedUser">
                    <label class="form-check-label" :for="'user-'+user.id">
                        {{ user.username }}
                    </label>
                </div>
            </div>
            <div v-if="searchCompleted">
                <button class="btn btn-primary mt-2" @click="handleSelectedUser">Select</button>
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