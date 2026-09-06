# frozen_string_literal: true

# La Piasse's navbar: the application-specific instance of Components::Navbar.
class Components::AppNavbar < Components::Base
  def initialize(current_user: nil)
    @current_user = current_user
  end

  def view_template
    Navbar(brand: "La Piasse") do
      if @current_user
        NavLink(label: "Home", href: "/", active: true)
        NavLink(label: "Link", href: "#", disabled: true)
      end
    end
  end
end
