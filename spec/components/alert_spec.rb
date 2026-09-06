require "rails_helper"

RSpec.describe Components::Alert do
  it "defaults to the primary colour" do
    expect(described_class.new.call).to include('class="alert alert-primary"')
  end

  it "applies the given colour" do
    expect(described_class.new(color: :warning).call).to include("alert-warning")
  end

  it "renders the block as its content and keeps the alert role" do
    html = described_class.new.call { "Signed out." }

    expect(html).to include('role="alert"').and include("Signed out.")
  end
end
