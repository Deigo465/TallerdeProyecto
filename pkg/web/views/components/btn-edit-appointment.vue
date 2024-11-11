<template>
    <div>
        <button @click="showModal = true" class="btn btn-tertiary">Editar cita</button>
        <modal :showModal="showModal" @close="showModal = false">
            <appointment-form :saved-appointment="savedAppointment" @close="showModal = false" @reload="showModal = false; $emit('reload')"/>
        </modal>
    </div>
</template>


<script>
export default {
    name: 'BtnEditAppointment',
    props: {
        savedAppointment: {
            type: Object,
            required: true
        },
    },
    data() {
        return {
            showModal: false,
        }
  },
        
}

</script>

<style scoped>
</style>