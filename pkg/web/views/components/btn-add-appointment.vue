<template>
    <div>
        <button @click="showModal = true" class="btn btn-primary text-sm py-2 px-6">Agregar citas</button>
        <modal :showModal="showModal" @close="showModal = false">
            <appointment-form @close="showModal = false" @reload="reload"/>
        </modal>
    </div>
</template>

<script>
export default {
    name: 'BtnAddAppointment',
    methods: {
        reload() {
            window.location.reload()
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