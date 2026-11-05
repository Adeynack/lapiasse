require "test_helper"
require_relative "../test_reporter"

module Import
  module Moneydance
    describe Importer do
      let(:imported_book) do
        Importer.new(
          reporter:,
          json_content:,
          book_owner_email:,
          default_currency_iso_code:,
          auto_delete_book:,
        ).import!
      end
      let(:reporter) { TestReporter.new }
      let(:json_content) { file_fixture("import/moneydance/with_metadata.json").read }
      let(:owner) { users(:joe) }
      let(:book_owner_email) { owner.email }
      let(:default_currency_iso_code) { "eur" }
      let(:auto_delete_book) { false }

      describe "without metadata" do
        let(:json_content) { "{}" }
        it "raises missing metadata" do
          e = assert_raise(Importer::MissingInformationError) { imported_book }
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
          e = assert_raise(Importer::MissingInformationError) { imported_book }
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
          e = assert_raise(Importer::MissingInformationError) { imported_book }
          assert_equal "Missing export_date (in metadata) in MoneyDance Export", e.message
        end
      end
      describe "no use exists with the email" do
        let(:book_owner_email) { "does_not_exist@example.com" }

        it "raises user cannot be found" do
          e = assert_raise(Importer::DataError) { imported_book }
          assert_equal "No user with E-Mail does_not_exist@example.com found", e.message
        end
      end

      describe "a book already exists" do
        let(:default_currency) { Currency.find_by!(iso_code: "eur") }
        let(:pre_existing_book) { owner.books.create! name: "Foobar (2026-10-31)", default_currency: }
        before { pre_existing_book }

        it "raises that a book already exists" do
          e = assert_raises(Importer::DataError) { imported_book }
          assert_equal "A book with name \"Foobar (2026-10-31)\" already exists.", e.message
          assert_predicate Book.where(id: pre_existing_book), :exists?, "The pre-existing book should not have been deleted"
        end

        describe "auto_delete_book is true" do
          let(:auto_delete_book) { true }

          it "deletes the existing book and creates a new one" do
            refute_equal imported_book.id, pre_existing_book.id, "The import should have created a new book"
            refute_predicate Book.where(id: pre_existing_book), :exists?, "The pre-existing book should have been deleted"
          end
        end
      end

      describe "the currency does not exist" do
        let(:default_currency_iso_code) { "FOO" }

        it "raises no currency with ISO code found" do
          e = assert_raise(Importer::DataError) { imported_book }
          assert_equal "No currency with ISO Code FOO found", e.message
        end
      end
    end
  end
end
