<template>
    <div class="p-4">
        <div class="">
            <h1 class="lg:text-4xl text-2xl font-bold mb-1">
                <span v-if="!isEdit">Agregar Cita</span>
                <span v-else>Editar Cita </span>
            </h1>
            <hr class="border-gray-800 mt-2">
            <p class="font-semibold text-gray-800 text-lg mt-4">Paciente</p>
        </div>
        <form class="grid grid-cols-1" @submit.prevent="isEdit ? updateAppointment() : addNewAppointment()">
            <div class="border-black mt-1 p-2 border">
                <div v-if="selectedPatient.id == undefined" class="">
                    <div class="">
                        <label for="name">Buscar por DNI</label>
                        <input type="number" name="dni" class="form-input" v-model="query">
                        <ul class="w-full bg-gray-100">
                            <li v-for="patient in filteredPatients" :key="patient.id"
                                class="border border-b flex justify-between py-2 px-5 w-full cursor-pointer"
                                @click="selectedPatient = patient">
                                <div class="font-bold">
                                    {{ patient.first_name }} {{ patient.father_last_name }}
                                </div>
                                <div class="text-lg text-gray-500">
                                    {{ patient.document_number }}
                                </div>
                            </li>
                        </ul>
                    </div>
                </div>
                <div class="flex justify-between items-center" v-else>
                    <div>
                        {{ selectedPatient.first_name }} {{ selectedPatient.father_last_name }} - {{
                            selectedPatient.document_number }}
                    </div>
                    <button type="button" class="btn btn-tertiary " @click="selectedPatient = {}" v-if="isEdit == false">Cambiar</button>
                </div>
            </div>

            <p class="font-semibold text-gray-800 text-lg mt-4">Doctor</p>
            <div class="md:flex md:space-x-10">
                <div class="relative mt-2 flex flex-col">
                    <label class="text-sm">Especialidad</label>
                    <select id="specialty" class="border-black appearance-none bg-white py-2 p-2 md:w-48 border "
                        required v-model="specialty" :disabled="isEdit">
                        <option v-for="specialty in specialties" :value="specialty">{{ specialty }}</option>
                    </select>
                    <!-- Aqui estoy agregando el icono -->
                    <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-black">
                        <svg class="fill-current h-4 w-4 mt-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
                            <path d="M5 8l5 5 5-5H5z" />
                        </svg>
                    </div>
                </div>

                <div class="relative mt-2 flex flex-col">
                    <label for="doctor" class="text-sm">Doctor</label>
                    <select id="doctor" class="border-black appearance-none bg-white py-2 p-2  md:w-96 border" required
                        v-model="doctorId" :disabled="isEdit">
                        <option disabled>Selecciona un doctor</option>
                        <option v-for="doctor in (isEdit ? doctors : filteredDoctors)" :value="doctor.id">{{ doctor.first_name }} {{
                            doctor.father_last_name }} {{ doctor.mother_last_name }}</option>
                    </select>
                    <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-black">
                        <svg class="fill-current h-4 w-4 mt-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
                            <path d="M5 8l5 5 5-5H5z" />
                        </svg>
                    </div>
                </div>
            </div>

            <div class="flex md:space-x-10 space-x-4">
                <div class="flex md:space-x-10 space-x-4">
                    <div class="flex flex-col">
                        <label for="date" class="text-sm mt-2">Fecha</label>
                        <input type="date" id="date" :class="invalidDate ? 'border-red-500' : 'border-black'"
                        class="appearance-none bg-white py-2 p-2 w-48 border"
                            required v-model="date" :min="minDate">
                            <p class="text-red-500 text-xs italic" v-show="invalidDate"> Ingresa una fecha válida.
                            </p>
                    </div>
                </div>

                <div class="relative mt-2 flex flex-col">
                    <label for="time" class="text-sm">Hora</label>
                    <div class="relative">
                        <select name="time" class="border-black appearance-none bg-white py-2 p-2 md:w-48 w-40 border"
                            v-model="time" required>
                            <option :value="String(i-1).padStart(2, '0')+ ':00'" v-for="i in 23">
                                    {{ String(i-1).padStart(2, '0')+ ':00' }} - {{ String(i).padStart(2, '0')+ ':00' }}
                            </option>
                            <option value="23:00">
                                23:00 - 00:00
                            </option>
                        </select>
                        <div class="absolute inset-y-0 right-0 flex items-center px-2 text-black pointer-events-none">
                            <svg class="fill-current h-4 w-4" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
                                <path d="M5 8l5 5 5-5H5z" />
                            </svg>
                        </div>
                    </div>
                </div>
            </div>

            <div class="grid grid-cols-1 mt-2">
                <label class="text-base font-bold">Motivo:</label>
                <textarea class="border-black border bg-white p-5 text-sm w-full"
                    placeholder="Ingrese el motivo de la cita" v-model="description" required :disabled="isEdit"></textarea>
            </div>
            <div class="mt-5">
                <label class="text-lg">Estado</label>
                <div class="grid md:grid-cols-5 grid-cols-3 gap-2 mt-2 font-bold">
                    <button type="button" @click="selectedStatus = status.value"
                        class="text-sm flex items-center space-x-2"
                        :class="selectedStatus == status.value ? 'btn btn-primary' : 'hover:bg-blue-300 bg-blue-100 text-blue-800 border border-blue-200'"
                        v-for="status in statuses">
                        <div :class="selectedStatus == status.value ? 'opacity-100' : 'opacity-0'">
                            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16px" height="16px"
                                baseProfile="basic">
                                <path fill="#2583ef"
                                    d="M8,0C3.582,0,0,3.582,0,8s3.582,8,8,8s8-3.582,8-8S12.418,0,8,0z" />
                                <polygon fill="#fff"
                                    points="7,12 3.48,8.48 4.894,7.066 7,9.172 11.71,4.462 13.124,5.876" />
                            </svg>
                        </div>
                        <div>
                            {{ status.label }}
                        </div>
                    </button>
                </div>
            </div>
            <div
                class="flex flex-col md:flex-row justify-center space-y-5 md:space-y-0 md:space-x-5 items-center mt-10">
                <button class="btn btn-primary relative px-10" type="submit" :disabled="isLoading">
                    <div>
                        {{ isEdit ?  'Guardar' : 'Agregar' }}
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
                <button type="button" class="btn btn-secondary px-10"
                    @click="$emit('close')">Cancelar</button>
            </div>
        </form>

    </div>
</template>

<script>
export default {
    name: 'AppointmentForm',
    data() {
        return {
            isLoading: false,
            isEdit: false,
            selectedStatus: 0,
            specialty: "Cardiología",
            doctorName: "",
            date: "",
            time: "08:00",
            description: "",
            selectedPatient: {},
            query: "",
            doctorId: 0,
            minDate: this.getMinDate(),
            patients: [],
            statuses: [
                { label: "Agendada", value: 0 },
                { label: "Pagada", value: 1 },
                { label: "En proceso", value: 2 },
                { label: "Terminada", value: 3 },
                { label: "Anulada", value: 4 }
            ],
            doctors: [],
            specialties: ["Cardiología", "Dermatología", "Endocrinología", "Gastroenterología", "Geriatría", "Ginecología", "Hematología", "Infectología", "Medicina interna", "Nefrología", "Neumología", "Neurología", "Nutriología", "Oftalmología", "Oncología", "Pediatría", "Psiquiatría", "Reumatología", "Traumatología", "Urología"],
        };
    },
    props: {
        savedAppointment: {
            type: Object,
            required: false,
        }
    },
    methods: {
        getMinDate() {
            // let yesterday = new Date(today.setDate(today.getDate() - 1))
            let today = new Date();
            today.setHours(0, 0, 0, 0);
            return today.toISOString().split('T')[0];
        },
        resetForm() {
            this.selectedPatient = {};
            this.specialty = "";
            this.doctorId = "";
            this.date = "";
            this.time = "";
            this.description = "";
            this.selectedStatus = 0;
        },
        addNewAppointment() {
            if (this.selectedPatient.id == undefined) {
                alert("Debe seleccionar un paciente");
                return;
            }
            if (this.doctorId == 0) {
                alert("Debe seleccionar un doctor");
                return;
            }
            if (this.invalidDate) {
                alert("Debe seleccionar una fecha válida");
                return;
            }
            this.isLoading = true;
            fetch('/api/v1/appointments', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    patient_id: this.selectedPatient.id,
                    specialty: this.specialty,
                    doctor_id: this.doctorId,
                    starts_at: this.date + "T" + this.time + ":05-05:00",
                    description: this.description,
                    status: this.selectedStatus
                })
            }).then(response => response.json())
                .then(data => {
                    console.log(data);
                    this.$emit('reload');
                    this.resetForm();
                    this.isLoading = false;
                })
                .catch(error => {
                    console.error('Error:', error);
                    alert("error ingresando la cita")
                })
        },

        fetchPatientWithDNI(dni) {
            fetch(`/api/v1/patients?document=${dni}`)
                .then(response => response.json())
                .then(data => {
                    this.patients = data.data;
                });
        },
        updateAppointment() {
            this.isLoading = true;
            fetch("/api/v1/appointments/" + this.id, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    patient_id: this.selectedPatient.id,
                    specialty: this.specialty,
                    doctor_id: this.doctorId,
                    starts_at: this.date + "T" + this.time + ":05-05:00",
                    description: this.description,
                    status: this.selectedStatus
                })
            }).then(response => response.json())
                .then(data => {
                    this.$emit('reload');
                    this.isLoading = false;
                })
                .catch(error => {
                    console.error('Error:', error);
                    alert("error actualizando la cita")
                })
        },
        fetchDoctors() {
            fetch("/api/v1/doctors")
                .then(response => response.json())
                .then(data => {
                    this.doctors = data.data;
                })
                .catch(error => {
                    console.error('Error:', error);
                    alert("error obteniendo los doctores")
                })
        }
    },
    computed: {
        invalidDate(){
            return this.date < this.minDate;
        },
        filteredPatients() {
            let patients = this.patients.filter(patient => patient.document_number.includes(this.query))
            return patients.slice(0, 4);
        },
        filteredDoctors(){
            return this.doctors.filter(doctor => doctor.specialty == this.specialty)
        }
    },
    mounted() {
        this.fetchPatientWithDNI()
        if (this.savedAppointment) {
            this.isEdit = true;
            this.id = this.savedAppointment.id;
            this.specialty = this.savedAppointment.specialty;
            this.date = this.savedAppointment.starts_at.split("T")[0];
            this.time = this.savedAppointment.starts_at.split("T")[1].split(":")[0] + ":00";
            this.description = this.savedAppointment.description;
            this.selectedStatus = this.savedAppointment.status;
            this.doctorId = this.savedAppointment.doctor_id;
            this.selectedPatient = {
                    id: this.savedAppointment.patient.id,
                    first_name: this.savedAppointment.patient.first_name,
                    father_last_name: this.savedAppointment.patient.father_last_name,
                    document_number: this.savedAppointment.patient.document_number
                };
        }
        this.fetchDoctors();
    }
};
</script>

<style scoped>
/* Your component styles go here */
</style>
