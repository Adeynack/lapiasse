module Import
  module Moneydance
    class Util
      STRING_DATE_MATCHER = /^(\d{4})(\d{2})(\d{2})$/

      def self.parse_string_date(source)
        r = STRING_DATE_MATCHER.match(source)
        raise ArgumentError, "expecting a 8 character string containing only numbers" if r.nil?

        Date.new r[1].to_i, r[2].to_i, r[3].to_i
      end
    end
  end
end
