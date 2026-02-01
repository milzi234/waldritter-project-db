module Authenticable
  extend ActiveSupport::Concern

  class_methods do
    def skip_authentication(*actions)
      skip_authentication_actions.concat(actions)
    end

    def skip_authentication_actions
      @skip_authentication_actions ||= []
    end
  end

  included do
    before_action :authenticate_request
  end

  private

  def authenticate_request
    return if skip_authentication?

    header = request.headers['Authorization']
    token = header.split(' ').last if header
    decoded_token = TokenValidator.validate(token)
    
    if decoded_token['active']
      # Find or create user by sub (subject identifier)
      @current_user = User.find_or_create_by(sub: decoded_token['sub']) do |u|
        u.email = decoded_token['email']
        u.name = decoded_token['name']
      end
      
      # Update user info if it has changed
      if @current_user.email != decoded_token['email'] || @current_user.name != decoded_token['name']
        @current_user.update(
          email: decoded_token['email'],
          name: decoded_token['name']
        )
      end
    else
      render json: { error: 'Invalid token' }, status: :unauthorized
    end
  end

  def skip_authentication?
    Rails.configuration.skip_authentication ||
      self.class.skip_authentication_actions.include?(action_name.to_sym)
  end
end