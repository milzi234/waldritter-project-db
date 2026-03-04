class AddHomepageToProjects < ActiveRecord::Migration[7.0]
  def change
    add_column :projects, :homepage, :string
  end
end
