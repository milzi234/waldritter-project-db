<script setup>
import { computed, watchEffect } from 'vue';
import MegaMenu from 'primevue/megamenu';
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
      }
    }
  }
`;

const { result } = useQuery(MENU_QUERY);

const items = computed(() => {
  if (!result.value) return [];

  const menuItems = result.value.menuItems;
  const itemMap = new Map(menuItems.map(item => [item.documentId, item]));
  const rootItem = menuItems.find(item => item.label === '$MAIN');
  if (!rootItem) return [];

  const resolvedItems = [];
  function resolveItem(item, depth = 0) {
    const copy = { ...itemMap.get(item.documentId) };
    if (copy.to) {
      copy.command = navigateTo(copy.to);
      delete copy.to;
    }
    if (copy.items) {
      copy.items = copy.items.map(subItem => resolveItem(subItem, depth + 1));
    }
    if (copy.items.length === 0) {
      delete copy.items;
    } else if (depth === 0) {
      let currentCol = [];
      let subitemCount = 0;
      const cols = [];
      copy.items.forEach(item => {
        subitemCount += item.items?.length || 1;
        if (subitemCount > 8) {
          cols.push(currentCol);
          currentCol = [];
          subitemCount = item.items?.length || 1;
        }
        currentCol.push(item);
      });
      if (currentCol.length > 0) {
        cols.push(currentCol);
      }
      copy.items = cols;
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
  <MegaMenu :model="items" orientation="horizontal" class="w-full flex" :pt="{
    root: { class: 'justify-between', style: { border: 'none', borderRadius: '0' } },
    menu: { class: 'ml-auto' },
  }">
    <template #start>
      <div>
        <img src="/logo.png" alt="Logo" class="h-8 ml-4" />
      </div>
    </template>
  </MegaMenu>
</template>
