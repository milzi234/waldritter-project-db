class CreateEvents < ActiveRecord::Migration[7.0]
  def change
    create_table :events do |t|
      t.datetime :start_date
      t.datetime :end_date
      t.boolean :weekly, default: false
      t.references :project, null: false, foreign_key: true
      t.timestamps
    end
  end
end
