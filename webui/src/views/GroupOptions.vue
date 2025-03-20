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
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
            this.id = sessionStorage.getItem("cs")
            this.chatId = parseInt(this.$route.params.id, 10);
			
			
		},
		async leavegroup(){
			await this.$axios.delete(`/conversation/leave/${this.chatId}`, {headers:{cs:this.id}});

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
                        let response = await this.$axios.put(`/conversation/group/name`,{nome_chat:this.group_name, id_chat:this.chatId}, {headers: {cs:this.id}});
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
                            'cs': this.id
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
                        let response = await this.$axios.put(`/conversation/group/name`,{nome_chat:this.group_name, id_chat:this.chatId}, {headers: {cs:this.id}});
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
                            'cs': this.id
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
	mounted() {
		this.refresh()
	}
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

	</div>
</template>