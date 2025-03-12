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
		}
	},
	methods: {
		async getUsers() {
			try{
                this.cs = sessionStorage.getItem("cs")
                console.log(this.chUsername)
                let response = await this.$axios.get(`/user/${this.chUsername}`, {headers: {cs:this.cs}});
                this.users = response.data
            }
            catch(error){
                this.error=error
                console.error('Error fetching user data:', error);
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
            <div v-for="user in users" :key="id">
                {{ user.username }}
            </div>
        </div>
    </div>
</template>