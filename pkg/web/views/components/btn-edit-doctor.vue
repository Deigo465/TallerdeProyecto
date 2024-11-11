<template>
    <div>
        <button @click="showModal = true" class="text-primary underline font-bold lg:text-base text-sm">Editar datos</button>
        <modal :showModal="showModal" @close="showModal = false">
            <doctor-form :saved-doctor="savedDoctor" @close="showModal = false; this.$emit('reload')"/>
        </modal>
    </div>
</template>

<script>
export default {
    name: 'BtnEditDoctor',
    props: {
        savedDoctor: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            showModal: false
        }
    }
}

</script>

<style scoped>
</style>