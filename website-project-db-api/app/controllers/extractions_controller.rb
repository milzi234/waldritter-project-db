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

  # POST /api/v1/extract_text
  # Proxies text extraction request to Python service
  def extract_text
    text = params[:text]

    if text.blank?
      render json: { error: "Text is required" }, status: :bad_request
      return
    end

    available_tags = Tag.includes(:category).map do |tag|
      {
        id: tag.id,
        name: tag.title,
        category_name: tag.category&.title || ""
      }
    end

    extractor_url = ENV.fetch("EXTRACTOR_SERVICE_URL", "http://localhost:8000")

    begin
      response = HTTParty.post(
        "#{extractor_url}/api/extract-text",
        body: {
          text: text,
          available_tags: available_tags
        }.to_json,
        headers: { "Content-Type" => "application/json" },
        timeout: 120
      )

      if response.success?
        render json: response.parsed_response
      else
        error_message = response.parsed_response&.dig("detail") || "Text extraction failed"
        render json: { error: error_message }, status: response.code
      end
    rescue HTTParty::Error, Errno::ECONNREFUSED => e
      render json: { error: "Extractor service unavailable: #{e.message}" }, status: :service_unavailable
    rescue Timeout::Error
      render json: { error: "Extraction timed out" }, status: :gateway_timeout
    end
  end

  # POST /api/v1/generate_image
  # Proxies image generation request to Python service
  def generate_image
    title = params[:title]
    description = params[:description]
    keywords = params[:keywords] || []
    variation = params[:variation] || 0

    if title.blank?
      render json: { error: "Title is required" }, status: :bad_request
      return
    end

    extractor_url = ENV.fetch("EXTRACTOR_SERVICE_URL", "http://localhost:8000")

    begin
      response = HTTParty.post(
        "#{extractor_url}/api/generate-image",
        body: {
          title: title,
          description: description || "",
          keywords: keywords,
          variation: variation.to_i
        }.to_json,
        headers: { "Content-Type" => "application/json" },
        timeout: 60 # Image generation may take time
      )

      if response.success?
        render json: response.parsed_response
      else
        error_message = response.parsed_response&.dig("detail") || "Image generation failed"
        render json: { error: error_message }, status: response.code
      end
    rescue HTTParty::Error, Errno::ECONNREFUSED => e
      render json: { error: "Extractor service unavailable: #{e.message}" }, status: :service_unavailable
    rescue Timeout::Error
      render json: { error: "Image generation timed out" }, status: :gateway_timeout
    end
  end
end
