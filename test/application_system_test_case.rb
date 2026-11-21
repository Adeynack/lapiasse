require "test_helper"

# Rails resolves the chromedriver path eagerly to avoid races in parallel tests,
# which also stops Selenium Manager from resolving the browser itself. Do that
# part here: it picks up an installed Chrome, and downloads a matching Chrome
# for Testing build on machines that have none.
Selenium::WebDriver::Chrome.path ||=
  Selenium::WebDriver::SeleniumManager.binary_paths("--browser", "chrome")["browser_path"]

# Chrome 158 rejects `Storage.clearDataForOrigin` with origin "*", which
# Capybara sends from `ChromeDriver#reset!` after every system test. Clear the
# cookies globally and the remaining storage for the app server's origin instead.
# Remove once a Capybara release includes one of the fixes linked from:
# https://github.com/teamcapybara/capybara/issues/2865
if Gem.loaded_specs["capybara"].version <= Gem::Version.new("3.40.0")
  module ClearStorageWithoutWildcardOrigin
    private

    def execute_cdp(cmd, params = {})
      return super unless cmd == "Storage.clearDataForOrigin" && params[:origin] == "*"

      types = params[:storageTypes].to_s.split(",")
      super("Network.clearBrowserCookies") if types.delete("cookies")

      origin = Capybara.current_session.server&.base_url
      super(cmd, params.merge(origin: origin, storageTypes: types.join(","))) if origin && types.any?
    end
  end
  Capybara::Selenium::Driver::ChromeDriver.prepend(ClearStorageWithoutWildcardOrigin)
end

class ApplicationSystemTestCase < ActionDispatch::SystemTestCase
  driven_by :selenium, using: :headless_chrome, screen_size: [1400, 800]

  def sign_in_as(user, password)
    visit new_user_session_path
    fill_in "Email", with: user.email
    fill_in "Password", with: password
    click_on "Log in"
  end
end
