<script>
    export default {
        data: function() {
            return {
                errormsg: null,
                loading: false,
                some_data: null,

                username: '',
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
            async login(){
                try {
                    
                    let response = await this.$axios.post("/user", {username: this.username});
                
                    sessionStorage.setItem("cs", response.data)
                    this.$router.push("/home")
                    
                   
                } catch (e) {
                    this.errormsg = e.toString();
                }
            },
        },
        mounted() {
            // this.refresh()

	    }
    }

</script>
<template>
    <form @submit.prevent="login">
        <div class="mb-3">
            <label for="username" class="form-label">username</label>
            <input type="text" class="form-control" id="username" aria-describedby="inserisci username" v-model="username">
            <div id="emailHelp" class="form-text">We'll never share your username with anyone else.</div>
        </div>
        <button type="submit" class="btn btn-primary">Submit</button>
    </form>
</template>