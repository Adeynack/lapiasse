module ApplicationHelper
  def link_button_to(
    target = nil,
    options = nil,
    disabled: false,
    **kwargs
  )
    link_to target, options,
      class: class_names("button", kwargs.delete(:class)),
      aria: {
        disabled: ("true" if disabled)
      },
      tabindex: ("-1" if disabled),
      **kwargs
  end
end
