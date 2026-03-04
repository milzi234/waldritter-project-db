# frozen_string_literal: true

module Types
  class TagType < Types::BaseObject
    field :id, ID, null: false
    field :title, String
    field :description, String
    field :category_id, Integer, null: false
    field :created_at, GraphQL::Types::ISO8601DateTime, null: false
    field :updated_at, GraphQL::Types::ISO8601DateTime, null: false

    field :category, Types::CategoryType, null: true
  end
end
