require "test_helper"

# Rails resolves the chromedriver path eagerly to avoid races in parallel tests,
# which also stops Selenium Manager from resolving the browser itself. Do that
# part here: it picks up an installed Chrome, and downloads a matching Chrome
# for Testing build on machines that have none.
Selenium::WebDriver::Chrome.path ||=
  Selenium::WebDriver::SeleniumManager.binary_paths("--browser", "chrome")["browser_path"]

class ApplicationSystemTestCase < ActionDispatch::SystemTestCase
  driven_by :selenium, using: :headless_chrome, screen_size: [1400, 800]

  def sign_in_as(user, password)
    visit new_user_session_path
    fill_in "Email", with: user.email
    fill_in "Password", with: password
    click_on "Log in"
  end
end
