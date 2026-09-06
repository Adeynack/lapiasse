require "rails_helper"

RSpec.describe Components::Card do
  it "renders the block as the card body" do
    html = described_class.new.call { "content" }

    expect(html).to include('<div class="card-body">content</div>')
  end

  it "renders a header when a title is given" do
    html = described_class.new(title: "Current User").call

    expect(html).to include('<div class="card-header">Current User</div>')
  end

  it "omits the header when no title is given" do
    expect(described_class.new.call).not_to include("card-header")
  end

  it "omits the header when the title is blank" do
    expect(described_class.new(title: "  ").call).not_to include("card-header")
  end

  it "escapes the title" do
    html = described_class.new(title: "<script>").call

    expect(html).to include("&lt;script&gt;")
    expect(html).not_to include("<script>")
  end
end
