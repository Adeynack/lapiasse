require "test_helper"
require_relative "../test_reporter"

module Import
  module Moneydance
    describe Importer do
      let(:book) { Importer.new(reporter:, json_content:, book_owner_email:, default_currency_iso_code:).import! }
      let(:reporter) { TestReporter.new }
      let(:json_content) { file_fixture("import/moneydance/with_metadata.json").read }
      let(:owner) { users(:joe) }
      let(:book_owner_email) { owner.email }
      let(:default_currency_iso_code) { "eur" }

      describe "without metadata" do
        let(:json_content) { "{}" }
        it "raises missing metadata" do
          e = assert_raise(Importer::MissingInformationError) { book }
          assert_equal "Missing metadata in MoneyDance Export", e.message
        end
      end

      describe "without a filename" do
        let(:json_content) do
          <<~JSON
            {
              "metadata": {
                "export_date": "20010101"
              }
            }
          JSON
        end

        it "raises missing export date" do
          e = assert_raise(Importer::MissingInformationError) { book }
          assert_equal "Missing filename (in metadata) in MoneyDance Export", e.message
        end
      end

      describe "without an export date" do
        let(:json_content) do
          <<~JSON
            {
              "metadata": {
                "file_name": "foo"
              }
            }
          JSON
        end

        it "raises missing export date" do
          e = assert_raise(Importer::MissingInformationError) { book }
          assert_equal "Missing export_date (in metadata) in MoneyDance Export", e.message
        end
      end
      describe "no use exists with the email" do
        let(:book_owner_email) { "does_not_exist@example.com" }

        it "raises user cannot be found" do
          e = assert_raise(Importer::DataError) { book }
          assert_equal "No user with E-Mail does_not_exist@example.com found", e.message
        end
      end

      describe "a book already exists" do
        it "deletes the existing book and creates a new one"
        it "raises that a book alreadu exists"
      end
    end
  end
end
