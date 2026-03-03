export interface Tag {
  id: string;
  title: string;
  description: string | null;
  categoryId: number;
  category?: Category;
}

export interface Category {
  id: string;
  title: string;
  description: string | null;
  tags: Tag[];
}

export interface Occurrence {
  id: string;
  startDate: string | null;
  endDate: string | null;
  eventId: number;
}

export interface Event {
  id: string;
  startDate: string | null;
  endDate: string | null;
  recurrenceType: string | null;
  title: string | null;
  description: string | null;
  projectId: number;
}

export interface Project {
  id: string;
  title: string;
  description: string | null;
  homepage: string | null;
  imageUrl: string | null;
  tags: Tag[];
  events: Event[];
  occurrences: Occurrence[];
  nextOccurrence: Occurrence | null;
  createdAt: string;
  updatedAt: string;
}
