<template>
    <!-- Header de las seccion doctores -->
    <section class="max-w-screen-xl mx-auto md:mt-10 mt-0 px-5 md:px-4 lg:px-5">
        <div class="lg:flex md:justify-between md:items-end">
            <div class="flex md:space-x-10 justify-between items-center">
                <p class="md:text-4xl text-2xl font-bold">Pacientes</p>
                <btn-add-patient @reload="fetchPatients()" />
            </div>
            <div class="flex items-end md:space-x-7 space-x-3 justify-end mt-5 md:mt-0">
                <div>
                    <label for="document" class="form-label text-sm md:text-base">Buscar por número de DNI</label>
                    <input id="document" name="document" type="document" class="form-input md:text-xs text-sm h-7"
                        v-model="query">
                </div>
                <div>
                    <button type="submit" class="btn btn-primary text-sm py-2 px-6">Buscar</button>
                </div>
            </div>
        </div>
    </section>

    <!-- tabla de doctores -->
    <section class="max-w-screen-xl mx-auto mt-10 px-5 md:px-4 lg:px-5">
        <table class="table-auto w-full">
            <thead>
                <tr class="font-bold text-left text-base border-b border-black">
                    <th class="py-2">Nombre</th>
                    <th class="py-2 px-2 hidden sm:table-cell">Género</th>
                    <th class="py-2 px-2 hidden sm:table-cell">F.Nacimiento</th>
                    <th class="py-2 px-2 relative">Doc.Identidad</th>
                </tr>
            </thead>
            <tbody>
                <!-- Si no hay DNI similares, mostrar eso -->
                <tr v-show="filteredPatients.length === 0 && patients.length >= 1">
                    <td colspan="4" class="py-4">
                        <div class="flex items-center text-gray-900">
                            <div>
                                <p class="font-semibold mb-2">No se encontraron pacientes con ese DNI '{{ query }}'.</p>
                                <button @click="query = ''" class="underline text-primary">Volver a buscar</button>
                            </div>
                        </div>
                    </td>
                </tr>
                <!-- Placeholder pacientes -->
                <tr v-show="isLoading == false && patients.length === 0" class="text-center">
                    <td colspan="4">
                        <div class="flex justify-center items-center mt-16">
                            <div class="w-1/2">
                                <img src="img/patient-placeholder.png" alt="No hay pacientes registrados"
                                    class="mx-auto object-contain">
                            </div>
                            <div class="w-1/2">
                                <p class="text-4xl mb-10">Parece que no<br> tienes a ningún<br> paciente registrado</p>
                                <button @click="showModal = true"
                                    class="underline text-primary text-2xl font-bold">Agregar a un paciente</button>
                                <modal :showModal="showModal" @close="showModal = false">
                                    <patient-form />
                                </modal>
                            </div>
                        </div>
                    </td>
                </tr>

                <tr class="border-b border-gray-400" v-for="patient in filteredPatients" :key="patient.document_number">
                    <td class="py-2 relative">
                        {{ patient.first_name }} {{ patient.father_last_name }} {{ patient.mother_last_name }}
                    </td>
                    <td class="py-2 px-2 hidden sm:table-cell">{{ patient.gender }}</td>
                    <td class="py-2 px-2 hidden sm:table-cell">{{ (new Date(patient.date_of_birth)).toLocaleDateString() }}</td>
                    <td class="py-2 px-2">
                        <span v-html="highlightMatch(patient.document_number)"></span>
                    </td>
                    <td class="py-2">
                        <btn-view-profile-patient :patient="patient"/>
                    </td>
                </tr>
            </tbody>
        </table>
    </section>
</template>

<script>
export default {
    name: 'PatientList',
    data() {
        return {
            isLoading: false,
            query: '',
            patients: [],
            showModal: false
        }
    },
    computed: {
        filteredPatients() {
            const lowerCaseQuery = this.query.toLowerCase();
            return this.patients.filter(patient =>
                patient.document_number.toLowerCase().includes(lowerCaseQuery)
            );
        }
    },
    mounted() {
        this.fetchPatients();
    },
    methods: {
        fetchPatients() {
            this.isLoading = true;
            fetch("/api/v1/patients")
                .then(response => response.json())
                .then(data => {
                    this.patients = data.data;
                    this.isLoading = false;
                })
                .catch(error => {
                    this.isLoading = false;
                    console.error(error);
                });
        },
        highlightMatch(document_number) {
            if (!this.query) return document_number;
            const regex = new RegExp(`(${this.query})`, 'gi');
            return document_number.replace(regex, '<span class="bg-red-400 text-white">$1</span>');
        }
    }
}
</script>

<style scoped>
/* Your component's CSS code goes here */
</style>
