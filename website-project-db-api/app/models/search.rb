class Search
  def self.search(search)
    if search[:tags].present?
      
      tags = Tag.where(id: search[:tags]).includes(:category)
      sorted_tags = tags.each_with_object({}) do |tag, hash|
        category_id = tag.category.id
        hash[category_id] ||= []
        hash[category_id] << tag.id
      end
      projects = []
      while !sorted_tags.empty?
        category_id_with_least_tags = sorted_tags.min_by { |_, v| v.size }.first
        if projects.empty?
          projects = Project.joins(:tags, events: :occurrences).where(tags: {id: sorted_tags[category_id_with_least_tags]}).where('occurrences.start_date >= ? AND occurrences.end_date <= ?', search[:start_date], search[:end_date]).distinct
        else
          projects = Project.joins(:tags).where(projects: {id: projects.map(&:id)}, tags: {id: sorted_tags[category_id_with_least_tags]}).distinct
        end
        if projects.empty?
          return {projects: {}, events: {}, occurrences: []}
        end
        sorted_tags.delete(category_id_with_least_tags)
      end
      if projects.empty?
        return {projects: {}, events: {}, occurrences: []}
      end

      return _massage({
        projects: projects, 
        events: projects.map(&:events).flatten.sort_by(&:start_date), 
        occurrences: projects.map(&:events).flatten.map(&:occurrences).flatten.sort_by(&:start_date)
      })
    else
      occurrences = Occurrence.includes(event: :project).where('start_date >= ? AND end_date <= ?', search[:start_date], search[:end_date]).order(:start_date)
      projects = occurrences.map { |occurrence| occurrence.event.project }.uniq
      events = occurrences.map(&:event).uniq
      return _massage({occurrences: occurrences, projects: projects, events: events})
    end
  end

  private

  def self._massage(data)
    massaged = {projects: {}, events: {}, occurrences: []}
    data[:projects].each do |project|
      massaged[:projects][project.id] = project
    end

    data[:events].each do |event|
      massaged[:events][event.id] = event
    end

    massaged[:occurrences] = data[:occurrences]

    return massaged
  end
end

