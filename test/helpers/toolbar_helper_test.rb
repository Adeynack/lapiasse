require "test_helper"

class ToolbarHelperTest < ActionView::TestCase
  test "renders a plain link with `to`" do
    render inline: <<~ERB
      <%= toolbar do |t| %>
        <%= t.to "Home", "/home" %>
      <% end %>
    ERB

    assert_dom "div > a[role='button'][href='/home']", text: "Home"
  end

  test "renders a submit button targeting the URL with the verb" do
    render inline: <<~ERB
      <%= toolbar do |t| %>
        <%= t.delete "Log out", "/logout", class: "extra" %>
      <% end %>
    ERB

    button = "a.extra[href='/logout'][role='button'][data-turbo-method='delete']"

    assert_dom "div > #{button}", text: "Log out"
  end

  test "keeps every item a direct child of the group" do
    render inline: <<~ERB
      <%= toolbar do |t| %>
        <%= t.to "Home", "/home" %>
        <%= t.post "Save", "/save" %>
        <%= t.delete "Log out", "/logout" %>
      <% end %>
    ERB

    assert_dom "div > *", count: 3
    assert_dom "div form", count: 0
  end
end
