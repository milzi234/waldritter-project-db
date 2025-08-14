<script setup>
import { computed, watchEffect } from 'vue';
import Menubar from 'primevue/menubar';
import { useRouter } from 'vue-router';
import { useQuery } from '@vue/apollo-composable';
import gql from 'graphql-tag';

const router = useRouter();

function navigateTo(to) {
  return () => {
    router.push(to);
  }
}

const MENU_QUERY = gql`
  query Menu {
    menuItems(pagination: {limit: 1000}) {
      documentId
      label
      to
      items {
        documentId
        label
        to
        items {
          documentId
          label
          to
        }
      }
    }
  }
`;

const { result } = useQuery(MENU_QUERY, null, {
  clientId: 'pagesAPI'
});

const items = computed(() => {
  if (!result.value) {
    return [];
  }

  const menuItems = result.value.menuItems;
  
  // Build a map of ALL menu items (including nested ones)
  const itemMap = new Map();
  function addToMap(items) {
    if (!items) return;
    items.forEach(item => {
      itemMap.set(item.documentId, item);
      if (item.items) {
        addToMap(item.items);
      }
    });
  }
  addToMap(menuItems);
  
  const rootItem = menuItems.find(item => item.label === '$MAIN');
  if (!rootItem) {
    return [];
  }

  const resolvedItems = [];
  function resolveItem(item, depth = 0) {
    const mapItem = itemMap.get(item.documentId);
    const copy = { ...mapItem };
    
    // Process subitems first if they exist
    if (copy.items) {
      if (copy.items.length > 0) {
        copy.items = copy.items.map(subItem => resolveItem(subItem, depth + 1));
      }
      
      // Delete items property if it's empty to prevent chevron display
      if (!copy.items || copy.items.length === 0) {
        delete copy.items;
      }
    }
    
    // Add navigation command for items with 'to' URL
    // Only add command if the item doesn't have subitems
    if (copy.to && (!copy.items || copy.items.length === 0)) {
      copy.command = navigateTo(copy.to);
      delete copy.to;
    } else if (copy.to) {
      // Remove the 'to' URL for items with subitems
      delete copy.to;
    }
    
    // Delete all keys except label, command, and items
    const allowedKeys = ['label', 'command', 'items'];
    Object.keys(copy).forEach(key => {
      if (!allowedKeys.includes(key)) {
        delete copy[key];
      }
    });
    
    return copy;
  }
  rootItem.items.forEach(item => {
    resolvedItems.push(resolveItem(item));
  });

  return resolvedItems;
});

</script>

<template>
  <Menubar :model="items" class="w-full flex" :pt="{
    root: { class: 'justify-between', style: { border: 'none', borderRadius: '0' } },
    menu: { class: 'ml-auto' },
  }">
    <template #start>
      <div>
        <img src="/logo.png" alt="Logo" class="h-8 ml-4" />
      </div>
    </template>
  </Menubar>
</template>
