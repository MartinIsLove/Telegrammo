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
            selectedUser: [], 
            if_selection: false,
            group_id: null,
            group_name: '',
            group_photo: null,
            
            chats: [],
            chChat: '',

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
			try {
				
				let response = await this.$axios.get("/conversations", {headers: {cs:this.id}});
				this.chats = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
        async getChats() {
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
                formData.append('membri', JSON.stringify(this.selectedUser.map(user => user.id)));

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
    <div style="position: sticky; top: 0; z-index: 1000; padding: 10px; border-bottom: 1px solid #ccc;">
        <form @submit.prevent="getChats">
            <div class="input-group mb-3">
                <input type="text" class="form-control" placeholder="Search users" v-model="chChat">
                <button class="btn btn-secondary" type="submit">Search</button>
            </div>
        </form>
    </div>

	<div class="custom-scroll">
		<div v-for="chat in chats" >
			<RouterLink :to="'/chat/' + chat.id_chat" class="nav-link">
			<div class="row g-0 border rounded my-2 p-2">


					
					<div class="col-md-1 col-2">
						<img v-if="chat.propic" :src="'data:image/png;base64,'+ chat.propic" class=" rounded-circle me-2" width="70" height="70"  role="img" focusable="false" style="object-fit: cover;">
					</div>
					
					<div class="col-md-11 col-10">
						<div class="card-body">
							
							<h5>{{chat.nome}}</h5>
							
							<div class="d-flex justify-content-between">
								<p v-if="chat.gruppo && chat.lastmsg !== ''" class="text-wrap">
									<span class="text-capitalize fw-bold">{{ chat.username }}</span>: {{ chat.lastmsg }}  
								</p>
								
								<p v-else> 
									<span>{{ chat.lastmsg  }}</span>
								</p>
							</div>
							
						</div>
					</div>
					
					
				</div>
			</RouterLink>
		
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