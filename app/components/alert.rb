# frozen_string_literal: true

# https://getbootstrap.com/docs/5.3/components/alerts/
class Components::Alert < Components::Base
  def initialize(color: :primary)
    @color = color
  end

  def view_template(&block)
    div(class: "alert alert-#{@color}", role: "alert", &block)
  end
end
