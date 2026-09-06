# frozen_string_literal: true

class Components::Base < Phlex::HTML
  # Rails helpers are opt-in per component in Phlex. These are the ones common
  # enough to be worth having everywhere; anything more specialised (asset tags,
  # CSRF, importmap) is included by the single component that needs it.
  include Phlex::Rails::Helpers::Routes
  include Phlex::Rails::Helpers::LinkTo
  include Phlex::Rails::Helpers::MailTo

  if Rails.env.development?
    def before_template
      comment { "Before #{self.class.name}" }
      super
    end

    def after_template
      super
      comment { "After #{self.class.name}" }
    end
  end
end
