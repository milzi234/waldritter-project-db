<script setup>
import { ref, computed } from 'vue';
import Button from 'primevue/button';

const props = defineProps({
  title: {
    type: String,
    required: true
  },
  description: {
    type: String,
    required: true
  },
  imageSrc: {
    type: String,
    required: true
  },
  imageAlt: {
    type: String,
    default: 'Hero image'
  },
  imageOnLeft: {
    type: Boolean,
    default: false
  },
  buttonColor: {
    type: String,
    default: 'primary'
  },
  buttonHoverColor: {
    type: String,
    default: 'bg-primary-600'
  },
  moreLink: {
    type: String,
    default: ''
  }
});

const containerClasses = computed(() => `
  flex 
  flex-col 
  ${props.imageOnLeft ? 'md:flex-row' : 'md:flex-row-reverse'} 
  items-stretch 
  w-full
`);

const textContainerClasses = computed(() => `
  w-full 
  md:w-1/2 
  flex
  flex-col
  justify-center
  py-12
  px-12
  md:px-16
  lg:px-18
`);

const imageContainerClasses = computed(() => `
  w-full 
  md:w-1/2
  flex
  items-center
  justify-center
  min-h-[400px]
`);

const buttonClasses = computed(() => [
  'p-button-outlined',
  `p-button-${props.buttonColor}`,
  props.buttonHoverColor,
  'font-bold',
  'py-2',
  'px-4',
  'rounded',
  'transition-colors',
  'duration-200'
]);
</script>

<template>
  <div :class="containerClasses">
    <div :class="textContainerClasses">
      <h1 class="
        text-xl 
        md:text-2xl 
        lg:text-3xl 
        font-heading 
        mb-4
      ">
        {{ title }}
      </h1>
      <p class="text-base text-sm mb-6">
        {{ description }}
      </p>
      <Button 
        v-if="moreLink"
        label="Mehr erfahren" 
        :pt="{
          root: { class: buttonClasses }
        }"
        :href="moreLink"
      />
    </div>
    <div :class="imageContainerClasses">
      <img 
        :src="imageSrc" 
        :alt="imageAlt" 
        class="w-full h-full object-cover object-center"
        height="400"
      />
    </div>
  </div>
</template>