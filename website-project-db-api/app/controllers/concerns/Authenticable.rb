# app/controllers/concerns/authenticatable.rb
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
      @current_user = User.find_or_create_by(sub: decoded_token['sub'])
    else
      render json: { error: 'Invalid token' }, status: :unauthorized
    end
  end

  def skip_authentication?
    self.class.skip_authentication_actions.include?(action_name.to_sym)
  end

end