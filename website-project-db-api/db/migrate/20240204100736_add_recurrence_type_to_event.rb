class AddRecurrenceTypeToEvent < ActiveRecord::Migration[7.0]
  def change
    add_column :events, :recurrence_type, :string, default: 'none'
    remove_column :events, :weekly
  end
end
