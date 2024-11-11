<template>
   <div @click="keepPopupOpen = !keepPopupOpen" class="px-2 py-1 relative h-full w-6 text-center"
   :class="statusMap[appointment.status]"
   >
        <div id="popup" class="absolute top-0 px-6 z-10 bg-gray-200 border-2 border-gray-500 md:w-80 w-72  pt-4 pb-4" v-show="showPopup || keepPopupOpen">
            <div class="flex justify-between">
                <p class="font-bold">{{toDay(appointment.starts_at)}}</p>
                <btn-edit-appointment :savedAppointment="appointment" @click.stop="" @reload="$emit('reload')" v-if="notInPast()"/>
            </div>   
            <div class="text-left text-sm">
                <p class="font-bold mt-5">Paciente</p>
                <p class="font-normal">{{ appointment.patient.first_name }} {{ appointment.patient.father_last_name }} {{ appointment.patient.mother_last_name }}</p>
                <p class="font-normal mt-2">DNI {{appointment.patient.document_number}}</p>
            </div>
            <div class="flex justify-between mt-5">
                <p class="font-bold text-sm">Doctor</p>
                <p class="font-bold text-sm">Especialidad</p>
            </div>  
            <div class="flex justify-between  text-sm">
                <p class="font-normal">{{appointment.doctor.first_name}} {{appointment.doctor.father_last_name}} {{appointment.doctor.mother_last_name}}</p>
                <p class="font-normal"> {{appointment.specialty}} </p>
            </div>  
            <div class="text-sm text-left">
                <p class="font-bold mt-2">Motivo</p>
                <p class="font-normal"> {{appointment.description}} </p>
            </div>
            <div class="text-left text-sm">
                <p class="font-bold mt-2 ">Estado</p>
                <p class="font-normal "> {{ statusLabel[appointment.status] }} </p>
                <hr class="border-black mt-2">
            </div>
        </div>
    </div>
</template>

<script>

export default {
    name: 'CalendarAppointment',
    props: {
        appointment: {
            type: Object,
            required: true
        },
    },
    data() {
        return {
            statusMap: {
                0: "bg-gray-400 hover:bg-gray-500", // pending
                1: "bg-green-300 hover:bg-green-400", // paid
                2: "bg-yellow-300 hover:bg-yellow-400", // in progress
                3: "bg-blue-300 hover:bg-blue-400", // done
                4: "bg-red-300 hover:bg-red-400" // cancelled
            },
            statusLabel: ["Pendiente", "Pagada", "En proceso", "Finalizada", "Terminada"],
            keepPopupOpen: false,
            showPopup: false,
            showModal: false,
        }
    },
    methods: {
        notInPast(){
            let yesterday = new Date()
            yesterday = yesterday.setDate(yesterday.getDate() - 1)
            const appointmentDate = new Date(this.appointment.starts_at);
            return yesterday < appointmentDate;
        },
        toDay(isoString) {
            const starts_at = new Date(isoString);
            const options = { weekday: 'long', day: 'numeric', timeZone: 'UTC' };
            let formattedStartsAt = starts_at.toLocaleDateString('es-ES', options);
            formattedStartsAt = formattedStartsAt.charAt(0).toUpperCase() + formattedStartsAt.slice(1);
            return formattedStartsAt
        }
    }
}
</script>

<style scoped>

</style>