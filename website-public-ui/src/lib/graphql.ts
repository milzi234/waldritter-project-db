import type { Project, Category, Tag } from './types';

const API_URL = import.meta.env.API_URL || 'http://localhost:3000';

async function query<T>(gql: string, variables: Record<string, unknown> = {}): Promise<T> {
  const res = await fetch(`${API_URL}/graphql`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ query: gql, variables }),
  });

  if (!res.ok) {
    throw new Error(`GraphQL request failed: ${res.status}`);
  }

  const json = await res.json();
  if (json.errors) {
    throw new Error(json.errors.map((e: { message: string }) => e.message).join(', '));
  }
  return json.data;
}

const PROJECT_FIELDS = `
  id
  title
  description
  homepage
  imageUrl
  createdAt
  updatedAt
  tags { id title description categoryId }
  events { id startDate endDate recurrenceType }
  nextOccurrence { id startDate endDate }
`;

export async function fetchProjects(opts: {
  tags?: string[];
  limit?: number;
  offset?: number;
} = {}): Promise<Project[]> {
  const args: string[] = [];
  if (opts.tags?.length) args.push(`tags: [${opts.tags.map(t => `"${t}"`).join(', ')}]`);
  if (opts.limit) args.push(`limit: ${opts.limit}`);
  if (opts.offset) args.push(`offset: ${opts.offset}`);
  const argStr = args.length ? `(${args.join(', ')})` : '';

  const data = await query<{ projects: Project[] }>(`{
    projects${argStr} { ${PROJECT_FIELDS} }
  }`);
  return data.projects;
}

export async function fetchProject(id: string): Promise<Project | null> {
  const data = await query<{ project: Project | null }>(`
    query($id: ID!) {
      project(id: $id) {
        ${PROJECT_FIELDS}
        occurrences { id startDate endDate }
      }
    }
  `, { id });
  return data.project;
}

export async function fetchCategories(): Promise<Category[]> {
  const data = await query<{ categories: Category[] }>(`{
    categories {
      id title description
      tags { id title description categoryId }
    }
  }`);
  return data.categories;
}

export async function fetchTags(): Promise<Tag[]> {
  const data = await query<{ tags: Tag[] }>(`{
    tags { id title description categoryId }
  }`);
  return data.tags;
}
