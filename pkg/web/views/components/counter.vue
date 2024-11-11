<template>
    <button @click="count++" class="btn btn-primary">Clicked {{ count }} times</button>
</template>

<script>
export default {
    data() {
        return {
            count: 0
        }
    }
}
</script>

<style scoped>
</style>