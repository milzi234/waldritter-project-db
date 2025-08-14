class EventException < ApplicationRecord
  belongs_to :event
  self.table_name = 'exceptions'
end
