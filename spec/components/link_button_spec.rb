require "rails_helper"

RSpec.describe Components::LinkButton do
  it "renders a primary button link by default" do
    html = described_class.new("/out").call { "Log out" }

    expect(html).to include('href="/out"')
    expect(html).to include('class="btn btn-primary"')
    expect(html).to include("Log out")
  end

  it "applies the given variant" do
    expect(described_class.new("/x", variant: :danger).call).to include("btn-danger")
  end

  it "marks a disabled button" do
    html = described_class.new("/x", disabled: true).call

    expect(html).to include("btn btn-primary disabled").and include('aria-disabled="true"')
  end

  it "passes extra attributes through to the anchor" do
    html = described_class.new("/out", data: {turbo_method: :delete}).call

    expect(html).to include('data-turbo-method="delete"')
  end
end
