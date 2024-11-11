<template>
    <div>
        <button @click="showModal = true" class="btn btn-tertiary">Ver más</button>
        <modal :showModal="showModal" @close="showModal = false">
            <record-detail :record="record" :patient="patient" @close="showModal = false"/>
        </modal>
    </div>
</template>


<script>
export default {
    name: 'BtnViewRecord',
    props: {
        record:{
            type: Object,
            required: true
        },
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