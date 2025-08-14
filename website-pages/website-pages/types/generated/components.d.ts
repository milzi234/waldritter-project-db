import type { Schema, Struct } from '@strapi/strapi';

export interface NavigationLink extends Struct.ComponentSchema {
  collectionName: 'components_navigation_links';
  info: {
    displayName: 'link';
    icon: 'link';
  };
  attributes: {
    external: Schema.Attribute.Boolean;
    label: Schema.Attribute.String;
    svgIcon: Schema.Attribute.Text;
    url: Schema.Attribute.String;
  };
}

export interface PageContact extends Struct.ComponentSchema {
  collectionName: 'components_page_contacts';
  info: {
    description: '';
    displayName: 'Contact';
  };
  attributes: {
    email: Schema.Attribute.String;
    image: Schema.Attribute.Media<'images' | 'files'>;
    name: Schema.Attribute.String;
    social_links: Schema.Attribute.Component<'navigation.link', true>;
  };
}

export interface PageHero extends Struct.ComponentSchema {
  collectionName: 'components_page_heroes';
  info: {
    description: '';
    displayName: 'Hero';
  };
  attributes: {
    description: Schema.Attribute.Text;
    image: Schema.Attribute.Media<'images' | 'files' | 'videos'>;
    more_link: Schema.Attribute.String;
    title: Schema.Attribute.String;
  };
}

export interface PageHighlight extends Struct.ComponentSchema {
  collectionName: 'components_page_highlights';
  info: {
    description: '';
    displayName: 'Highlight';
    icon: 'emotionHappy';
  };
  attributes: {
    description: Schema.Attribute.Text;
    image: Schema.Attribute.Media<'images' | 'files' | 'videos'>;
    title: Schema.Attribute.String;
    to: Schema.Attribute.Component<'navigation.link', true>;
  };
}

export interface PageHighlightReel extends Struct.ComponentSchema {
  collectionName: 'components_page_highlight_reels';
  info: {
    displayName: 'HighlightReel';
    icon: 'plane';
  };
  attributes: {
    highlights: Schema.Attribute.Component<'page.highlight', true>;
    projectFilter: Schema.Attribute.Component<'projects.project-filter', false>;
    projectsToHighlight: Schema.Attribute.Integer;
  };
}

export interface PageMarkdown extends Struct.ComponentSchema {
  collectionName: 'components_page_markdowns';
  info: {
    displayName: 'Markdown';
    icon: 'layer';
  };
  attributes: {
    markdown: Schema.Attribute.RichText & Schema.Attribute.Required;
  };
}

export interface ProjectsProjectFilter extends Struct.ComponentSchema {
  collectionName: 'components_projects_project_filters';
  info: {
    displayName: 'ProjectFilter';
    icon: 'database';
  };
  attributes: {
    hiddenCategories: Schema.Attribute.JSON;
    selectedTags: Schema.Attribute.JSON;
  };
}

export interface ProjectsProjectSearch extends Struct.ComponentSchema {
  collectionName: 'components_projects_project_searches';
  info: {
    displayName: 'ProjectSearch';
    icon: 'search';
  };
  attributes: {
    filter: Schema.Attribute.Component<'projects.project-filter', false>;
  };
}

declare module '@strapi/strapi' {
  export module Public {
    export interface ComponentSchemas {
      'navigation.link': NavigationLink;
      'page.contact': PageContact;
      'page.hero': PageHero;
      'page.highlight': PageHighlight;
      'page.highlight-reel': PageHighlightReel;
      'page.markdown': PageMarkdown;
      'projects.project-filter': ProjectsProjectFilter;
      'projects.project-search': ProjectsProjectSearch;
    }
  }
}
