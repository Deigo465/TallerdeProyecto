<template>

    <div class="mb-10">
        <div class="flex gap-2">
            <button class="btn btn-secondary" @click="changeWeek(-7)">
                <!-- Chevron left icon -->
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="size-6 text-blue-800">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5 8.25 12l7.5-7.5" />
                </svg>

            </button>
            <button class="btn btn-secondary" @click="changeWeek(7)">
                <!-- Chevron right icon -->
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="size-6 text-blue-800">
                    <path stroke-linecap="round" stroke-linejoin="round" d="m8.25 4.5 7.5 7.5-7.5 7.5" />
                </svg>

            </button>
        </div>
        <!--Este es el segundo div para mostrar los dias-->
        <div class="flex md:hidden space-x-2 justify-center">
            <button type="button" class="px-2 py-1 bg-blue-100" @click="selectedDate = day" v-for="day in days"> {{
                day.substring(0, 2) }} </button>
        </div>
        <div class="flex md:grid md:grid-cols-8 md:p-0 p-4">
            <div class="w-8 md:w-full relative mt-10">
                <div v-for="i in hours">
                    <time class="block h-8 px-2 py-1 text-gray-600 pb-5">{{i}}</time>
                    <hr class="border-b border-gray-100 w-full">
                </div>
                <div class="border-l-2 border-red-500 h-full w-8 absolute top-0 right-0"></div>
            </div>

            <calendar-day @reload="fetchAppointments()" :appointments="filterByDate(day)"
                :class="selectedDate == day ? 'grow' : 'hidden sm:block'" :date="day" v-for="day in days"
                class="underline font-bold"></calendar-day>

        </div>
    </div>
</template>
<script>
export default {
    data() {
        return {
            selectedDate: "Lunes 17",
            days: [],
            hours: [],
            appointments: [],
            reference: null
        }
    },
    methods: {
        getMonday(d) {
            d = new Date(d);
            var day = d.getDay(),
                diff = d.getDate() - day + (day == 0 ? -6 : 1); // adjust when day is sunday
            return new Date(d.setDate(diff));
        },
        changeWeek(id) {
            this.reference.setDate(this.reference.getDate() + id);
            this.days = [];
            for (let i = 0; i < 7; i++) {
                let start = this.getMonday(this.reference)
                let day = new Date(start.setDate(start.getDate() + i));
                day.setHours(10, 0, 0, 0)
                this.days.push(day.toISOString().substring(0, 10));
            }
        },
        fetchAppointments() {
            fetch("/api/v1/appointments")
                .then(response => response.json())
                .then(data => {
                    console.log(data.data)

                    this.days = []
                    for (let i = 0; i < 7; i++) {
                        let start = this.getMonday(this.reference)
                        let day = new Date(start.setDate(start.getDate() + i));
                        day.setHours(10, 0, 0, 0)
                        this.days.push(day.toISOString().substring(0, 10));
                    }

                    this.appointments = data.data;
                });
        },
        filterByDate(date) {
            return this.appointments.filter(appointment => appointment.starts_at.substring(0, 10) === date);
        }
    },
    mounted() {
        for (let i = 0; i < 24; i++) {
            let start = new Date();
            start.setUTCHours(i, 0, 0, 0);
            this.hours.push(start.toISOString().substring(11, 16));
        }


        let today = new Date();
        this.reference = today;
        this.fetchAppointments();

        window.setInterval(() => {
            this.fetchAppointments();
        }, 60000) // get appointments every 60 seconds
    }
}
</script>

<style scoped></style>