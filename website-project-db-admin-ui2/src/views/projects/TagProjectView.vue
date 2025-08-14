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
<h1>Projekt Taggen</h1>
<h2>{{ title }}</h2>
<TagChooser :selected="projectTagIDs" @update="updateProjectTagIDs"/>
<br />
<button type="submit" class="btn btn-primary" @click.prevent="setTags">Speichern</button>&nbsp;
<button class="btn btn-secondary" @click.prevent="router.go(-1)">Zurück</button>

</template>

