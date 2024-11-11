<template>
    <button @click="showModal = true" class="underline hover:text-blue-500">Ver Perfil</button>
    <modal :showModal="showModal" @close="showModal = false">
        <section class="p-5">
            <!-- Header -->
            <div class="flex justify-between items-center mb-6">
                <h1 class="lg:text-2xl text-xl font-bold">Perfil</h1>
                <div class="flex text-primary underline font-bold items-center">
                    <btn-edit-patient :savedPatient="patient"/>
                </div>
            </div>

            <!-- nombre -->
            <h1 class="lg:text-3xl md:text-2xl font-bold mb-5"> {{ fullName }}</h1>
            
            <!-- datos -->
            <div class="flex justify-center">
                <table class="w-full mb-2 justify-center">
                    <thead>
                        <tr class="font-bold text-center lg:text-base text-sm border-b border-black">
                            <th class="md:px-6 px-1 py-1">DNI</th>
                            <th class="md:px-6 px-1 py-1">F. Nac.</th>
                            <th class="md:px-6 px-1 py-1 hidden lg:table-cell">Teléfono</th>
                            <th class="md:px-6 px-1 py-1 hidden lg:table-cell">Correo electrónico</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr class="text-center">
                            <td class="md:px-6 px-1 py-1">{{ patient.document_number }}</td>
                            <td class="md:px-6 px-1 py-1">{{ (new Date(patient.date_of_birth)).toLocaleDateString() }}</td>
                            <td class="md:px-6 px-1 py-1 hidden lg:table-cell">{{ patient.phone }}</td>
                            <td class="md:px-6 px-1 py-1 hidden lg:table-cell">{{ patient.email }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <!-- responsive de la tabla -->
            <div class="flex justify-center">
                    <table class="w-full mb-2 justify-center">
                        <thead>
                            <tr class="font-bold text-center lg:text-base text-sm border-b border-black lg:hidden">
                                <th class="md:px-6 px-1 py-1">Teléfono</th>
                            <th class="md:px-6 px-1 py-1">Correo electrónico</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr class="text-center lg:hidden lg:text-base text-sm">
                                <td class="md:px-6 px-1 py-1">{{ patient.phone }}</td>
                            <td class="md:px-6 px-1 py-1">{{ patient.email }}</td>
                            </tr>
                        </tbody>
                    </table>
            </div>


            <!-- Historial de citas -->
            <p class="text-xl font-bold mb-2">
                    Historial de citas
            </p>
            <table class="w-full md:mb-10 mb-3">
                <thead>
                    <tr class="font-bold text-center  md:text-sm text-xs border-b border-black">
                        <th class="py-2 lg:px-2 md:px-4 px-1 text-left">Doctor/a</th>
                        <th class="py-2 lg:px-2 md:px-4 px-1 text-left">Especialidad</th>
                        <th class="py-2 lg:px-2 md:px-4 px-1 text-left hidden sm:table-cell">Fecha</th>
                        <th class="py-2 lg:px-2 md:px-4 px-1 text-left hidden lg:table-cell">Motivo</th>
                        <th class="py-2 lg:px-2 md:px-4 px-1">Estado</th>
                    </tr>
                </thead>
                <tbody>
                    <tr class="border-b border-gray-400 text-center lg:text-base md:text-md text-xs" v-for="appointment in appointments">
                        <td class="py-2 lg:px-2 md:px-4 px-1 text-left">{{ appointment.doctor.first_name + " " + appointment.doctor.father_last_name}}</td> 
                        <td class="py-2 lg:px-2 md:px-4 px-1 text-left">{{ appointment.specialty }}</td>
                        <td class="py-2 lg:px-2 md:px-4 px-1 text-left hidden sm:table-cell">{{ (new Date(appointment.starts_at)).toLocaleDateString() }}</td> 
                        <td class="py-2 lg:px-2 md:px-4 px-1 text-left hidden lg:table-cell">{{ appointment.description }}</td>
                        <td class="py-2 lg:px-2 md:px-4 px-1">
                            <span class="inline-block btn rounded-2xl lg:text-sm text-xs lg:w-24 w-20" :class="statusColorBg[appointment.status]"> 
                                {{ appointment.status}}
                            </span>
                        </td>
                    </tr>
                </tbody>
            </table>
            <div class="flex items-center space-x-5">
                    <button class="btn btn-secondary mt-10 px-10 lg:text-xl md:text-sm text-xs" @click="showModal = false">Atrás</button>
            </div>
        </section>
    </modal>

</template>


<script>
export default {
    name: 'BtnViewProfilePatient',
    props: {
        patient: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            showModal: false,
            statusColorBg: {
                "Agendada": "bg-yellow-400",
                "Pagada": "bg-green-400",
                "Terminada": "bg-gray-400",
                "En proceso": "bg-blue-400",
                "Anulada": "bg-red-400",
            },
            appointments: []
        }
    },
    computed: {
        fullName() {
            return `${this.patient.first_name} ${this.patient.father_last_name} ${this.patient.mother_last_name}`;
        }
    },
    mounted() {
        this.fetchPatientData();
    },
    methods: {
        fetchPatientData(){
            fetch("/api/v1/patients/" + this.patient.id + "/appointments")
                .then(response => response.json())
                .then(data => {
                    this.appointments = data.data;
                })
                .catch(error => console.error(error))
        }
    }

}

</script>

<style scoped>
</style>