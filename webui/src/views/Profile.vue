<script>
export default {
	data: function() {
		return {
            error: '',
            message: '',

            id: null,
            username: '',
            propic: null,

            chUsername: '',
            chPhoto: null
        }
	},
	methods:{
		profileButtonHandler() {
			this.$router.push("/profile");
		},
        async fetchUserData(){
            try{
                this.id = sessionStorage.getItem("cs")
                let response = await this.$axios.get(`/users/${this.id}`, {headers: {cs:this.id}});
                this.username = response.data.username;
                this.propic = 'data:image/png;base64,' + response.data.propic;
            }
            catch(error){
                this.error=error
                console.error('Error fetching user data:', error);
            }
        },
        handleFileUpload(event){
            this.chPhoto=event.target.files[0]
        },
        async changeUserDetails(){
            if (this.chPhoto === null && this.chUsername !== '' ){
                try{
                    let response = await this.$axios.put(`/user/name`,{username:this.chUsername}, {headers: {cs:this.id}});
                    if (response.status === 200){
                        this.message = 'Username changed successfully';
                    }
                    

                }
                catch(error){
                    this.error= error
                    console.error('error change username')
                }

            }
            if (this.chPhoto !== null && this.chUsername === ''){
                let formData = new FormData();
                formData.append('photo', this.chPhoto);
                try{
                    
                    let response = await this.$axios.put(`/user/photo`, formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                        'cs': this.id
                    }});
                    if (response.status === 200){
                        this.message = 'photo changed successfully';
                    }
                    

                }
                catch(error){
                    this.error= error
                    console.error('error change photo')
                }
            }
            if (this.chPhoto !== null && this.chUsername !== ''){
                let formData = new FormData();
                formData.append('photo', this.chPhoto);
                try{
                    let response = await this.$axios.put(`/user/name`,{username:this.chUsername}, {headers: {cs:this.id}});
                    if (response.status === 200){
                        this.message = 'Username changed successfully';
                    }
                    

                }
                catch(error){
                    this.error= error
                    console.error('error in change username')
                }
                try{
                    let response2 = await this.$axios.put(`/user/photo`, formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                        'cs': this.id
                    }});
                    if (response2.status === 200){
                        this.message = 'photo changed successfully';
                    }
                    
                }
                catch(error){
                    this.error= error
                    console.error('error in change user photo')
                }
            }
            this.fetchUserData()
        }
	},
    async mounted() {
        this.fetchUserData();
    }
}
</script>
<template>
    <div class="row">
        <div class="col-4"></div>
        <div class="col-2 mt-2">
            <img :src="propic" class="bd-placeholder-img rounded-circle" width="140" height="140" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="Placeholder" preserveAspectRatio="xMidYMid slice" focusable="false"><title>Placeholder</title><rect width="100%" height="100%" fill="var(--bs-secondary-color)"></rect>
            <h2 class="fw-normal">{{username}}</h2>
            <form class="mt-5" @submit.prevent="changeUserDetails">
                <label for="username" class="form-label ">New Username</label>
                <input type="text" class="form-control mb-2" id="username" placeholder="insert new username" v-model="chUsername">
                <label for="photo" class="form-label" >Change Photo</label>
                <input type="file" class="form-control mb-2" id="photo" @change="handleFileUpload">
                <button class="btn btn-secondary mt-2"  >Change</button>
            </form>
            <p v-if=error>
                {{ error }}
            </p>
            <p v-if=message>
                {{ messagge }}
            </p>
            <div class="col-6"></div>
        </div>
    </div>
        
</template>

