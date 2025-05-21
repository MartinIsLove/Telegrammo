<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,

			
			id: null,
            chatId: null,
			group_name: "",
			chPhoto: null,
            users: [],
		}
	},
    mounted() {
        this.id = sessionStorage.getItem("authorization")
		if (this.id == null){
			this.$router.push("/");
		}
		this.refresh()
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
            this.id = sessionStorage.getItem("authorization");
            this.chatId = parseInt(this.$route.params.id, 10);
			this.loadusers();
		},
        async loadusers(){
            let response = await this.$axios.get(`/conversation/${this.chatId}/group/users`, {headers:{authorization:this.id}});
			this.users = response.data;
            
        },
		async leavegroup(){
			await this.$axios.delete(`/conversation/leave/${this.chatId}`, {headers:{authorization:this.id}});

			this.$router.push("/home");
		},
		addPartecipants(){
			this.$router.push(`/options/group/`+this.chatId+`/partecipats/`);
		},
		async changeGroupDetails(){
			this.message=''
            this.error=''
                if (this.chPhoto === null && this.group_name !== '' ){
                    try{
                        let response = await this.$axios.put(`/conversation/group/name`,{nome_chat:this.group_name, id_chat:this.chatId}, {headers: {authorization:this.id}});
                        if (response.status === 200){
                            this.message = 'group name changed successfully';
                        }
                        

                    }
                    catch(error){
                        this.error= error
                        console.error('error change group name')
                    }

                }
                if (this.chPhoto !== null && this.group_name === ''){
                    let formData = new FormData();
                    formData.append('propic', this.chPhoto);
					formData.append('id_chat', this.chatId);
                    try{
                        
                        let response = await this.$axios.put(`/conversation/group/photo`, formData, {
                        headers: {
                            'Content-Type': 'multipart/form-data',
                            'authorization': this.id
                        }});
                        if (response.status === 200){
                            this.message = 'group photo changed successfully';
                        }
                        

                    }
                    catch(error){
                        this.error= error
                        console.error('error change group photo')
                    }
                }
                if (this.chPhoto !== null && this.group_name !== ''){
                    let formData = new FormData();
                    formData.append('propic', this.chPhoto);
					formData.append('id_chat', this.chatId);

                    try{
                        let response = await this.$axios.put(`/conversation/group/name`,{nome_chat:this.group_name, id_chat:this.chatId}, {headers: {authorization:this.id}});
                        if (response.status === 200){
                            this.message = 'group name changed successfully';
                        }
                        

                    }
                    catch(error){
                        this.error= error
                        console.error('error in change group name')
                    }
                    try{
                        let response2 = await this.$axios.put(`conversation/group/photo`, formData, {
                        headers: {
                            'Content-Type': 'multipart/form-data',
                            'authorization': this.id
                        }});
                        if (response2.status === 200){
                            this.message = 'group photo changed successfully';
                        }
                        
                    }
                    catch(error){
                        this.error= error
                        console.error('error in change group photo')
                    }
                }
				this.$router.push("/chat/"+ this.chatId);

		},
		handleFileUpload(event){
            this.chPhoto=event.target.files[0]
        },
	},
	
}
</script>
<template>
	<div class="mt-2">
		<div>
		<button @click="leavegroup" class="btn btn-secondary">
			-leave group
		</button>
		<button @click="addPartecipants" class="btn btn-secondary ms-2">
			+Add Partecipants
		</button>
		</div>
		<form @submit.prevent="changeGroupDetails">
                <div class="mb-3">
                    <label for="group_name" class="form-label">Group Name</label>
                    <input type="text" class="form-control" id="group_name" aria-describedby="inserisci username" v-model="group_name">
					<input type="file" class="form-control mb-2 mt-2" id="photo" @change="handleFileUpload">
                </div>
                <button type="submit" class="btn btn-primary">Submit</button>
        </form>
        <h1>membri del gruppo:</h1>
        <div class="custom-scroll mt-4">
            <div v-for="user in users" :key="user.id" class="row g-0 border rounded my-2 p-2 align-items-center">
                <div class="col-md-1 col-2">
                    <img v-if="user.propic" :src="'data:image/png;base64,' + user.propic" class="rounded-circle me-2" width="70" height="70" role="img" focusable="false" style="object-fit: cover;">
                </div>
                <div class="col-md-11 col-10">
                    <div class="card-body">
                        <h5 class="fw-normal">{{ user.username }}</h5>
                    </div>
                </div>
            </div>
        </div>
	</div>
</template>