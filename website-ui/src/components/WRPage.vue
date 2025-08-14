<script setup>
import { ref, markRaw, defineProps, computed } from 'vue';
import { useQuery } from '@vue/apollo-composable';
import gql from 'graphql-tag';
import WRContact from './WRContact.vue';
import WRHero from './WRHero.vue';
import WRRichText from './WRRichText.vue';
import WRProjectSearch from './WRProjectSearch.vue';
import WRHighlightReel from './WRHighlightReel.vue';

const props = defineProps({
  documentId: {
    type: String,
    required: true
  }
});

const LOAD_PAGE_QUERY = gql`
  query LoadPage($documentId: ID!) {
    page(documentId: $documentId) {
      content {
        ... on ComponentPageMarkdown {
          id
          markdown
          __typename
        }
        ... on ComponentPageHero {
          id
          title
          description
          more_link
          image {
            caption
            url
          }
          __typename
        }
        ... on ComponentPageContact {
          id
          __typename
          name
          email
          image {
            caption
            url
          }
          social_links {
            external
            label
          }
        }
        ... on ComponentProjectsProjectSearch {
          filter {
            hiddenCategories
            selectedTags
          }
        }
        ... on ComponentPageHighlightReel {
          highlights {
            description
            image {
              url
              caption
            }
          }
          projectsToHighlight
          projectFilter {
            hiddenCategories
            selectedTags
          }
        }
      }
    }
  }
`;

const { result } = useQuery(LOAD_PAGE_QUERY, () => ({
  documentId: props.documentId
}), {clientId: "pagesAPI"});

const mapComponentType = (typename) => {
  switch (typename) {
    case 'ComponentPageMarkdown':
      return markRaw(WRRichText);
    case 'ComponentPageHero':
      return markRaw(WRHero);
    case 'ComponentPageContact':
      return markRaw(WRContact);
    case 'ComponentProjectsProjectSearch':
      return markRaw(WRProjectSearch);
    case 'ComponentPageHighlightReel':
      return markRaw(WRHighlightReel);
    default:
      return null;
  }
};

const mapComponentProps = (component) => {
  switch (component.__typename) {
    case 'ComponentPageMarkdown':
      return { markdown: component.markdown };
    case 'ComponentPageHero':
      return {
        title: component.title,
        description: component.description,
        imageSrc: "http://localhost:1337" + component.image?.url,
        imageAlt: component.imageAlt || component.image?.caption,
        moreLink: component.more_link
      };
    case 'ComponentPageContact':
      return {
        name: component.name,
        email: component.email,
        imageSrc: "http://localhost:1337" + component.image?.url,
        imageAlt: component.image?.caption,
        socials: component.social_links.map(link => ({
          name: link.label,
          url: link.external ? link.label : `https://${link.label.toLowerCase()}.com/waldritter`
        }))
      };
    case 'ComponentProjectsProjectSearch':
      return {
        filter: component.filter
      };
    case 'ComponentPageHighlightReel':
      return {
        highlights: component.highlights.map(highlight => ({
          ...highlight,
          image: {
            ...highlight.image,
            url: "http://localhost:1337" + highlight.image.url
          }
        })),
        projectsToHighlight: component.projectsToHighlight,
        projectFilter: component.projectFilter
      };
    default:
      return {};
  }
};

const colors = [
  { bg: 'bg-emerald-700', text: 'text-white', hover: 'hover:bg-emerald-800' },
  { bg: 'bg-white', text: 'text-emerald-900', hover: 'hover:bg-emerald-100' },
  { bg: 'bg-emerald-900', text: 'text-white', hover: 'hover:bg-emerald-800' },
  { bg: 'bg-green-100', text: 'text-emerald-900', hover: 'hover:bg-green-200' },
  { bg: 'bg-emerald-800', text: 'text-white', hover: 'hover:bg-emerald-700' },
];

const contactBgColor = 'bg-emerald-100';
const contactTextColor = 'text-emerald-900';

const pageStructure = computed(() => {
  if (result.value && result.value.page) {
    let heroIndex = 0;
    return result.value.page.content.map((component, index) => {
      const colorSet = colors[index % colors.length];
      const mappedComponent = {
        type: mapComponentType(component.__typename),
        props: mapComponentProps(component),
        bgColor: component.__typename === 'ComponentPageContact' ? contactBgColor : 
                 component.__typename === 'ComponentPageMarkdown' ? 'bg-white' : colorSet.bg,
        textColor: component.__typename === 'ComponentPageContact' ? contactTextColor : 
                   component.__typename === 'ComponentPageMarkdown' ? 'text-emerald-900' : colorSet.text,
      };
      
      if (component.__typename === 'ComponentPageHero') {
        mappedComponent.props.imageOnLeft = heroIndex % 2 === 0;
        mappedComponent.props.buttonColor = colorSet.text.replace('text-', '');
        mappedComponent.props.buttonHoverColor = colorSet.hover.replace('hover:', '');
        heroIndex++;
      }
      
      return mappedComponent;
    });
  }
  return [];
});

const componentClasses = (component) => {
  if (component.type === markRaw(WRHero) || component.type === markRaw(WRRichText)) {
    return 'w-full';
  }
  return 'container mx-auto px-4';
};
</script>

<template>
  <div>
    <template v-for="(component, index) in pageStructure" :key="index">
      <div :class="[component.bgColor, component.textColor, 'w-full']">
        <div :class="componentClasses(component)">
          <component 
            :is="component.type" 
            v-bind="component.props"
          />
        </div>
      </div>
    </template>
  </div>
</template>