module ApplicationHelper
  def link_button_to(
    target = nil,
    options = nil,
    disabled: false,
    **kwargs
  )
    link_to target, options,
      class: "button",
      aria: {
        disabled: ("true" if disabled)
      },
      tabindex: ("-1" if disabled),
      **kwargs
  end
end
