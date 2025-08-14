<script setup>
import { ref } from 'vue';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import FloatLabel from 'primevue/floatlabel';

const email = ref('');

const subscribe = () => {
  // TODO: Implement newsletter subscription logic
  console.log('Subscribing with email:', email.value);
  email.value = '';
};
</script>

<template>
  <div class="flex flex-col space-y-6">
    <h3 class="text-lg font-heading text-white">NEWSLETTER</h3>

      <InputText id="email" v-model="email" class="w-full" placeholder="E-Mail eingeben" />

    <Button label="Eintragen" @click="subscribe" class="p-button-lg w-full mt-2" />
  </div>
</template>
