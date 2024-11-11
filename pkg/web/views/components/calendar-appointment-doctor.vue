<template>
    <button @click.stop="keepPopupOpen = !keepPopupOpen" class="px-2 py-1 relative h-full w-full"
        :class="statusMap[appointment.status]">
        <div class="text-xs uppercase">
        {{ appointment.patient.first_name }} {{ appointment.patient.father_last_name }}
        </div>
        <div id="popup" class="absolute top-5 left-10 px-6 z-10 bg-gray-200 border-2 border-gray-500 md:w-96 w-72 pt-4 pb-4 cursor-default"
            @click.stop="" v-show="showPopup || keepPopupOpen">
            <div class="flex text-left justify-between">
                <p class="font-bold text-2xl">{{ appointment.patient.first_name }} {{ appointment.patient.father_last_name }}</p>
                <btn-view-all :appointment="appointment" v-if="loadedRecords" @reload="appointment.patient.records = []"/>
            </div>
            <div class="flex">
                <p class="font-bold text-sm">F. Nac.:</p>
                <p class="font-normal text-sm ml-2">{{(new Date(appointment.patient.date_of_birth )).toLocaleDateString() }}</p>
            </div>
            <p class="font-normal text-sm bg-gray-300 px-1 py-1 mt-2">{{ appointment.description }}</p>
            <div class="flex justify-between text-sm mt-5">
                <p class="font-bold text-lg">Historia clínica</p>
            </div>
            <div class="mx-auto" v-show="records.length > 0">
                <table class="w-full">
                    <tbody>
                        <tr class="border-b border-black" v-for="record in appointment.patient.records">
                            <td class="font-normal py-2 text-sm text-gray-600">{{ (new Date(record.created_at)).toLocaleDateString() }}</td>
                            <td class="px-2 font-normal py-2 text-sm">{{ record.specialty }}</td>
                            <td class="font-normal text-sm">
                                <btn-view-record :patient="appointment.patient" :record="record" />
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
            <div class="mt-5" v-show="!loadedRecords">
                <button class="btn btn-primary w-full relative" @click="fetchPatient">
                    <div>
                        Obtener datos de paciente
                    </div>
                    <div v-if="isLoading" class="absolute right-2 inset-y-0 flex items-center">
                        <!-- Spinning circle -->
                        <svg class="animate-spin  h-5 w-5 text-blue-100" xmlns="http://www.w3.org/2000/svg" fill="none"
                            viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4">
                            </circle>
                            <path class="opacity-75" fill="currentColor"
                                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                            </path>
                        </svg>
                    </div>
                </button>
            </div>
            <div class="mt-5" v-if="loadedRecords">
                <btn-add-register @click.stop="" :patient="appointment.patient" @reload="fetchPatient()"/>
            </div>
        </div>
    </button>
</template>

<script>

export default {
    name: 'CalendarAppointmentDoctor',
    props: {
        appointment: {
            type: Object,
            required: true
        }
    },
    methods: {
        fetchPatient() {
            this.isLoading = true
            window.setTimeout(()=>{
                fetch(`/api/v1/patients/${this.appointment.patient.id}/records`)
                    .then(response => response.json())
                    .then(response => {
                        if (response.status > 400) {
                            alert("No tienes permisos para acceder a esta historia clínica")
                            return
                        }
                        this.records = response.data
                        this.loadedRecords = true
                        if (this.records){
                            this.records = this.records.reverse()
                            this.loadedRecords = true
                        }  else{
                            this.records = []
                        }
                    })
                    .finally(() => {
                        this.isLoading = false
                        this.appointment.patient.records = this.records
                    })
            }, 100 )
        }
    },
    data() {
        return {
            isLoading: false,
            loadedRecords: false,
            records: [],
            statusMap: {
                0: "bg-gray-400 hover:bg-gray-500", // pending
                1: "bg-green-300 hover:bg-green-400", // paid
                2: "bg-yellow-300 hover:bg-yellow-400", // in progress
                3: "bg-blue-300 hover:bg-blue-400", // done
                4: "bg-red-300 hover:bg-red-400" // cancelled
            },
            keepPopupOpen: false,
            showPopup: false,
            showModal: false,
        }
    }
}
</script>

<style scoped></style>