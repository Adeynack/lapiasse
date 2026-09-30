require "test_helper"

class ToolbarHelperTest < ActionView::TestCase
  test "renders a plain link with `to`" do
    render inline: <<~ERB
      <%= toolbar do |t| %>
        <%= t.to "Home", "/home" %>
      <% end %>
    ERB

    assert_dom "div.toolbar > a.button.toolbar__item[href='/home']", text: "Home"
  end

  test "renders a button in a form with the verb" do
    render inline: <<~ERB
      <%= toolbar do |t| %>
        <%= t.delete "Log out", "/logout", class: "extra" %>
      <% end %>
    ERB

    assert_dom "div.toolbar > form[action='/logout']" do
      assert_dom "input[name='_method'][value='delete']"
      assert_dom "button.button.toolbar__item.extra", text: "Log out"
    end
  end
end
