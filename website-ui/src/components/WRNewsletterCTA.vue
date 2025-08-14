<script setup>
import Button from 'primevue/button';

const mailchimpUrl = 'https://waldritter.us13.list-manage.com/subscribe?u=8cc2e000e5dc017f0476b002c&id=f01b223522';

const openMailchimp = () => {
  window.open(mailchimpUrl, '_blank', 'noopener,noreferrer');
};
</script>

<template>
  <div class="flex flex-col space-y-6">
    <h3 class="text-lg font-heading text-white">NEWSLETTER</h3>
    <p class="text-sm text-white/80">Erhalte alle aktuellen Neuigkeiten und Informationen zu unseren Projekten.</p>
    <Button 
      label="Zum Newsletter anmelden" 
      @click="openMailchimp" 
      class="p-button-lg w-full"
      icon="pi pi-external-link"
      iconPos="right"
    />
  </div>
</template>
