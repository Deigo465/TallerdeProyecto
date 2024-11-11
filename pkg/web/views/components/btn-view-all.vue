<template>
    <div>
        <button @click="showModal = true" class="btn btn-tertiary">Ver más</button>
        <modal :showModal="showModal" @close="showModal = false">
            <section class="p-5">
                <div class="flex justify-end space-x-2">
                    <btn-add-register :patient="appointment.patient" />
                    <div>
                        <button v-if="appointment.status == 2 && !isLoading" @click="endAppointment" class="btn bg-red-700 text-white">Finalizar cita</button>
                        <div v-if="appointment.status == 3">Finalizada</div>
                    </div>
                </div>
                <p class="text-2xl my-5">
                    {{ appointment.patient.first_name }}
                    {{ appointment.patient.father_last_name }}
                    {{ appointment.patient.mother_last_name }}
                </p>
                <!-- datos -->
                <div class="flex justify-center">
                    <table class="w-full justify-center">
                        <thead>
                            <tr class="font-bold text-center lg:text-base text-sm border-b border-black">
                                <th class="md:px-6 px-1 py-1">DNI</th>
                                <th class="md:px-6 px-1 py-1">F. Nac.</th>
                                <th class="md:px-6 px-1 py-1">Teléfono</th>
                                <th class="md:px-6 px-1 py-1">Correo electró</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr class="text-center">
                                <td class="md:px-6 px-1 py-1">{{ appointment.patient.document_number }}</td>
                                <td class="md:px-6 px-1 py-1">{{ (new Date(appointment.patient.date_of_birth)).toLocaleDateString() }}</td>
                                <td class="md:px-6 px-1 py-1">{{ appointment.patient.phone }}</td>
                                <td class="md:px-6 px-1 py-1    ">{{ appointment.patient.email }}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <div class="mt-5">
                    <h2 class="md:text-xl text-sm mb-4">Descripción</h2>
                    <textarea class="text-gray-700 border rounded border-black font-normal w-full h-32 p-4"
                        :value="appointment.description"></textarea>
                </div>

                <!-- Historial de registros -->
                <p class="lg:text-3xl md:text-2xl text-xl font-bold text-left my-5">
                    Historial de registros
                </p>
                <table class="w-full md:mb-10 mb-5" v-if="appointment.status == 2">
                    <thead>
                        <tr class="font-bold text-center lg:text-xl md:text-sm text-xs border-b border-black">
                            <th class="py-2 lg:px-0 md:px-5 px-1 text-left">Doctor/a</th>
                            <th class="py-2 lg:px-5 md:px-5 px-1 text-left">Especialidad</th>
                            <th class="py-2 lg:px-5 md:px-5 px-1 text-left hidden sm:table-cell">Fecha</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr class="border-b border-gray-400 text-center lg:text-base md:text-sm text-xs" v-for="record in appointment.patient.records">
                            <td class="py-2 lg:px-0 md:px-5 px-1 text-left">{{  record.doctor.first_name  }}</td>
                            <td class="py-2 lg:px-5 md:px-5 px-1 text-left font-normal">{{ record.specialty }}</td>
                            <td class="py-2 lg:px-5 md:px-5 px-1 text-left font-normal hidden sm:table-cell">
                                {{ record.created_at }}
                            </td>
                            <td>
                                <btn-view-record :patient="appointment.patient" :record="record" />
                            </td>
                        </tr>
                    </tbody>
                </table>
                <div v-else> Cita finalizada</div>

                <div class="flex items-center space-x-5">
                    <button class="btn btn-secondary mt-10 px-10 lg:text-xl md:text-sm text-xs"
                        @click="showModal = false">Cerrar</button>
                </div>


            </section>
        </modal>
    </div>
</template>

<script>
export default {
    name: 'BtnViewAll',
    props: {
        appointment: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            isLoading: false,
            showModal: false
        }
    },
    methods: {
        endAppointment() {
            this.isLoading = true
            fetch(`/api/v1/appointments/${this.appointment.id}/end`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json'
                }
            }).then(response => {
                if (response.ok) {
                    this.appointment.status = 3
                    this.emit('reload')
                } else {
                    alert('No se pudo finalizar la cita')
                }
                this.isLoading = false
            });
        }
    }
}

</script>

<style scoped></style>