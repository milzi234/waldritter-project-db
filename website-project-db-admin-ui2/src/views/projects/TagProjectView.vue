<script setup>
  import { ref, watchEffect } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import useProjectAPI from '../../api/Project';
  import TagViewer from '../../components/TagViewer.vue'
  import TagChooser from '../../components/TagChooser.vue'
  const route = useRoute()
  const router = useRouter()

  const title = ref('')
  const projectTagIDs = ref([])

  const projectAPI = useProjectAPI();

  watchEffect(async () => {
    const project = await projectAPI.get(route.params.id);
    title.value = project.title
    var projectTags = await projectAPI.getProjectTags(route.params.id)
    projectTagIDs.value = projectTags.map(projectTag => projectTag.id)
  });

  const updateProjectTagIDs = (newProjectTagIDs) => {
    console.log("updateProjectTagIDs", newProjectTagIDs)
    projectTagIDs.value = newProjectTagIDs
  }

  const setTags = async () => {
    await projectAPI.setProjectTags(route.params.id, projectTagIDs.value);
    router.push(`/projects/${route.params.id}`);
  }

</script>

<template>
<h1 class="text-2xl font-display font-bold text-wald-300 mb-2">Projekt Taggen</h1>
<h2 class="text-lg text-gray-400 mb-6">{{ title }}</h2>
<div class="section-panel mb-6">
  <TagChooser :selected="projectTagIDs" @update="updateProjectTagIDs"/>
</div>
<div class="flex gap-3">
  <button type="submit" class="btn-cyber" @click.prevent="setTags">Speichern</button>
  <button class="btn-cyber-secondary" @click.prevent="router.go(-1)">Zurück</button>
</div>
</template>
