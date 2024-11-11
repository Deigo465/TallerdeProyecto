<template>
    <div>
        <button @click="showModal = true" class="btn btn-primary md:text-sm text-xs py-2 px-2 md:px-6">Agregar doctor</button>
        <modal :showModal="showModal" @close="showModal = false">
            <doctor-form :saved-doctor="savedDoctor" @close="showModal = false; this.$emit('reload')"/>
        </modal>
    </div>
</template>


<script>
export default {
    name: 'BtnAddDoctor',
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