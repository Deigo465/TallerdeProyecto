<template>
    <div>
        <button @click="showModal = true" class="btn btn-secondary text-sm py-2 px-4">+ Nuevo paciente</button>
        <modal :showModal="showModal" @close="showModal = false">
            <patient-form  @close="showModal = false; this.$emit('reload')"/>
        </modal>
    </div>
</template>


<script>
export default {
    name: 'BtnAddPatient',
    props: { },
    data() {
        return {
            showModal: false
        }
    }
}
</script>

<style scoped>
</style>