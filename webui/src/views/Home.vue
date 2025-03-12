<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/");
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		createChatButtonHandler() {
			this.$router.push("/chat");
		}
	},
	mounted() {
		this.refresh()
	}
}
</script>
<template>
    
		<div class="position-absolute top-1 start-1">
			<button class="btn btn-secondary ms-2 mb-2" @click="createChatButtonHandler">
				Create Chat
			</button>
		</div>
	
</template>