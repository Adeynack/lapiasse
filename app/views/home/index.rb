# frozen_string_literal: true

class Views::Home::Index < Views::Base
  def initialize(current_user:)
    @current_user = current_user
  end

  def view_template
    Card(title: "Current User") do
      p do
        plain "Currently logged as #{@current_user.display_name} ("
        mail_to @current_user.email
        plain ")."
      end

      LinkButton(destroy_user_session_path, data: {turbo_method: :delete}) { "Log out" }
    end
  end
end
