class TokenValidator

  TOKEN_CACHE = LruRedux::Cache.new(500, 1.hour.to_i)

  def self.validate(token)
    raise ArgumentError, "Token cannot be empty or nil" if token.nil? || token.empty?

    TOKEN_CACHE.getset(token) do
      puts "Fetching token #{token} from DEX"
      response = HTTParty.post(Rails.configuration.dex_introspect_url,
        headers: { 'Content-Type' => 'application/x-www-form-urlencoded' },
        body: {
          token: token,
          client_id: Rails.configuration.dex_client_id,
          client_secret: Rails.configuration.dex_client_secret
        }
      )
      puts "Response from DEX: #{response.body}"
      JSON.parse(response.body)
    end
  end
end