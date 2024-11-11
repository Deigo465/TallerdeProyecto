<template>

    <!-- Header de las seccion doctores -->
    <section class="max-w-screen-xl mx-auto px-5 md:px-4 lg:px-5">
        <div class="md:flex justify-between items-center">
            <div class="md:flex gap-10 items-end">
                <p class="md:text-4xl text-lg font-bold md:p-0">Doctores</p>
                <!-- <p class="md:text-sm text-xs">en Clinica San Felipe - Jesús María</p> -->
            </div>
            <div class="flex justify-end mt-2 md:mt-0 md:p-0 p-2">
                <btn-add-doctor @reload="fetchDoctors()"></btn-add-doctor>
            </div>
            <div class="flex items-end md:space-x-7 space-x-3 justify-end mt-5 md:mt-0">
                <div>
                    <label for="document" class="form-label text-sm md:text-base">Buscar por CMP</label>
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
                    <th class="py-1 lg:py-2">Nombre</th>
                    <th class="py-1 lg:py-2 px-2 hidden lg:block">Género</th>
                    <th class="py-1 lg:py-2 px-2">Especialidad</th>
                    <th class="py-1 lg:py-2 px-2 hidden sm:table-cell">CMP</th>
                    <th class="py-1 lg:py-2 px-2 hidden sm:table-cell">Correo</th>
                </tr>
            </thead>
            <tbody>
                <tr v-if="isLoading">
                    <td colspan="4">
                        Cargando...
                    </td>
                </tr>
                <!-- Placeholder -->
                <tr v-if="isLoading == false && doctors.length === 0" class="text-center">
                    <td colspan="4">
                        <div class="flex justify-center items-center mt-16">
                            <div class="w-1/2">
                                <img src="img/doctor-placeholder.png" alt="No hay doctores registrados"
                                    class="mx-auto object-contain">
                            </div>
                            <div class="w-1/2">
                                <p class="text-4xl mb-10">Parece que no tienes a <br> ningún doctor registrado</p>
                                <button @click="showModal = true"
                                    class="underline text-primary text-2xl font-bold">Agregar a un doctor</button>
                                <modal :showModal="showModal" @close="showModal = false">
                                    <doctor-form @close="showModal = false; window.location.reload()"/>
                                </modal>
                            </div>
                        </div>
                    </td>
                </tr>
                <!-- LISTA DE DOCTORES -->
                <tr class="border-b border-gray-400" v-for="doctor in filteredDoctors">
                    <td class="py-1 lg:py-2">{{ doctor.first_name }} {{ doctor.father_last_name }} {{
                        doctor.mother_last_name }}</td>
                    <td class="py-1 lg:py-2 px-2 hidden lg:block">{{ doctor.gender }}</td>
                    <td class="py-1 lg:py-2 px-2">{{ doctor.specialty }}</td>
                    <td class="py-1 lg:py-2 px-2 hidden sm:table-cell">
                        <span v-html="highlightMatch(doctor.cmp)"></span>
                    </td>
                    <td class="py-1 lg:py-2 px-2 hidden sm:table-cell">{{ doctor.email }}</td>
                    <td class="py-1 lg:py-2">
                        <btn-view-profile-doctor :user="doctor" />
                    </td>
                </tr>
            </tbody>
        </table>
    </section>
</template>

<script>
export default {
    name: 'DoctorList',
    data() {
        return {
            isLoading: false,
            doctors: [],
            showModal: false,
            query: '',
        }
    },
    computed: {
        filteredDoctors() {
            const lowerCaseQuery = this.query.toLowerCase();
            return this.doctors.filter(doctor =>
                doctor.cmp.toLowerCase().includes(lowerCaseQuery)
            );
        }
    },
    mounted() {
        this.fetchDoctors()
    },
    methods: {
        fetchDoctors() {
            this.isLoading = true;
            fetch("/api/v1/doctors")
                .then(response => response.json())
                .then(data => {
                    this.doctors = data.data
                    this.isLoading = false
                })
                .catch(error => {
                    this.isLoading = false
                })
        },
        highlightMatch(cmp) {
            if (!this.query) return cmp;
            const regex = new RegExp(`(${this.query})`, 'gi');
            return cmp.replace(regex, '<span class="bg-red-400 text-white">$1</span>');
        }
    }
}
</script>

<style scoped>
/* Your component's CSS code goes here */
</style>