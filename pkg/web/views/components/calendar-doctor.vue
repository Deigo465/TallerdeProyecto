<template>
    <div>
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
                day.substring(8, 10) }} </button>
        </div>
        <div class="flex md:grid md:grid-cols-8 md:p-0 p-4 relative group" id="calendarWrapper">
            <div class="w-8 md:w-full mt-10 text-xs text-gray-400" id="hoursWrapper">
                <div v-for="i in hours">
                    <time class="block h-8 px-2 py-1 text-gray-600 pb-5">{{i}}</time>
                    <hr class="border-b border-gray-100 w-full">
                </div>
            </div>
            <!-- day line -->
            <div id="dayLine"
                class="bg-red-500 transition-opacity opacity-20 group-hover:opacity-100 h-full w-px absolute top-0 left-0">
            </div>
            <!-- hour line -->
            <div id="hourLine"
                class="bg-red-500 transition-opacity opacity-20 group-hover:opacity-100 w-full h-px absolute top-0 right-0">
            </div>
            <calendar-day-doctor :appointments="filterByDate(day)"
                :class="selectedDate == day ? 'grow' : 'hidden sm:block' " :date="day" v-for="day in days"
                </calendar-day-doctor>

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

                    // this.appointments = data.data;
                    for (let i = 0; i < data.data.length; i++) {
                        let found = false;
                        for (let j = 0; j < this.appointments.length; j++) {
                            // only update the appointment if it already exists
                            if (this.appointments[j].id === data.data[i].id) {
                                this.appointments[j].status = data.data[i].status;
                                found = true;
                            }
                        }
                        if (!found) {
                            this.appointments.push(data.data[i]);
                        }
                    }
                });
        },
        filterByDate(date) {
            return this.appointments.filter(appointment => appointment.starts_at.substring(0, 10) === date);
        },
        updateLines() {
            let hourLine = document.getElementById('hourLine');
            let hoursWrapper = document.getElementById('hoursWrapper');
            let date = new Date();
            let hour = date.getHours();
            let minutes = date.getMinutes();
            let top = (hour - 8) * 64 + minutes * 1 + 50;
            hourLine.style.top = `${top}px`;
            // hoursWrapper.style.top = `-${top}px`;

            // day line
            let dayLine = document.getElementById('dayLine');
            let day = date.getDay();
                let calendarWrapper = document.getElementById('calendarWrapper');
                let left = (day) * 14.285714285714286;
                dayLine.style.left = `${left}%`;
                calendarWrapper.scrollLeft = left;

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

            this.updateLines()
            window.setInterval(() => {
                this.updateLines()
            }, 1000)


            window.setInterval(() => {
                this.fetchAppointments();
            }, 60*1000) // get appointments every 60 seconds
        }
    }
    </script>
    
    <style scoped></style>