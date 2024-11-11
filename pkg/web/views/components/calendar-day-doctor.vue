<template>
    <div class="h-full">
        <h2 class="font-medium text-center text-sm">{{ toDay(date) }}</h2>
        <div class="border border-gray-100 mt-[18px]">
            <div v-for="hour in hours" :key="hour" class="h-[34px] py-0.5 border-b border-gray-100" :data-time="hour">
                <div class="flex gap-1 items-stretch h-full">
                    <calendar-appointment-doctor
                        v-for="appointment in filterByTime(hour)" 
                        :key="appointment.name" 
                        :appointment="appointment" 
                        class="text-center py-1 px-2 h-full"
                    ></calendar-appointment-doctor>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    props: {
        date: {
            type: String,
            required: true
        },
        appointments: {
            type: Array,
            required: true
        }
    },
    data() {
        return {
            hours: [],
        };
    },
    methods: {
        toDay(isoString){
            const date = new Date(isoString);
            const options = { weekday: 'long', day: 'numeric', timeZone: 'UTC'};
            let formattedDate = date.toLocaleDateString('es-ES', options);
            formattedDate = formattedDate.charAt(0).toUpperCase() + formattedDate.slice(1);
            return formattedDate
        },

        filterByTime(time) {
            return this.appointments.filter(event => event.starts_at.substring(11, 16) === time);
        }
    },
    mounted() {
        for (let i = 0; i < 24; i++) {
            let start = new Date();
            start.setUTCHours(i, 0, 0, 0);
            this.hours.push(start.toISOString().substring(11, 16));
        }
    }
};
</script>

<style scoped>
</style>
