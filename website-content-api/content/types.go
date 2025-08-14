package content

import "time"

// PageMetadataYAML represents the YAML frontmatter of a page
type PageMetadataYAML struct {
	DocumentID   string                 `yaml:"documentId"`
	URL          string                 `yaml:"url"`
	Title        string                 `yaml:"title"`
	Description  string                 `yaml:"description"`
	Author       string                 `yaml:"author"`
	PublishDate  *time.Time             `yaml:"publishDate"`
	LastModified *time.Time             `yaml:"lastModified"`
	Tags         []string               `yaml:"tags"`
	// Components field removed - components now parsed from markdown body
}

// ComponentDefinition represents a component in the YAML frontmatter
type ComponentDefinition struct {
	Type                string                       `yaml:"type"`
	ID                  string                       `yaml:"id"`
	Title               string                       `yaml:"title,omitempty"`
	Description         string                       `yaml:"description,omitempty"`
	Content             string                       `yaml:"content,omitempty"`
	Markdown            string                       `yaml:"markdown,omitempty"`
	MoreLink            string                       `yaml:"more_link,omitempty"`
	Name                string                       `yaml:"name,omitempty"`
	Email               string                       `yaml:"email,omitempty"`
	Image               *ImageDefinition             `yaml:"image,omitempty"`
	SocialLinks         []SocialLinkDefinition       `yaml:"social_links,omitempty"`
	Filter              *FilterDefinition            `yaml:"filter,omitempty"`
	Highlights          []HighlightDefinition        `yaml:"highlights,omitempty"`
	ProjectsToHighlight []string                     `yaml:"projectsToHighlight,omitempty"`
	ProjectFilter       *FilterDefinition            `yaml:"projectFilter,omitempty"`
}

// ImageDefinition represents an image in the YAML frontmatter
type ImageDefinition struct {
	URL     string `yaml:"url"`
	Caption string `yaml:"caption,omitempty"`
	Width   int    `yaml:"width,omitempty"`
	Height  int    `yaml:"height,omitempty"`
}

// SocialLinkDefinition represents a social link in the YAML frontmatter
type SocialLinkDefinition struct {
	Label    string `yaml:"label"`
	URL      string `yaml:"url"`
	External bool   `yaml:"external"`
	SVGIcon  string `yaml:"svgIcon,omitempty"`
}

// FilterDefinition represents a project filter in the YAML frontmatter
type FilterDefinition struct {
	HiddenCategories []string `yaml:"hiddenCategories,omitempty"`
	SelectedTags     []string `yaml:"selectedTags,omitempty"`
}

// HighlightDefinition represents a highlight in the YAML frontmatter
type HighlightDefinition struct {
	Description string           `yaml:"description"`
	Image       *ImageDefinition `yaml:"image"`
}

// MenuMetadata represents the YAML frontmatter of a menu file
type MenuMetadata struct {
	Type       string               `yaml:"type"`
	DocumentID string               `yaml:"documentId"`
	MenuItems  []MenuItemDefinition `yaml:"menuItems"`
}

// MenuItemDefinition represents a menu item in the YAML frontmatter
type MenuItemDefinition struct {
	DocumentID string               `yaml:"documentId"`
	Label      string               `yaml:"label"`
	To         string               `yaml:"to,omitempty"`
	Items      []MenuItemDefinition `yaml:"items,omitempty"`
	Order      int                  `yaml:"order,omitempty"`
}

// FooterMetadata represents the YAML frontmatter of a footer file
type FooterMetadata struct {
	Type                  string                       `yaml:"type"`
	DocumentID            string                       `yaml:"documentId"`
	ShortDescriptionTitle string                       `yaml:"shortDescriptionTitle"`
	ShortDescriptionText  string                       `yaml:"shortDescriptionText"`
	Links1                []FooterLinkDefinition       `yaml:"links_1"`
	Links2                []FooterLinkDefinition       `yaml:"links_2"`
	SocialLinks           []FooterSocialLinkDefinition `yaml:"social_links"`
}

// FooterLinkDefinition represents a footer link in the YAML frontmatter
type FooterLinkDefinition struct {
	Label string `yaml:"label"`
	URL   string `yaml:"url"`
}

// FooterSocialLinkDefinition represents a footer social link in the YAML frontmatter
type FooterSocialLinkDefinition struct {
	Label   string `yaml:"label"`
	URL     string `yaml:"url"`
	SVGIcon string `yaml:"svgIcon,omitempty"`
}

// GraphQL model types

// Page represents a page in the GraphQL schema
type Page struct {
	DocumentID string         `json:"documentId"`
	URL        string         `json:"url"`
	Content    []Component    `json:"content"`
	Metadata   *PageMetadata  `json:"metadata,omitempty"`
}

// PageMetadata represents page metadata in the GraphQL schema (renamed from PageMetadataGraphQL)
type PageMetadata struct {
	Title        string     `json:"title,omitempty"`
	Description  string     `json:"description,omitempty"`
	LastModified *time.Time `json:"lastModified,omitempty"`
	Author       string     `json:"author,omitempty"`
}

// Component types for GraphQL schema

// Component interface for GraphQL union type
type Component interface {
	IsComponent()
}

type ComponentPageMarkdown struct {
	ID       string `json:"id"`
	Markdown string `json:"markdown"`
}

func (ComponentPageMarkdown) IsComponent() {}

type ComponentPageHero struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MoreLink    string `json:"more_link,omitempty"`
	Image       *Image `json:"image,omitempty"`
}

func (ComponentPageHero) IsComponent() {}

type ComponentPageContact struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Email       string        `json:"email"`
	Image       *Image        `json:"image,omitempty"`
	SocialLinks []*SocialLink `json:"social_links"`
}

func (ComponentPageContact) IsComponent() {}

type ComponentProjectsProjectSearch struct {
	ID     string         `json:"id"`
	Filter *ProjectFilter `json:"filter,omitempty"`
}

func (ComponentProjectsProjectSearch) IsComponent() {}

type ComponentPageHighlightReel struct {
	ID                  string         `json:"id"`
	Highlights          []*Highlight   `json:"highlights"`
	ProjectsToHighlight []string       `json:"projectsToHighlight,omitempty"`
	ProjectFilter       *ProjectFilter `json:"projectFilter,omitempty"`
}

func (ComponentPageHighlightReel) IsComponent() {}

// Supporting types

type Image struct {
	URL     string `json:"url"`
	Caption string `json:"caption,omitempty"`
	Width   int    `json:"width,omitempty"`
	Height  int    `json:"height,omitempty"`
}

type SocialLink struct {
	Label    string `json:"label"`
	URL      string `json:"url"`
	External bool   `json:"external"`
	SVGIcon  string `json:"svgIcon,omitempty"`
}

type ProjectFilter struct {
	HiddenCategories []string `json:"hiddenCategories,omitempty"`
	SelectedTags     []string `json:"selectedTags,omitempty"`
}

type Highlight struct {
	Description string `json:"description"`
	Image       *Image `json:"image"`
}

// MenuItem represents a menu item in the GraphQL schema
type MenuItem struct {
	DocumentID string      `json:"documentId"`
	Label      string      `json:"label"`
	To         string      `json:"to,omitempty"`
	Items      []*MenuItem `json:"items,omitempty"`
	Order      int         `json:"order,omitempty"`
}

// Footer represents the footer in the GraphQL schema
type Footer struct {
	ShortDescriptionTitle string               `json:"shortDescriptionTitle"`
	ShortDescriptionText  string               `json:"shortDescriptionText"`
	Links1                []*FooterLink        `json:"links_1"`
	Links2                []*FooterLink        `json:"links_2"`
	SocialLinks           []*FooterSocialLink  `json:"social_links"`
}

// FooterLink represents a footer link in the GraphQL schema
type FooterLink struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// FooterSocialLink represents a footer social link in the GraphQL schema
type FooterSocialLink struct {
	Label   string `json:"label"`
	URL     string `json:"url"`
	SVGIcon string `json:"svgIcon,omitempty"`
}