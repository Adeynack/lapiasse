require "rails_helper"

# End-to-end cover for the Phlex layout. The ERB shim in
# app/views/layouts/application.html.erb has to serve two different kinds of
# view, so both are exercised here.
RSpec.describe "Application layout", type: :request do
  fixtures :users

  shared_examples "the application chrome" do
    it "renders the document shell" do
      # Phlex emits a lowercase doctype; HTML5 treats it case-insensitively.
      expect(response.body).to start_with("<!doctype html>")
      expect(response.body).to include('<html data-bs-theme="dark">')
      expect(response.body).to include("<title>Lapiasse</title>")
    end

    it "renders the navbar with the brand" do
      expect(response.body).to include('<a class="navbar-brand" href="#">La Piasse</a>')
    end

    # csrf_meta_tags and csp_meta_tag are deliberately not asserted here: both
    # return nil in the test environment (forgery protection is off, and no CSP
    # is configured), so they can only be checked against a booted app.
    it "renders the asset tags" do
      expect(response.body).to include('<link rel="stylesheet" href="/assets/application-')
      expect(response.body).to include('data-turbo-track="reload"')
      expect(response.body).to include('<script type="importmap"')
    end
  end

  context "with a Devise view, which is still ERB" do
    before { get new_user_session_path }

    include_examples "the application chrome"

    it "renders no nav links when signed out" do
      expect(response.body).not_to include("nav-link")
    end
  end

  context "with a Phlex view" do
    before do
      sign_in users(:joe)
      get root_path
    end

    include_examples "the application chrome"

    it "renders the nav links when signed in" do
      expect(response.body).to include("nav-link active").and include("nav-link disabled")
    end

    it "renders the page content inside the container" do
      expect(response.body).to include("Current User")
      expect(response.body).to include("Joe")
      expect(response.body).to include("joe@example.com")
      expect(response.body).to include("Log out")
    end
  end

  context "when there is a flash" do
    it "renders a warning alert for flash[:alert]" do
      post user_session_path, params: {user: {email: "joe@example.com", password: "wrong"}}

      expect(response.body).to include('class="alert alert-warning" role="alert"')
    end

    # Signing in redirects straight to root and renders there, so the notice
    # survives. Signing out bounces root -> sign-in, and the flash is swept
    # before the second hop renders anything.
    it "renders a primary alert for flash[:notice]" do
      post user_session_path, params: {user: {email: "joe@example.com", password: "joe"}}
      follow_redirect!

      expect(response.body).to include('class="alert alert-primary" role="alert"')
    end
  end
end
