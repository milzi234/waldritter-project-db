# frozen_string_literal: true

module Types
  class QueryType < Types::BaseObject
    # Projects
    field :projects, [Types::ProjectType], null: false, description: "Get all projects with optional filtering" do
      argument :tag_ids, [ID], required: false, description: "Filter by tag IDs (AND logic)"
      argument :tags, [String], required: false, description: "Filter by tag titles (case-insensitive, AND logic)"
      argument :limit, Integer, required: false, description: "Maximum number of projects to return"
      argument :offset, Integer, required: false, description: "Number of projects to skip"
    end

    def projects(tag_ids: nil, tags: nil, limit: nil, offset: nil)
      result = Project.includes(:tags).all

      # Filter by tag IDs if provided (AND logic - project must have ALL specified tags)
      if tag_ids.present?
        tag_ids.each do |tag_id|
          result = result.joins(:tags).where(tags: { id: tag_id })
        end
        result = result.distinct
      end

      # Filter by tag titles if provided (case-insensitive, AND logic)
      if tags.present?
        tags.each do |tag_title|
          result = result.joins(:tags).where("LOWER(tags.title) = LOWER(?)", tag_title)
        end
        result = result.distinct
      end

      # Apply offset
      result = result.offset(offset) if offset.present?

      # Apply limit
      result = result.limit(limit) if limit.present?

      result
    end

    # Project count (for pagination)
    field :projects_count, Integer, null: false, description: "Get total count of projects matching filters" do
      argument :tag_ids, [ID], required: false, description: "Filter by tag IDs (AND logic)"
      argument :tags, [String], required: false, description: "Filter by tag titles (case-insensitive, AND logic)"
    end

    def projects_count(tag_ids: nil, tags: nil)
      result = Project.all

      if tag_ids.present?
        tag_ids.each do |tag_id|
          result = result.joins(:tags).where(tags: { id: tag_id })
        end
        result = result.distinct
      end

      if tags.present?
        tags.each do |tag_title|
          result = result.joins(:tags).where("LOWER(tags.title) = LOWER(?)", tag_title)
        end
        result = result.distinct
      end

      result.count
    end

    # Single project by ID
    field :project, Types::ProjectType, null: true, description: "Get a single project by ID" do
      argument :id, ID, required: true
    end

    def project(id:)
      Project.includes(:tags).find_by(id: id)
    end

    # Categories
    field :categories, [Types::CategoryType], null: false, description: "Get all categories"
    def categories
      Category.all
    end

    field :tags, [Types::TagType], null: false, description: "Get all tags"
    def tags
      Tag.all
    end

    # Search
    field :search, Types::SearchResultType, null: false, description: "Search for projects, events, and occurrences" do
      argument :start_date, GraphQL::Types::ISO8601DateTime, required: true
      argument :end_date, GraphQL::Types::ISO8601DateTime, required: true
      argument :tags, [ID], required: false
    end

    def search(start_date:, end_date:, tags: nil)
      search_params = {
        start_date: start_date,
        end_date: end_date,
        tags: tags
      }
      results = Search.search(search_params)

      {
        projects: results[:projects].values,
        events: results[:events].values,
        occurrences: results[:occurrences]
      }
    end
  end
end
