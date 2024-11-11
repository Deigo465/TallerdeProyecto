<template>
    <div>
        <button @click="showModal = true" class="btn btn-primary w-full px-10 text-sm">Nuevo registro</button>
        <modal :showModal="showModal" @close="showModal = false">
            <register-form :patient="patient" @close="showModal = false" @reload="showModal = false; $emit('reload')"/>
        </modal>
    </div>
</template>


<script>
export default {
    name: 'BtnAddRegister',
    props: {
        patient: {
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