# frozen_string_literal: true

# https://getbootstrap.com/docs/5.3/components/navbar/
class Components::Navbar < Components::Base
  MENU_ID = "navbarSupportedContent"

  def initialize(brand: nil)
    @brand = brand
  end

  # &block: The nav items, normally Components::NavLink.
  def view_template(&block)
    nav(class: "navbar navbar-expand-sm bg-body-tertiary") do
      div(class: "container-fluid") do
        a(class: "navbar-brand", href: "#") { @brand } if @brand.present?
        toggler
        div(class: "collapse navbar-collapse", id: MENU_ID) do
          ul(class: "navbar-nav me-auto mb-2 mb-sm-0", &block)
        end
      end
    end
  end

  private

  def toggler
    button(
      class: "navbar-toggler",
      type: "button",
      data: {bs_toggle: "collapse", bs_target: "##{MENU_ID}"},
      aria: {
        controls: MENU_ID,
        # Strings, not booleans: a `false` attribute value is dropped entirely.
        expanded: "false",
        label: "Toggle navigation"
      },
    ) do
      span(class: "navbar-toggler-icon")
    end
  end
end
