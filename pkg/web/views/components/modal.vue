<template lang="html">
    <div>
          <transition :duration="{enter: 200, leave: 400}">
            <div v-show="showModal">
            <div class="fixed z-70 inset-0 overflow-y-auto">
                <div class="flex items-center justify-center min-h-screen sm:pt-4 px-2 sm:px-4 sm:pb-20 text-center">
                    <!--
                    Background overlay, show/hide based on modal state.
                    -->
                    <transition
                        enter-active-class="ease-out"
                        enter-class="opacity-0"
                        enter-to-class="opacity-100"
                        leave-active-class="delay-200 ease-in"
                        leave-class="opacity-100"
                        leave-to-class="opacity-0"
                        >
                        <div class="fixed inset-0 transition-opacity duration-200" aria-hidden="true" v-show="showModal" @click="closeModal()">
                            <div class="absolute inset-0 bg-gray-800 opacity-90"></div>
                        </div>
                    </transition>

                    <!-- This element is to trick the browser into centering the modal contents. -->
                    <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
                    <!--
                    Modal panel, show/hide based on modal state.
                    -->
                    <transition
                    enter-active-class="ease-out"
                    enter-class="delay-150 opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                    enter-to-class="opacity-100 translate-y-0 sm:scale-100"
                    leave-active-class="ease-in"
                    leave-class="opacity-100 translate-y-0 sm:scale-100"
                    leave-to-class="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                    >
                        <div class="inline-block align-bottom bg-white rounded-lg sm:px-4 sm:pt-5 pb-4 text-left shadow-xl transform duration-200 transition-all sm:align-middle sm:p-6"
                             :class="classSize" role="dialog" aria-modal="true" aria-labelledby="modal-headline" v-show="showModal">
                             <slot></slot>
                        </div>
                    </transition>
                </div>
            </div>
            </div>
        </transition>
    </div>
</template>
<script>
export default {
    props: {
      showModal:{},
      classSize:{
        type: String,
      }
    },
    methods: {
        closeModal(){
            this.$emit('close');
        }
    },
}
</script>
<style lang="css">
    .z-70{
    z-index: 70;
    }

</style>
