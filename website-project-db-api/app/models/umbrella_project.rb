class UmbrellaProject < ApplicationRecord
  has_many :projects
  has_one_attached :image
end
