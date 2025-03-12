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
		async getUsers(){
            
        }
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
            <h2 class="fw-normal">username</h2>
            <form class="mt-5" @submit.prevent="getUsers">
                <label for="username" class="form-label ">search user</label>
                <input type="text" class="form-control mb-2" id="username" placeholder="insert username" v-model="chUsername">
                <button class="btn btn-secondary mt-2"  >find</button>
            </form>
            <div class="col-6"></div>
        </div>
    </div>

</template>