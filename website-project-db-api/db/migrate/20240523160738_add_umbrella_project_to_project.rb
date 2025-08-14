class AddUmbrellaProjectToProject < ActiveRecord::Migration[7.0]
  def change
    add_reference :projects, :umbrella_project, null: true, foreign_key: true
  end
end
