# frozen_string_literal: true

module Types
  class EventType < Types::BaseObject
    field :id, ID, null: false
    field :start_date, GraphQL::Types::ISO8601DateTime
    field :end_date, GraphQL::Types::ISO8601DateTime
    field :project_id, Integer, null: false
    field :created_at, GraphQL::Types::ISO8601DateTime, null: false
    field :updated_at, GraphQL::Types::ISO8601DateTime, null: false
    field :recurrence_type, String
  end
end
