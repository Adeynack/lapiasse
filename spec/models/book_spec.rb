require 'rails_helper'

# == Schema Information
#
# Table name: books
#
#  id                  :integer          not null, primary key
#  name                :string           not null
#  created_at          :datetime         not null
#  updated_at          :datetime         not null
#  default_currency_id :integer          not null
#  owner_id            :integer          not null
#
# Indexes
#
#  index_books_on_default_currency_id  (default_currency_id)
#  index_books_on_owner_id             (owner_id)
#
# Foreign Keys
#
#  default_currency_id  (default_currency_id => currencies.id)
#  owner_id             (owner_id => users.id)
#
RSpec.describe Book do
  fixtures :all

  def book_with_defaults(**attr)
    Book.new(
      name: "Test Book",
      default_currency: currencies(:eur),
      owner: users(:joe),
      **attr,
    )
  end

  describe "is valid" do
    it "with all needed information" do
      expect(book_with_defaults).to be_valid
    end
  end
end
