require_relative '../config/environment'

require 'yaml'

# Clear the tables
Category.destroy_all
Tag.destroy_all
Project.destroy_all

# Load the YAML file
categories = YAML.load_file('import/categories.yaml')

# Import the data into the database
categories['categories'].each do |category_data|
  category = Category.create(title: category_data['title'])
  puts "Importing #{category.title}"
  category_data['tags'].each do |tag_data|
    category.tags.create(title: tag_data['title'])
    puts "  - Importing #{tag_data['title']}"
  end
end
