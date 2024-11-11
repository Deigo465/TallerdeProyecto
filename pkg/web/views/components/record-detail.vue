<template>
    <div class="p-5 min-w-96 max-w-screen-lg">
        <h1 class="lg:text-4xl text-2xl font-bold mb-5">
            <span>{{ record.specialty }}</span>
        </h1>
        <p class="lg:text-2xl md:text-xl text-xl font-bold">Fecha: 01/02/24</p>
        <p class="lg:text-xl md:text-sm text-xs">
            Dr/a {{ record.doctor.first_name }} {{ record.doctor.father_last_name }} ({{ record.doctor.specialty }}) - {{ record.doctor.cmp }}
        </p>
        <p class="text-gray-700 border rounded border-black font-normal w-full h-32 p-4">
            {{ record.body }}
        </p>
        <div class="flex justify-between items-center mt-4 gap-4 flex-wrap">
            <div v-for="file in record.files" class="bg-gray-100 px-4 py-4 border max-w-lg">
                {{  file.name.substring(0, 70) }} - ({{  file.mime_type }})
                <a :href="'/downloads'+file.url" target="_bank" class="text-blue-500 underline">Descargar</a>
                ( {{ Math.floor(file.file_size / 1024) }} KB)
            </div>
        </div>

        <button type="button" @click="$emit('close')" class="btn btn-secondary mt-4 w-full">Cerrar</button>
    </div>
</template>

<script>
export default {
    name: 'RecordDetail',
    props: {
        record: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
        }
    },
    methods: {
    },

};
</script>

<style scoped>
/*Your component's CSS styles go here */
</style>
