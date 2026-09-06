# frozen_string_literal: true

# https://getbootstrap.com/docs/5.3/components/card/
class Components::Card < Components::Base
  # title: Optional title, displayed in the card's header.
  def initialize(title: nil)
    @title = title
  end

  # &block: Content of the card's body.
  def view_template(&block)
    div(class: "card") do
      div(class: "card-header") { @title } if @title.present?
      div(class: "card-body", &block)
    end
  end
end
