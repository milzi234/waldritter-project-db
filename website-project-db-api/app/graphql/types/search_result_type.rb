# frozen_string_literal: true

module Types
  class SearchResultType < Types::BaseObject
    field :projects, [Types::ProjectType], null: false
    field :events, [Types::EventType], null: false
    field :occurrences, [Types::OccurrenceType], null: false
  end
end
