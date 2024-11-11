<template>
    <div class="p-5">
        <h1 class="lg:text-4xl text-2xl font-bold mb-5">
            <span v-if="!isEdit">Agregar doctor</span>
            <span v-else>Editar doctor</span>
        </h1>
        <p class="lg:text-2xl md:text-xl text-xl font-bold">Datos personales</p>
        <p class="lg:text-xl md:text-sm text-xs">
            Ingresa los datos del doctor, como aparecen en su documento de identidad
        </p>
        <form @submit.prevent="submitForm">
            <div class="grid md:grid-cols-3 grid-cols-2 gap-4">
                <div>
                    <label class="form-label lg:text-xl text-sm" for="grid-first-name">
                        Nombres
                    </label>
                    <input v-model="name" name="nombres" class="form-input"
                        :class="invalidName ? 'border-red-500' : 'border-gray-200'" type="text"
                        placeholder="Ingresa tus nombres" required>
                    <p class="text-red-500 text-xs italic" v-show="invalidName">
                        Por favor rellena este campo.
                    </p>
                </div>

                <div>
                    <label class="form-label lg:text-xl text-sm" for="grid-father-last-name">
                        Apellido Paterno
                    </label>
                    <input v-model="fatherName" name="apellidoPapa" class="form-input"
                        :class="invalidFatherName ? 'border-red-500' : 'border-gray-200'" type="text"
                        placeholder="Ingresa tu apellido paterno" required>
                    <p class="text-red-500 text-xs italic" v-show="invalidFatherName">Por favor rellena
                        este
                        campo.</p>
                </div>

                <div>
                    <label class="form-label lg:text-xl text-sm" for="grid-mother-last-name">
                        Apellido Materno
                    </label>
                    <input v-model="motherName" name="apellidoMama" class="form-input"
                        :class="invalidMotherName ? 'border-red-500' : 'border-gray-200'" type="text"
                        placeholder="Ingresa tu apellido materno" required>
                    <p class="text-red-500 text-xs italic" v-show="invalidMotherName">Por favor rellena
                        este campo.</p>
                </div>

                <div>
                    <label class="form-label lg:text-xl text-sm" for="grid-date">
                        Fecha de nacimiento
                    </label>
                    <input v-model="dob" name="fecha" class="form-input"
                        :class="invalidDob ? 'border-red-500' : 'border-gray-200'" type="date"
                        placeholder="Ingresa tu fecha de nacimiento" required min="1824-01-01" :max="new Date().toISOString().split('T')[0]">
                    <p class="text-red-500 text-xs italic" v-show="invalidDob">
                        Ingresa una fecha de nacimiento válida.
                    </p>
                </div>

                <div>
                    <label class="form-label lg:text-xl text-sm" for="grid-identity-number">
                        Doc. Identidad
                    </label>
                    <input v-model="documentNumber" name="docIdentidad" class="form-input"
                        :class="invalidDocumentNumber ? 'border-red-500' : 'border-gray-200'"
                        type="text" 
                        @keydown="event => String(event.target.value).length == 8 && (event.key != 'Backspace' && event.key != 'Delete') ? event.preventDefault() : ''"
                        placeholder="Ingresar número de documento" required>
                    <p class="text-red-500 text-xs italic" v-show="invalidDocumentNumber">
                        Ingresa un número de documento válido.
                    </p>
                </div>

                <div>
                    <label class="form-label lg:text-xl text-sm" for="grid-gender">
                        Género
                    </label>
                    <select v-model.lazy="gender" name="genero" class="form-input"
                        :class="invalidGender ? 'border-red-500' : 'border-gray-200'" required>
                        <option value="">Selecciona tu género</option>
                        <option value="Female">Femenino</option>
                        <option value="Male">Masculino</option>
                        <option value="Other">Prefiero no decirlo</option>
                    </select>
                    <p class="text-red-500 text-xs italic" v-show="invalidGender">Por favor rellena este campo.
                    </p>
                </div>

                <div>
                    <label class="form-label lg:text-xl text-sm" for="grid-cmp-number">
                        CMP
                    </label>
                    <input v-model="cmp" name="CMP" class="form-input"
                        :class="invalidCMP ? 'border-red-500' : 'border-gray-200'" type="number"
                        placeholder="Ingresar número de CMP" required>
                    <p class="text-red-500 text-xs italic" v-show="invalidCMP">Por favor rellena
                        este
                        campo.</p>
                </div>

                <div class="relative mt-2 flex flex-col">
                    <label class="text-sm">Especialidad</label>
                    <select id="specialty" class="border-black bg-white py-2 p-2 md:w-48 border "
                        required v-model="specialty" :disabled="isEdit">
                        <option v-for="specialty in specialties" :value="specialty">{{ specialty }}</option>
                    </select>
                </div>

                <div v-if="!isEdit">
                    <label for="" class="text-sm">Contraseña por defecto</label>
                    <div class="text-red-600 border-black bg-white py-2 p-2 md:w-48 border">
                        123456
                    </div>
                </div> 
            </div>

            <hr class="border-1.5 border-black my-5">

            <p class="lg:text-2xl text-xl font-bold">Información de contacto</p>
            <p class="lg:text-xl md:text-sm text-xs">Ingresa al menos un correo y un telefono para poder contactar al
                cliente</p>

            <div class="flex mb-10">

                <div>
                    <label class="form-label lg:text-xl text-sm" for="grid-identity-number">
                        Correo electrónico
                    </label>
                    <input v-model="email" name="Correo" class="form-input border-black"
                        :class="invalidEmail ? 'border-red-500' : 'border-gray-200'" type="email"
                        placeholder="Ingresa tu correo electrónico" required>
                    <p class="text-red-500 text-xs italic" v-show="invalidEmail">Ingresa un correo válido.
                    </p>
                </div>

                <div class="ml-10">
                    <label class="form-label lg:text-xl text-sm" for="grid-identity-number">
                        Teléfono
                    </label>
                    <input v-model="phone" name="Telefono" class="form-input border-black"
                        :class="invalidPhone ? 'border-red-500' : 'border-gray-200'" type="number"
                        @keydown="event => String(event.target.value).length == 9 && (event.key != 'Backspace' && event.key != 'Delete') ? event.preventDefault() : ''"
                        placeholder="Ingresa tu telefono" required>
                    <p class="text-red-500 text-xs italic" v-show="invalidPhone"> Ingresa un teléfono válido.
                    </p>
                </div>

            </div>

            <div class="flex items-center space-x-5">
                <button class="btn btn-primary mt-10 px-10 lg:text-xl md:text-sm text-xs" type="submit">{{ isEdit ?
                    'Guardar' : 'Agregar' }}</button>
                <button class="btn btn-secondary mt-10 px-10 lg:text-xl md:text-sm text-xs"
                    @click.prevent="$emit('close')">Cancelar</button>
            </div>

        </form>
    </div>
</template>

<script>
export default {
    name: 'DoctorForm',
    props: {
        savedDoctor: {
            type: Object,
            required: false,
        }
    },
    data() {
        return {
            isEdit: false,
            id: null,

            name: "",
            validateName: false,

            fatherName: "",
            validateFatherName: false,

            motherName: "",
            validateMotherName: false, 

            dob: "",
            validateDob: false,

            documentNumber: "",
            validateDocumentNumber: false,

            gender: "",
            validateGender: false,

            cmp: "",
            validateCMP: false,

            specialty: "Cardiología",

            email: "",
            validateEmail: false,

            phone: "",
            validatePhone: false,
            specialties: ["Cardiología", "Dermatología", "Endocrinología", "Gastroenterología", "Geriatría", "Ginecología", "Hematología", "Infectología", "Medicina interna", "Nefrología", "Neumología", "Neurología", "Nutriología", "Oftalmología", "Oncología", "Pediatría", "Psiquiatría", "Reumatología", "Traumatología", "Urología"],
        };
    },
    computed: {
        invalidName() {
            return this.validateName && this.name.length === 0;
        },
        invalidFatherName() {
            return this.validateFatherName && this.fatherName.length === 0;
        },
        invalidMotherName() {
            return this.validateMotherName && this.motherName.length === 0;
        },
        invalidDob() {
            let validDob = false;
            if (this.dob.length > 0) {
                let dob = new Date(this.dob + "T23:00:00Z"); // assume time is 23:00:00 so that TZ doesn't change the date
                let today = new Date();
                validDob = dob < today;
                // check that the date is newer than 1824
                validDob = validDob && dob.getFullYear() >= 1824;
            }
            return this.validateDob && !validDob;
        },
        invalidDocumentNumber() {
            if ( this.validateDocumentNumber && !this.isNumeric(this.documentNumber)) {
                return true
            }
            return this.validateDocumentNumber && String(this.documentNumber).length !== 8
        },
        invalidGender() {
            return this.validateGender && this.gender.length === 0;
        },
        invalidCMP() {
            return this.validateCMP && String(this.cmp).length === 0;
        },
        invalidEmail() {
            // check if email is valid
            let validEmail = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(this.email);
            return this.validateEmail && !validEmail;
        },
        invalidPhone() {
            return this.validatePhone && String(this.phone).length !== 9 && isNaN(parseInt(this.phone))
        }
    },
    watch: {
        name: function (val) {
            this.validateName = true;
        },
        fatherName: function (val) {
            this.validateFatherName = true;
        },
        motherName: function (val) {
            this.validateMotherName = true;
        },
        dob: function (val) {
            this.validateDob = true;
        },
        documentNumber: function (val) {
            this.validateDocumentNumber = true;
        },
        gender: function (val) {
            this.validateGender = true;
        },
        cmp: function (val) {
            this.validateCMP = true;
        },
        email: function (val) {
            this.validateEmail = true;
        },
        phone: function (val) {
            this.validatePhone = true;
        },
    },
    methods: {
        isNumeric(str){
            return /^\d+$/.test(str);
        },
        submitForm() {
            if (this.invalidName || this.invalidFatherName || this.invalidMotherName || this.invalidDob || this.invalidDocumentNumber || this.invalidGender || this.invalidCMP || this.invalidEmail || this.invalidPhone) {
                alert("Por favor, rellena todos los campos correctamente");
                return;
            }
            if (this.id != null) {
                this.sendUpdateDoctor()
            } else {
                this.sendNewDoctor()
            }
        },
        sendUpdateDoctor() {
            fetch("/api/v1/doctors/" + this.id, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    first_name: this.name.toString(),
                    father_last_name: this.fatherName.toString(),
                    mother_last_name: this.motherName.toString(),
                    date_of_birth: this.dob.toString(),
                    document_number: this.documentNumber.toString(),
                    gender: this.gender.toString(),
                    cmp: this.cmp.toString(),
                    specialty: this.specialty.toString(),
                    contact_email: this.email.toString(),
                    phone: this.phone.toString(),
                }),
            })
                .then((response) => response.json())
                .then((data) => {
                    this.$emit("close", data.data);
                    window.location.reload();
                })
                .catch((error) => {
                    alert("1002: Error al actualizar datos del doctor");
                    console.error("Error:", error);
                });
        },
        sendNewDoctor() {
            fetch("/api/v1/doctors", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    first_name: this.name.toString(),
                    father_last_name: this.fatherName.toString(),
                    mother_last_name: this.motherName.toString(),
                    date_of_birth: this.dob.toString(),
                    document_number: this.documentNumber.toString(),
                    gender: this.gender.toString(),
                    cmp: this.cmp.toString(),
                    specialty: this.specialty.toString(),
                    contact_email: this.email.toString(),
                    phone: this.phone.toString()
                }),
            })
                .then((response) => response.json())
                .then((data) => {
                    if (data.status == 201) {
                        console.log("Success:", data);
                        this.$emit("close");
                        this.resetForm();
                    } else {
                        if (data.status == 400 && data.message == "doctor data doesn't match") {
                            alert("CMP y datos del doctor no coinciden, revisa los datos y vuelve a intentarlo")
                        } else if (data.status == 400 && data.message == "document number already exists") {
                            alert("Ya hay un doctor con este DNI")
                        }
                    }
                })
                .catch((error) => {
                    alert("1001: Error al crear nuevo doctor");
                    console.error("Error:", error);
                });
        },
        resetForm() {
            this.name = "";
            this.fatherName = "";
            this.motherName = "";
            this.dob = "";
            this.documentNumber = "";
            this.gender = "";
            this.cmp = "";
            this.specialty = "";
            this.email = "";
            this.phone = "";
        },
    },
    mounted() {
        if (this.savedDoctor) {
            this.isEdit = true;
            this.id = this.savedDoctor.id;
            this.name = this.savedDoctor.first_name;
            this.fatherName = this.savedDoctor.father_last_name;
            this.motherName = this.savedDoctor.mother_last_name;
            this.dob = this.savedDoctor.date_of_birth.split("T")[0];
            this.documentNumber = this.savedDoctor.document_number;
            this.gender = this.savedDoctor.gender;
            this.cmp = this.savedDoctor.cmp;
            this.specialty = this.savedDoctor.specialty;
            this.email = this.savedDoctor.email;
            this.phone = this.savedDoctor.phone;
        }
    }
};
</script>

<style scoped>
/*Your component's CSS styles go here */
</style>
