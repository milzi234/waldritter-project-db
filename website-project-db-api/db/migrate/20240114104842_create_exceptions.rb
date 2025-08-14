class CreateExceptions < ActiveRecord::Migration[7.0]
  def change
    create_table :exceptions do |t|
      t.datetime :start_date
      t.datetime :end_date
      t.references :event, null: false, foreign_key: true
      t.timestamps
    end
  end
end
