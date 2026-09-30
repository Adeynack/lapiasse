# Builds the items of a `toolbar`. See `ToolbarHelper#toolbar`.
class ToolbarBuilder
  ITEM_CLASS = "button toolbar__item"

  def initialize(template)
    @template = template
  end

  # A plain link (GET).
  def to(name, url, **options)
    @template.link_to name, url, **options, class: item_class(options[:class])
  end

  # A button in its own form, submitted with the given HTTP verb.
  %i[post put patch delete].each do |verb|
    define_method(verb) do |name, url, **options|
      @template.button_to name, url, **options, method: verb, class: item_class(options[:class])
    end
  end

  private

  def item_class(extra)
    @template.class_names(ITEM_CLASS, extra)
  end
end
