require "test_helper"

module Import
  module Moneydance
    class UtilTest < ActiveSupport::TestCase
      class ParseStringDate < ActiveSupport::TestCase
        def expect_err(source, error_type, matcher)
          e = assert_raise(error_type) { Util.parse_string_date(source) }
          assert_match matcher, e.message
        end

        test "raise when input is empty" do
          expect_err "", ArgumentError, /expecting a 8 character string containing only numbers/
        end

        test "raise when input is not long enough" do
          expect_err "202610", ArgumentError, /expecting a 8 character string containing only numbers/
        end

        test "raise when input is too long" do
          expect_err "2026103101", ArgumentError, /expecting a 8 character string containing only numbers/
        end

        test "raise when input is not numbers" do
          expect_err "2026ABCD", ArgumentError, /expecting a 8 character string containing only numbers/
        end

        test "raise when date is not valid" do
          expect_err "20261401", Date::Error, /invalid date/
        end

        test "parse date properly" do
          assert_equal Date.new(2026, 10, 31), Util.parse_string_date("20261031")
        end
      end
    end
  end
end
