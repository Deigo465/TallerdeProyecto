<template>
    <div>
        <button @click="showModal = true" class="text-primary underline font-bold">Editar datos</button>
        <modal :showModal="showModal" @close="showModal = false">
            <patient-form :saved-patient="savedPatient" @close="showModal = false; this.$emit('reload')"/> 
        </modal>
    </div>
</template>

<script>
export default {
    name: 'BtnEditPatient',
    props: {
        savedPatient: {
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