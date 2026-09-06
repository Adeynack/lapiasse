# frozen_string_literal: true

# The application chrome. Rendered from app/views/layouts/application.html.erb,
# which stays ERB so that Rails' layout pipeline keeps serving Devise's
# gem-provided views alongside this app's Phlex views.
class Components::Layout < Components::Base
  include Phlex::Rails::Helpers::ContentFor
  include Phlex::Rails::Helpers::CSPMetaTag
  include Phlex::Rails::Helpers::CSRFMetaTags
  include Phlex::Rails::Helpers::Flash
  include Phlex::Rails::Helpers::JavaScriptImportmapTags
  include Phlex::Rails::Helpers::StyleSheetLinkTag

  def initialize(current_user: nil)
    @current_user = current_user
  end

  def view_template
    doctype
    html(data: {bs_theme: "dark"}) do
      head do
        title { content_for(:title) || "Lapiasse" }
        meta(name: "viewport", content: "width=device-width,initial-scale=1")
        meta(name: "apple-mobile-web-app-capable", content: "yes")
        meta(name: "application-name", content: "La Piasse")
        meta(name: "mobile-web-app-capable", content: "yes")
        csrf_meta_tags
        csp_meta_tag

        # `raw` accepts a SafeBuffer and no-ops on nil, so this covers the
        # former `yield :head` whether or not a view filled it in.
        raw(content_for(:head))

        # Enable PWA manifest for installable apps (also enable it in config/routes.rb!)
        # link(rel: "manifest", href: pwa_manifest_path(format: :json))

        link(rel: "icon", href: "/icon.png", type: "image/png")
        link(rel: "icon", href: "/icon.svg", type: "image/svg+xml")
        link(rel: "apple-touch-icon", href: "/icon.png")

        # Includes all stylesheet files in app/assets/stylesheets
        stylesheet_link_tag(:app, "data-turbo-track": "reload")
        javascript_importmap_tags
      end

      body do
        AppNavbar(current_user: @current_user)

        div(class: "container-fluid mt-2") do
          Alert(color: :primary) { flash[:notice] } if flash[:notice].present?
          Alert(color: :warning) { flash[:alert] } if flash[:alert].present?
          yield
        end
      end
    end
  end
end
