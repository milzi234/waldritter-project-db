class CreateOccurrences < ActiveRecord::Migration[7.0]
  def change
    create_table :occurrences do |t|
      t.datetime :start_date
      t.datetime :end_date
      t.references :event, null: false, foreign_key: true
      t.timestamps
    end
  end
end
