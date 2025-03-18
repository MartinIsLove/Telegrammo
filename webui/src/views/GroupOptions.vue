<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,

			
			id: null,
            chatId: null,
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
		}
	},
	mounted() {
		this.refresh()
	}
}
</script>
<template>
	<button @click="leavegroup" class="btn btn-secondary">
		leave group
	</button>
</template>