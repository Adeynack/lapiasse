# frozen_string_literal: true

# A link styled as a button.
# https://getbootstrap.com/docs/5.3/components/buttons/
#
# Replaces the former ApplicationHelper#link_button_to. The label moves from a
# leading positional argument to the block, which is the Phlex idiom.
class Components::LinkButton < Components::Base
  def initialize(href, variant: :primary, disabled: false, **attributes)
    @href = href
    @variant = variant
    @disabled = disabled
    @attributes = attributes
  end

  # &block: The button's label.
  def view_template(&block)
    a(
      href: @href,
      class: ["btn", "btn-#{@variant}", ("disabled" if @disabled)],
      aria: {disabled: ("true" if @disabled)},
      tabindex: ("-1" if @disabled),
      **@attributes,
      &block
    )
  end
end
