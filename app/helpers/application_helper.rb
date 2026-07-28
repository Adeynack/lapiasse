module ApplicationHelper
  def link_button_to(
    target = nil,
    options = nil,
    variant: :primary,
    disabled: false,
    **kwargs
  )
    link_to target, options,
      class: class_names("btn", "btn-#{variant}", disabled:),
      aria: {
        disabled: ("true" if disabled)
      },
      tabindex: ("-1" unless disabled),
      **kwargs
  end
end
