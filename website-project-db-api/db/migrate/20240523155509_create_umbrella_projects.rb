class CreateUmbrellaProjects < ActiveRecord::Migration[7.0]
  def change
    create_table :umbrella_projects do |t|
      t.string :name
      t.string :title
      t.text :description
      t.datetime :start_date
      t.datetime :end_date

      t.timestamps
    end
  end
end
