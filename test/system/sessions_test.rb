require "application_system_test_case"

class SessionsTest < ApplicationSystemTestCase
  test "signing in shows the current user on the home page" do
    sign_in_as users(:joe), "joe"

    assert_text "Currently logged as Joe"
  end

  test "signing in with a wrong password keeps the visitor on the log in page" do
    visit new_user_session_path
    fill_in "Email", with: users(:joe).email
    fill_in "Password", with: "not-joe"
    click_on "Log in"

    assert_text "Invalid email or password."
  end

  test "signing out brings the user back to the log in page" do
    sign_in_as users(:joe), "joe"
    click_link "Log out"

    assert_text "Log in"
    assert_no_text "Currently logged as"
  end
end
