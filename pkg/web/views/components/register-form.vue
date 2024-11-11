<template>
    <div class="p-5">
        <div class="flex justify-between items-center">
            <span class="lg:text-4xl text-2xl font-bold">Nuevo registro</span>
            <span class="lg:text-xl text-sm font-bold">
                <p>Paciente:</p>
                <p>
                    {{ patient.first_name }}
                    {{ patient.father_last_name }}
                    {{ patient.mother_last_name }}
                </p>
            </span>
        </div>

        <hr class="border-1 border-gray-400 my-5">

        <form @submit.prevent="submit">
            <div class="flex">
                <fieldset>
                    <div class="flex space-x-2">
                        <button type="button" @click="selectedType='child'" class="btn btn-secondary">Niño/niña</button>
                        <button type="button" @click="selectedType='teenager'" class="btn btn-secondary">Adolescente</button>
                        <button type="button" @click="selectedType='young'" class="btn btn-secondary">Joven</button>
                        <button type="button" @click="selectedType='adult'" class="btn btn-secondary">Adulto</button>
                        <button type="button" @click="selectedType='senior'" class="btn btn-secondary">Adulto mayor</button>
                    </div>
                    <div class="bg-gray-100 border border-black mt-4">
                        <ul>
                            <li class="border-b border-gray-500 py-4 px-4" v-for="pdf in pdfList[selectedType]">
                                <div class="flex justify-between items-center">
                                    {{ pdf.name }}
                                    <a target="_blank" :href="pdf.url" class="btn btn-tertiary">Descargar</a>
                                </div>
                                <input name="file[]" type="file" accept=".pdf" class="font-normal text-sm mt-2 max-w-lg">
                            </li>
                        </ul>
                    </div>
                </fieldset>
                <fieldset class="md:ml-10 mt-5 md:mt-0 flex flex-col">
                    <!-- Section for file upload -->
                    <div class="mb-6">
                        <label class="form-label text-xl mb-2">Adicionales</label>
                        <input class="max-w-lg" name="file[]" type="file" accept=".pdf,.jpg,.png" multiple>
                    </div>
                    <div class="grow flex flex-col">
                        <label class="font-semibold mb-2">Observaciones adicionales</label>
                        <textarea
                            class="grow w-full border border-gray-400 rounded-md p-2 md:h-44 h-20 block font-normal"
                            placeholder="Escribe aquí..." v-model="body" required></textarea>
                    </div>
                </fieldset>
            </div>
            <div class="flex items-center space-x-5">
                <button class="btn btn-primary mt-10 px-10 lg:text-xl md:text-sm text-xs" type="submit">Guardar</button>
                <button type="button" class="btn btn-secondary mt-10 px-10 lg:text-xl md:text-sm text-xs"
                    @click="$emit('close')">Cancelar</button>
            </div>
        </form>
    </div>

</template>

<script>
export default {
    name: 'RegisterForm',
    props: {
        patient: {
            type: Object,
            required: true
        },
        record: {
            type: Object,
            required: false
        }
    },
    data() {
        return {
            selectedType: "young",
            pdfList: {
                child: [
                    {
                        name: 'Formato atención integral',
                        url: '/pdf/child-01-formato-atencion-integral.pdf',
                        attachment: '#'
                    },
                    {
                        name: 'Evaluación alimentación',
                        url: '/pdf/child-02-evaluacion-alimentacion-ninha-ninho.pdf',
                        attachment: '#'
                    },
                    {
                        name: 'Consulta',
                        url: '/pdf/child-03-consulta.pdf',
                        attachment: '#'
                    },
                ],
                teenager: [
                    {
                        name: 'Formato atención integral',
                        url: '/pdf/teenager-01-formato-atencion-integral-adolscente.pdf',
                        attachment: '#'
                    },
                    {
                        name: 'Formato atención integral - Datos generales',
                        url: '/pdf/teenager-01-formato-atencion-integral-datos-generales.pdf',
                        attachment: '#'
                    },
                    {
                        name: 'Seguimiento de riesgos',
                        url: '/pdf/teenager-02-seguimiento-riesgos.pdf',
                        attachment: '#'
                    },
                    {
                        name: 'Formato integral',
                        url: '/pdf/teenager-03-formato-integral-adolescente.pdf',
                        attachment: '#'
                    },
                    {
                        name: 'Formato integral',
                        url: '/pdf/teenager-04-formato-integral-adolescente.pdf',
                        attachment: '#'
                    },
                ],
                young: [
                    {
                        name: 'Formato atención integral',
                        url: '/pdf/young-01-formato-atencion-integral-joven.pdf',
                        attachment: '#'
                    },
                    {
                        name: 'Datos generales',
                        url: '/pdf/young-02-datos-generales.pdf',
                        attachment: '#'
                    },
                    {
                        name: 'Ultimo form',
                        url: '/pdf/young-03-ultimo-form.pdf',
                        attachment: '#'
                    },
                ],
                adult: [
                    {
                        name: 'Formato atención integral',
                        url: '/pdf/adult-01-formato-atencion-integral-datos-generales.pdf',
                        attachment: '#'
                    },
                    {
                        name: 'Consulta',
                        url: '/pdf/adult-02-formato-atencion-integral-datos-generales.pdf',
                        attachment: '#'
                    },
                ],
                senior: [
                    {
                        name: 'Formato atención integral',
                        url: '/pdf/senior-01-formato-atencion-integral-adulto-mayor.pdf',
                        attachment: '#'
                    },
                ]
            },
            body: "",
        }
    },
    methods: {
        submit(){
            this.sendForm()
        },
        sendForm(){
            const formData = new FormData()
            var fileInputs = document.querySelectorAll("input[type='file']");
            fileInputs.forEach(file => {
                if (file.files.length > 0){
                    for (let i = 0; i < file.files.length; i++) {
                        formData.append('file[]', file.files[i])
                    }
                }
            })

            formData.append('body', this.body)

            fetch(`/api/v1/patients/${this.patient.id}/records`, {
                method: 'POST',
                body: formData
            })
                .then(r => r.json())
                .then(data => {
                    if (data.status == 400){
                        alert("1004: Error al guardar el registro, tiene que subir un archivo")
                        return
                    }
                    console.log(data)
                    this.$emit('reload')
                })
                .catch(err => {
                    alert("1003: Error al guardar el registro")
                    console.error(err)
                })
        }
    },
    mounted() {
        if (this.record) {
            console.log(this.record);
        }
    }
};
</script>

<style scoped>
/*Your component's CSS styles go here */
</style>
