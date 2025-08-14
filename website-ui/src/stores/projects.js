import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { useQuery, provideApolloClient } from '@vue/apollo-composable';
import gql from 'graphql-tag';
import { apolloClientProjects } from '../apollo';

export const useProjectStore = defineStore('projects', () => {
  return provideApolloClient(apolloClientProjects)(() => {
    const today = new Date();
    const oneYearFromNow = new Date(today.getFullYear() + 1, today.getMonth(), today.getDate());

    const LOAD_PROJECTS_DATA = gql`
      query LoadProjectsData($startDate: ISO8601DateTime!, $endDate: ISO8601DateTime!) {
        categories {
          id
          title
        }
        
        tags {
          id
          categoryId
          title
        }
        
        search(
          startDate: $startDate
          endDate: $endDate
        ) {
          events {
            id
          }
          projects {
            id
            title
            description
            imageUrl
            tags {
              id
              title
              categoryId
            }
          }
          occurrences {
            id
            eventId
            startDate
            endDate
          }
        }
      }
    `;

    const lastWeek = new Date(today.getTime() - 7 * 24 * 60 * 60 * 1000);

    const { result } = useQuery(LOAD_PROJECTS_DATA, 
      () => ({
        startDate: lastWeek.toISOString(),
        endDate: oneYearFromNow.toISOString()
      })
    );

    console.log("result", result);

    const categories = computed(() => {
      if (!result.value) return [];
      return result.value.categories.map(category => ({
        id: category.id,
        title: category.title,
        tags: result.value.tags
          .filter(tag => tag.categoryId === parseInt(category.id))
          .map(tag => ({ id: tag.id, title: tag.title }))
      }));
    });

    const projects = computed(() => {
      console.log(result.value);
      if (!result.value) return [];
      return result.value.search.projects.map(project => {
        const projectOccurrences = result.value.search.occurrences.filter(
          occurrence => occurrence.eventId === parseInt(project.id)
        );
        return {
          id: project.id,
          title: project.title,
          startDateTime: projectOccurrences[0]?.startDate || null,
          endDateTime: projectOccurrences[0]?.endDate || null,
          description: project.description,
          image: { url: project.imageUrl || '/fallback.webp', caption: project.title },
          occurrences: projectOccurrences.map(occurrence => ({
            startDateTime: occurrence.startDate,
            endDateTime: occurrence.endDate
          })),
          categories: project.tags.reduce((acc, tag) => {
            const category = result.value.categories.find(c => c.id === tag.categoryId.toString());
            if (category) {
              const existingCategory = acc.find(c => c.id === category.id);
              if (existingCategory) {
                existingCategory.tags.push({ id: tag.id, title: tag.title });
              } else {
                acc.push({
                  id: category.id,
                  title: category.title,
                  tags: [{ id: tag.id, title: tag.title }]
                });
              }
            }
            return acc;
          }, [])
        };
      });
    });

    const createFilter = () => {
      const filter = ref([]);
      const updateFilter = (newFilter) => {
        filter.value = newFilter;
      };
      const filteredProjects = computed(() => {
        if (filter.value.length === 0) return projects.value;
        return projects.value.filter(project => 
          project.categories.some(category => 
            category.tags.some(tag => filter.value.includes(tag.id))
          )
        );
      });
      return { filter, updateFilter, filteredProjects };
    };

    const defaultFilter = createFilter();

    return {
      categories,
      projects,
      createFilter,
      filter: defaultFilter.filter,
      filteredProjects: defaultFilter.filteredProjects,
      updateFilter: defaultFilter.updateFilter
    };
  });
});
