require "rails_helper"

RSpec.describe Components::NavLink do
  it "renders a plain nav item by default" do
    html = described_class.new(label: "Home", href: "/").call

    expect(html).to include('<a href="/" class="nav-link">Home</a>')
    expect(html).not_to include("aria-")
  end

  it "marks the active link as the current page" do
    html = described_class.new(label: "Home", href: "/", active: true).call

    expect(html).to include("nav-link active").and include('aria-current="page"')
  end

  it "marks a disabled link" do
    html = described_class.new(label: "Link", disabled: true).call

    expect(html).to include("nav-link disabled").and include('aria-disabled="true"')
  end

  it "defaults href to an anchor" do
    expect(described_class.new(label: "Link").call).to include('href="#"')
  end
end
