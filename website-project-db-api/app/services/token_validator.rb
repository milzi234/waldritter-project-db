require 'googleauth'

class TokenValidator
  TOKEN_CACHE = LruRedux::Cache.new(500, 1.hour.to_i)
  ALLOWED_DOMAIN = 'waldritter.de'
  
  def self.validate(token)
    Rails.logger.info "TokenValidator: Validating token: #{token&.first(20)}..."
    raise ArgumentError, "Token cannot be empty or nil" if token.nil? || token.empty?
    
    TOKEN_CACHE.getset(token) do
      validate_google_token(token)
    end
  end
  
  private
  
  def self.validate_google_token(token)
    client_id = Rails.configuration.google_client_id
    Rails.logger.info "Google OAuth: Validating with client_id: #{client_id}"
    
    begin
      # Use the correct Google Auth method
      payload = Google::Auth::IDTokens.verify_oidc(token, aud: client_id)
      
      Rails.logger.info "Google OAuth: Token payload: #{payload.inspect}"
      
      # Check domain restriction
      if payload['hd'] != ALLOWED_DOMAIN
        Rails.logger.info "Google OAuth: Invalid domain #{payload['hd']}"
        return { 'active' => false, 'error' => 'Invalid domain. Only waldritter.de emails are allowed.' }
      end
      
      Rails.logger.info "Google OAuth: Successfully validated token for #{payload['email']}"
      
      {
        'active' => true,
        'sub' => payload['sub'],
        'email' => payload['email'],
        'name' => payload['name'],
        'picture' => payload['picture']
      }
    rescue => e
      Rails.logger.error "Google OAuth validation failed: #{e.message}"
      Rails.logger.error "Google OAuth validation backtrace: #{e.backtrace.first(5).join("\n")}"
      { 'active' => false, 'error' => e.message }
    end
  end
end