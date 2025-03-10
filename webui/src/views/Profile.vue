<script>
export default {
	data: function() {
		return {}
	},
	methods:{
		profileButtonHandler() {
			this.$router.push("/profile");
		}
	}
}
</script>
<template>
    <div class="row">
        <div class="col-4"></div>
        <div class="col-2 mt-2">
            <svg class="bd-placeholder-img rounded-circle" width="140" height="140" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="Placeholder" preserveAspectRatio="xMidYMid slice" focusable="false"><title>Placeholder</title><rect width="100%" height="100%" fill="var(--bs-secondary-color)"></rect></svg>
            <h2 class="fw-normal">Heading</h2>
            <p>Some representative placeholder content for the three columns of text below the carousel. This is the first column.</p>
            <p><a class="btn btn-secondary" href="#">View details »</a></p>
            <div class="col-6"></div>
        </div>
    </div>
        
</template>

