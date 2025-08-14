class CreateTagProjectJoinTable < ActiveRecord::Migration[7.0]
  def change
    create_join_table :tags, :projects do |t|
      # t.index [:tag_id, :project_id]
      # t.index [:project_id, :tag_id]
    end
  end
end
