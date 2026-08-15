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
RSpec.describe Book, type: :model do
  pending "add some examples to (or delete) #{__FILE__}"
end
