# frozen_string_literal: true

module Types
  class QueryType < Types::BaseObject
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
