# frozen_string_literal: true

# https://getbootstrap.com/docs/5.3/components/navbar
class Components::NavLink < Components::Base
  # href:     Target of the link.
  # label:    Text to display for the link.
  # active:   If the link is showed as the current page.
  # disabled: Should the link be disabled?
  def initialize(href: "#", label: nil, active: false, disabled: false)
    @href = href
    @label = label
    @active = active
    @disabled = disabled
  end

  def view_template
    li(class: "nav-item") do
      a(
        href: @href,
        # nil entries are dropped from class arrays. `false` is not: it raises,
        # so these must be `(x if cond)` and never `cond && x`.
        class: ["nav-link", ("active" if @active), ("disabled" if @disabled)],
        aria: {
          current: ("page" if @active),
          disabled: ("true" if @disabled)
        },
      ) { @label }
    end
  end
end
