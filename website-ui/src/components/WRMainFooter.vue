<script setup>
import { computed } from 'vue';
import { useQuery } from '@vue/apollo-composable';
import gql from 'graphql-tag';
import WRFooterLinks from './WRFooterLinks.vue';
import WRNewsletterCTA from './WRNewsletterCTA.vue';
import WRShortDescription from './WRShortDescription.vue';
import WRSocials from './WRSocials.vue';

const FOOTER_QUERY = gql`
  query Footer {
    footer {
      shortDescriptionTitle
      shortDescriptionText
      links_1 {
        label
        url
      }
      links_2 {
        label
        url
      }
      social_links {
        label
        url
        svgIcon
      }
    }
  }
`;

const { result } = useQuery(FOOTER_QUERY);

const footerData = computed(() => result.value?.footer || {});

const shortDescriptionTitle = computed(() => footerData.value.shortDescriptionTitle || '');
const shortDescriptionText = computed(() => footerData.value.shortDescriptionText || '');
const footerLinksColumn1 = computed(() => footerData.value.links_1 || []);
const footerLinksColumn2 = computed(() => footerData.value.links_2 || []);
const socialLinks = computed(() => footerData.value.social_links || []);

const currentYear = new Date().getFullYear();
</script>

<template>
  <div class="bg-green-800 text-white p-4 md:p-16 lg:p-24">
    <div class="container mx-auto">
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8 lg:gap-16 mb-8">
        <WRShortDescription :title="shortDescriptionTitle" :description="shortDescriptionText" />
        <WRFooterLinks :links="footerLinksColumn1" class="mt-6 lg:mt-10" />
        <WRFooterLinks :links="footerLinksColumn2" class="mt-6 lg:mt-10" />
      </div>
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8 lg:gap-16 items-center">
        <div class="flex flex-col">
          <WRSocials :socialLinks="socialLinks" />
          <div class="text-sm mt-4 text-[10px]">
            &copy; {{ currentYear }} Waldritter. Alle Rechte vorbehalten.
          </div>
        </div>
        <div class="hidden lg:block"></div>
        <WRNewsletterCTA class="mt-8 lg:mt-0" />
      </div>
    </div>
  </div>
</template>
