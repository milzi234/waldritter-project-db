import type { Struct, Schema } from '@strapi/strapi';

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

export interface ProjectsProjectFilter extends Struct.ComponentSchema {
  collectionName: 'components_projects_project_filters';
  info: {
    displayName: 'ProjectFilter';
    icon: 'database';
  };
  attributes: {
    selectedTags: Schema.Attribute.JSON;
    hiddenCategories: Schema.Attribute.JSON;
  };
}

export interface NavigationLink extends Struct.ComponentSchema {
  collectionName: 'components_navigation_links';
  info: {
    displayName: 'link';
    icon: 'link';
  };
  attributes: {
    label: Schema.Attribute.String;
    url: Schema.Attribute.String;
    external: Schema.Attribute.Boolean;
    svgIcon: Schema.Attribute.Text;
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

export interface PageHighlight extends Struct.ComponentSchema {
  collectionName: 'components_page_highlights';
  info: {
    displayName: 'Highlight';
    icon: 'emotionHappy';
    description: '';
  };
  attributes: {
    title: Schema.Attribute.String;
    description: Schema.Attribute.Text;
    image: Schema.Attribute.Media<'images' | 'files' | 'videos'>;
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
    projectFilter: Schema.Attribute.Component<'projects.project-filter', false>;
    projectsToHighlight: Schema.Attribute.Integer;
    highlights: Schema.Attribute.Component<'page.highlight', true>;
  };
}

export interface PageHero extends Struct.ComponentSchema {
  collectionName: 'components_page_heroes';
  info: {
    displayName: 'Hero';
    description: '';
  };
  attributes: {
    title: Schema.Attribute.String;
    description: Schema.Attribute.Text;
    image: Schema.Attribute.Media<'images' | 'files' | 'videos'>;
    more_link: Schema.Attribute.String;
  };
}

export interface PageContact extends Struct.ComponentSchema {
  collectionName: 'components_page_contacts';
  info: {
    displayName: 'Contact';
    description: '';
  };
  attributes: {
    name: Schema.Attribute.String;
    email: Schema.Attribute.String;
    image: Schema.Attribute.Media<'images' | 'files'>;
    social_links: Schema.Attribute.Component<'navigation.link', true>;
  };
}

declare module '@strapi/strapi' {
  export module Public {
    export interface ComponentSchemas {
      'projects.project-search': ProjectsProjectSearch;
      'projects.project-filter': ProjectsProjectFilter;
      'navigation.link': NavigationLink;
      'page.markdown': PageMarkdown;
      'page.highlight': PageHighlight;
      'page.highlight-reel': PageHighlightReel;
      'page.hero': PageHero;
      'page.contact': PageContact;
    }
  }
}
