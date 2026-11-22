module ToolbarHelper
  # Renders a `toolbar`, yielding a `ToolbarBuilder` to add its items.
  #
  #   <%= toolbar do |t| %>
  #     <%= t.to "Home", root_path %>
  #     <%= t.delete "Log out", destroy_user_session_path %>
  #   <% end %>
  def toolbar(**options, &)
    builder = ToolbarBuilder.new(self)
    tag.div(**options) do
      capture(builder, &)
    end
  end

  class ToolbarBuilder
    def initialize(template)
      @template = template
    end

    # A plain link (GET).
    def to(name, url, **options)
      @template.link_to name, url, **options, role: :button
    end

    # A button in its own form, submitted with the given HTTP verb.
    %i[post put patch delete].each do |verb|
      define_method(verb) do |name, url, **options|
        @template.link_to name, url, **options, role: :button, data: {turbo_method: verb}
      end
    end
  end
end
