# frozen_string_literal: true

module Types
  class ProjectType < Types::BaseObject
    field :id, ID, null: false
    field :title, String
    field :description, String
    field :created_at, GraphQL::Types::ISO8601DateTime, null: false
    field :updated_at, GraphQL::Types::ISO8601DateTime, null: false
    field :umbrella_project_id, Integer

    field :image_url, String, null: true
    def image_url
      object.image.attached? ? Rails.application.routes.url_helpers.url_for(object.image) : nil
    end

    field :tags, [Types::TagType], null: false
    def tags
      object.tags
    end
  end
end
