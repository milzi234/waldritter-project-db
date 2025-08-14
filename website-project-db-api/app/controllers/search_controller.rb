class SearchController < ApplicationController
  skip_authentication :search
  
  def search
    @search = search_params
    @results = Search.search(@search)
    
    render json: json_with_image_urls(@results)
  end

  private
    
    def json_with_image_urls(results)
      {
        projects: results[:projects].map do |id, project| 
          tag_ids = project.tags.map(&:id)
          if project.image.attached?
            [id, project.as_json.merge({image_url: url_for(project.image), tag_ids: tag_ids})]
          else
            [id, project.as_json.merge({tag_ids: tag_ids})]
          end
        end.to_h,
        events: results[:events].map { |id, event| [id, event.as_json] }.to_h,
        occurrences: results[:occurrences].map(&:as_json)
      }
    end

    # Only allow a list of trusted parameters through.
    def search_params
      params.require(:search).permit(:start_date, :end_date, tags: [])
    end
end

