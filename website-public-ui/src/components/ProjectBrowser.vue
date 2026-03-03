<template>
  <div>
    <!-- Tag Filter -->
    <div class="mb-10 section-panel">
      <p class="font-mono text-xs text-gray-500 mb-4">FILTER::THEMEN</p>
      <div class="space-y-4">
        <div v-for="category in categories" :key="category.id" class="space-y-2">
          <h4 class="font-display text-xs tracking-wider text-ritter-300 uppercase">
            {{ category.title }}
          </h4>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="tag in category.tags"
              :key="tag.id"
              @click="toggleTag(tag.title)"
              :class="[
                'inline-block px-3 py-1 text-xs font-mono tracking-wider uppercase rounded-full border transition-all duration-300 cursor-pointer',
                selectedTags.has(tag.title)
                  ? 'border-wald-400 bg-wald-500/20 text-wald-300 shadow-[0_0_10px_rgba(0,255,136,0.2)]'
                  : 'border-wald-500/30 text-gray-400 bg-wald-950/60 hover:border-wald-400/50 hover:text-wald-300'
              ]"
            >
              {{ tag.title }}
            </button>
          </div>
        </div>

        <div v-if="selectedTags.size > 0" class="flex items-center gap-3 pt-2">
          <button
            @click="clearAll"
            class="text-xs font-mono text-gray-500 hover:text-ritter-400 transition-colors cursor-pointer"
          >
            Filter zurücksetzen
          </button>
        </div>
      </div>
    </div>

    <!-- Active filters -->
    <div v-if="selectedTags.size > 0" class="mb-6 flex items-center gap-2">
      <span class="font-mono text-xs text-gray-500">Aktive Filter:</span>
      <span
        v-for="tag in selectedTags"
        :key="tag"
        class="inline-block px-2 py-0.5 text-xs font-mono text-wald-300 border border-wald-500/30 rounded-full bg-wald-950/60"
      >
        {{ tag }}
      </span>
    </div>

    <!-- Result count -->
    <div class="flex items-center gap-2 mb-6 text-xs font-mono text-gray-500">
      <span class="w-1 h-1 rounded-full bg-wald-400"></span>
      <span>{{ projects.length }} Projekte gefunden</span>
    </div>

    <!-- Loading indicator -->
    <div v-if="loading" class="text-center py-12">
      <p class="font-mono text-xs text-wald-400 animate-pulse">LOADING::PROJECTS...</p>
    </div>

    <!-- Project Grid -->
    <div v-else-if="projects.length > 0" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
      <a
        v-for="project in projects"
        :key="project.id"
        :href="`/projects/${project.id}`"
        class="group block section-panel p-0 overflow-hidden hover:border-wald-400/30 transition-all duration-500"
      >
        <div v-if="project.imageUrl" class="aspect-video">
          <img
            :src="project.imageUrl"
            :alt="project.title"
            class="w-full h-full object-cover"
            loading="lazy"
          />
        </div>
        <div v-else class="aspect-video bg-wald-950/80 flex items-center justify-center">
          <span class="font-display text-wald-700 text-xs tracking-[0.3em]">KEIN BILD</span>
        </div>

        <div class="p-5 space-y-3">
          <h3 class="font-display text-sm tracking-wider text-gray-100 group-hover:text-wald-300 transition-colors line-clamp-2">
            {{ project.title }}
          </h3>

          <div v-if="project.description" class="text-xs text-gray-400 line-clamp-2 font-body prose prose-invert prose-xs prose-p:m-0 prose-headings:m-0 max-w-none" v-html="renderMarkdown(project.description)"></div>

          <div v-if="project.tags.length > 0" class="flex flex-wrap gap-1.5">
            <span
              v-for="tag in project.tags.slice(0, 3)"
              :key="tag.id"
              class="inline-block px-2 py-0.5 text-[10px] font-mono tracking-wider uppercase border border-wald-500/30 rounded-full text-wald-300 bg-wald-950/60 hover:border-wald-400/50 hover:shadow-[0_0_10px_rgba(0,255,136,0.15)] transition-all duration-300"
            >
              {{ tag.title }}
            </span>
            <span v-if="project.tags.length > 3" class="text-[10px] text-gray-500 font-mono self-center">
              +{{ project.tags.length - 3 }}
            </span>
          </div>

          <div
            v-if="project.nextOccurrence?.startDate"
            class="flex items-center gap-2 text-[10px] font-mono text-ritter-400"
          >
            <span class="w-1 h-1 rounded-full bg-ritter-400 animate-pulse"></span>
            <span>Nächster Termin: {{ formatDate(project.nextOccurrence.startDate) }}</span>
          </div>
        </div>
      </a>
    </div>

    <!-- Empty state -->
    <div v-else class="text-center py-20">
      <p class="font-mono text-gray-500 text-sm">Keine Projekte gefunden.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { marked } from 'marked';

interface Tag {
  id: string;
  title: string;
  description: string | null;
  categoryId: number;
}

interface Category {
  id: string;
  title: string;
  tags: Tag[];
}

interface Occurrence {
  id: string;
  startDate: string | null;
  endDate: string | null;
}

interface Project {
  id: string;
  title: string;
  description: string | null;
  imageUrl: string | null;
  tags: Tag[];
  nextOccurrence: Occurrence | null;
}

const props = defineProps<{
  initialProjects: string;
  initialCategories: string;
  initialTags: string;
}>();

const apiUrl = import.meta.env.PUBLIC_API_URL || 'http://localhost:3000';

const categories = ref<Category[]>([]);
const projects = ref<Project[]>([]);
const selectedTags = ref<Set<string>>(new Set());
const loading = ref(false);

const PROJECT_FIELDS = `
  id
  title
  description
  imageUrl
  tags { id title description categoryId }
  nextOccurrence { id startDate endDate }
`;

onMounted(() => {
  // Hydrate from SSR data
  projects.value = JSON.parse(props.initialProjects);
  categories.value = JSON.parse(props.initialCategories);

  // Parse initial tags from URL
  const tagsParam = new URLSearchParams(window.location.search).get('tags');
  if (tagsParam) {
    tagsParam.split(',').forEach(t => selectedTags.value.add(decodeURIComponent(t)));
  }
});

function renderMarkdown(text: string): string {
  return marked.parse(text, { async: false }) as string;
}

function formatDate(dateStr: string | null): string {
  if (!dateStr) return '';
  return new Date(dateStr).toLocaleDateString('de-DE', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  });
}

async function fetchProjectsByTags(tags: string[]) {
  const args: string[] = [];
  if (tags.length) args.push(`tags: [${tags.map(t => `"${t}"`).join(', ')}]`);
  const argStr = args.length ? `(${args.join(', ')})` : '';

  const res = await fetch(`${apiUrl}/graphql`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      query: `{ projects${argStr} { ${PROJECT_FIELDS} } }`,
    }),
  });
  const json = await res.json();
  return json.data.projects as Project[];
}

async function toggleTag(title: string) {
  const next = new Set(selectedTags.value);
  if (next.has(title)) {
    next.delete(title);
  } else {
    next.add(title);
  }
  selectedTags.value = next;
  await updateProjects();
}

async function clearAll() {
  selectedTags.value = new Set();
  await updateProjects();
}

async function updateProjects() {
  const tags = Array.from(selectedTags.value);

  // Update URL without navigation
  const search = tags.length ? `?tags=${tags.map(encodeURIComponent).join(',')}` : '';
  history.replaceState(null, '', `/projects/${search}`);

  // Fetch and update
  loading.value = true;
  try {
    projects.value = await fetchProjectsByTags(tags);
  } catch (e) {
    console.error('Failed to fetch projects:', e);
  } finally {
    loading.value = false;
  }
}
</script>
