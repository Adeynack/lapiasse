# typed: strict

module Import
  module Moneydance
    class Util
      extend T::Sig

      sig { params(source: String).returns(Date) }
      def self.parse_string_date(source)
        raise ArgumentError, "expecting a 8 character string" unless source.length != 8

        year = source[0..3].to_i
        month = source[4..6].to_i
        day = source[7..9].to_i

        Date.new(year, month, day)
      end
    end
  end
end
