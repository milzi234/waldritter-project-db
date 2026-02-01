class ExtractionsController < ApplicationController
  # POST /api/v1/extract
  # Proxies extraction request to Python service
  def extract
    url = params[:url]

    if url.blank?
      render json: { error: "URL is required" }, status: :bad_request
      return
    end

    # Fetch all tags with their categories for matching
    available_tags = Tag.includes(:category).map do |tag|
      {
        id: tag.id,
        name: tag.title,
        category_name: tag.category&.title || ""
      }
    end

    # Build request to Python service
    extractor_url = ENV.fetch("EXTRACTOR_SERVICE_URL", "http://localhost:8000")

    begin
      response = HTTParty.post(
        "#{extractor_url}/api/extract",
        body: {
          url: url,
          available_tags: available_tags
        }.to_json,
        headers: { "Content-Type" => "application/json" },
        timeout: 120 # Allow time for exploration and LLM calls
      )

      if response.success?
        render json: response.parsed_response
      else
        error_message = response.parsed_response&.dig("detail") || "Extraction failed"
        render json: { error: error_message }, status: response.code
      end
    rescue HTTParty::Error, Errno::ECONNREFUSED => e
      render json: { error: "Extractor service unavailable: #{e.message}" }, status: :service_unavailable
    rescue Timeout::Error
      render json: { error: "Extraction timed out" }, status: :gateway_timeout
    end
  end
end
