module ToolbarHelper
  # Renders a `toolbar`, yielding a `ToolbarBuilder` to add its items.
  #
  #   <%= toolbar do |t| %>
  #     <%= t.to "Home", root_path %>
  #     <%= t.delete "Log out", destroy_user_session_path %>
  #   <% end %>
  def toolbar(**options, &)
    builder = ToolbarBuilder.new(self)
    tag.div(**options, class: class_names("toolbar", options[:class])) do
      capture(builder, &)
    end
  end
end
