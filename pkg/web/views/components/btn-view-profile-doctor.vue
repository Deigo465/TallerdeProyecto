<template>
    <div>
        <button @click="showModal = true" class="underline hover:text-blue-500">Ver Perfil</button>
        <modal :showModal="showModal" @close="showModal = false; user = $event">
            <section class="p-5">

                <!-- Header -->
                <div class="flex justify-between items-center mb-6">
                    <h1 class="lg:text-2xl text-xl font-bold">Perfil</h1>
                    <div class="md:flex text-primary underline font-bold md:space-x-4 md:items-center">
                        <btn-edit-doctor :savedDoctor="user" />
                        <button class="lg:text-base text-sm" @click="resetPassword(user.id)">Restablecer
                            contraseña</button>
                    </div>
                </div>

                <!-- nombre -->
                <h1 class="lg:text-3xl md:text-2xl font-bold mb-5"> {{ fullName }}</h1>

                <!-- datos -->
                <div class="flex justify-center">
                    <table class="w-full lg:mb-10 mb-0 justify-center">
                        <thead>
                            <tr class="font-bold text-center lg:text-base text-sm border-b border-black">
                                <th class="md:px-6 px-1 py-1">DNI</th>
                                <th class="md:px-6 px-1 py-1">F. Nac.</th>
                                <th class="md:px-6 px-1 py-1">Teléfono</th>
                                <th class="md:px-6 px-1 py-1 hidden lg:table-cell">Correo electró</th>
                                <th class="md:px-6 px-1 py-1 hidden lg:table-cell">CMP</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr class="text-center lg:text-base text-sm">
                                <td class="md:px-6 px-1 py-1">{{ user.document_number }}</td>
                                <td class="md:px-6 px-1 py-1">{{ (new Date(user.date_of_birth)).toLocaleDateString() }}
                                </td>
                                <td class="md:px-6 px-1 py-1">{{ user.phone }}</td>
                                <td class="md:px-6 px-1 py-1 hidden lg:table-cell">{{ user.email }}</td>
                                <td class="md:px-6 px-1 py-1 hidden lg:table-cell">{{ user.cmp }}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <!-- responsive de la tabla -->
                <div class="flex justify-center">
                    <table class="w-full md:mb-10 mb-5 justify-center">
                        <thead>
                            <tr class="font-bold text-center lg:text-base text-sm border-b border-black lg:hidden">
                                <th class="md:px-6 px-1 py-1">Correo electrónico</th>
                                <th class="md:px-6 px-1 py-1">CMP</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr class="text-center lg:hidden lg:text-base text-sm">
                                <td class="md:px-6 px-1 py-1">{{ user.email }}</td>
                                <td class="md:px-6 px-1 py-1">{{ user.cmp }}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <!-- Historial de citas -->
                <p class="lg:text-3xl md:text-2xl text-xl font-bold text-center mb-5">
                    Historial de citas
                </p>

                <table class="w-full md:mb-10 mb-5">
                    <thead>
                        <tr class="font-bold text-center lg:text-xl md:text-sm text-xs border-b border-black">
                            <th class="py-2 lg:px-0 md:px-5 px-1 text-left">Fecha</th>
                            <th class="py-2 lg:px-0 md:px-5 px-1 text-left">Paciente</th>
                            <th class="py-2 lg:px-0 md:px-5 px-1">Estado</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr class="border-b border-gray-400 text-center lg:text-base md:text-sm text-xs"
                            v-for="appointment in appointments" :key="appointment.id">
                            <td class="py-2 lg:px-0 md:px-5 px-1 text-left">{{ appointment.date }}</td>
                            <td class="py-2 lg:px-0 md:px-5 px-1 text-left">{{ appointment.patient }}</td>
                            <td class="py-2 lg:px-0 md:px-5 px-1 ">
                                <span class="inline-block btn rounded-2xl lg:text-sm text-xs lg:w-24 w-20"
                                    :class="statusColorBg[appointment.status]">
                                    {{ appointment.status }}
                                </span>
                            </td>
                        </tr>
                    </tbody>
                </table>

                <div class="flex items-center space-x-5">
                    <button class="btn btn-secondary mt-10 px-10 lg:text-xl md:text-sm text-xs"
                        @click="showModal = false">Atrás</button>
                </div>

            </section>
        </modal>
    </div>
</template>

<script>
export default {
    name: 'BtnViewProfileDoctor',
    props: {
        user: {
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
            return `${this.user.first_name} ${this.user.father_last_name} ${this.user.mother_last_name}`;
        }
    },
    methods: {
        resetPassword(id) {
            fetch(`/api/v1/doctors/${id}/reset-password`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                }})
                .then(response => response.json())
                .then(data => {
                    if (data.status == 200){
                        alert('Se ha restablecido la contraseña: ' + data.data.password);
                    } else {
                        alert('Error al restablecer la contraseña');
                    }
                })
        }
    }
}
</script>

<style scoped>
/* Your component's CSS code goes here */
</style>
