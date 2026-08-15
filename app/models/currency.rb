# == Schema Information
#
# Table name: currencies
#
#  id              :integer          not null, primary key
#  iso_code        :string(3)
#  name            :string           not null
#  subunit_to_unit :integer          default(100), not null
#  symbol          :string
#  symbol_first    :boolean          default(FALSE), not null
#  created_at      :datetime         not null
#  updated_at      :datetime         not null
#
# Indexes
#
#  index_currencies_on_iso_code  (iso_code) UNIQUE
#
class Currency < ApplicationRecord
  validates :iso_code, uniqueness: true
  validates :name, presence: true

  def self.seed!
    # Some ISO Codes in the Money gem are repeated. Creating the Currency object with the unique
    # codes will ensure the proper version is used.
    Money::Currency.map(&:iso_code).uniq!.sort!.map! { Money::Currency.new(it) }.each do |c|
      Currency.find_or_create_by!(iso_code: c.iso_code) do
        it.assign_attributes(
          name: c.name,
          symbol: c.symbol,
          symbol_first: c.symbol_first,
          subunit_to_unit: c.subunit_to_unit,
        )
      end
    end

    nil
  end
end
