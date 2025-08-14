class Event < ApplicationRecord
  belongs_to :project
  has_many :exceptions, dependent: :destroy, class_name: 'EventException'
  has_many :occurrences, dependent: :destroy

  def recalculate_occurrences
    self.occurrences.destroy_all
    if self.recurrence_type != 'no-repeat' &&  self.recurrence_type != 'none' #fix me
      exceptions = self.exceptions.pluck(:start_date)
      if self.recurrence_type == 'weekly'
        156.times do |week|
          start_date = self.start_date.advance({weeks: week })
          end_date = self.end_date.advance({weeks: week })
          skip = exceptions.any? { |exception| exception.day == start_date.day && exception.month == start_date.month && exception.year == start_date.year }
          self.occurrences.create(start_date: start_date, end_date: end_date) unless skip
        end
      elsif self.recurrence_type == 'monthly-day'
        start_day_of_week = self.start_date.wday
        start_week_of_month = (self.start_date.day - 1) / 7 + 1
        36.times do |month|
          start_date = self.start_date.advance({months: month })
          end_date = self.end_date.advance({months: month })
          start_date = start_date.advance({days: (start_week_of_month - 1) * 7 + start_day_of_week})
          end_date = end_date.advance({days: (start_week_of_month - 1) * 7 + start_day_of_week})
          skip = exceptions.any? { |exception| exception.day == start_date.day && exception.month == start_date.month && exception.year == start_date.year }
          self.occurrences.create(start_date: start_date, end_date: end_date) unless skip
        end
      elsif self.recurrence_type == 'monthly-date'
        36.times do |month|
          start_date = self.start_date.advance({months: month })
          end_date = self.end_date.advance({months: month })
          skip = exceptions.any? { |exception| exception.day == start_date.day && exception.month == start_date.month && exception.year == start_date.year }
          self.occurrences.create(start_date: start_date, end_date: end_date) unless skip
        end
      end
    else
      self.occurrences.create(start_date: self.start_date, end_date: self.end_date)
    end
  end
end
