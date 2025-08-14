class Project < ApplicationRecord
  has_and_belongs_to_many :tags
  has_many :events, dependent: :destroy
  has_many :occurrences, through: :events
  has_one_attached :image
  belongs_to :umbrella_project, optional: true

end
