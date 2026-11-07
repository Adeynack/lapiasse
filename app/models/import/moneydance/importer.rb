module Import
  module Moneydance
    class Importer
      ##
      # md_json: an already parsed JSON content (Hash).
      def initialize(reporter:, md_json:, book_owner_email:, default_currency_iso_code:, auto_delete_book:)
        @reporter = reporter
        @md_json = md_json
        @book_owner_email = book_owner_email
        @default_currency_iso_code = default_currency_iso_code
        @auto_delete_book = auto_delete_book
      end

      class MissingInformationError < StandardError
        def initialize(missing_information)
          @missing_information = missing_information
        end

        def message = "Missing #{@missing_information} in MoneyDance Export"
      end

      class DataError < StandardError; end

      def import!
        create_book!
      end

      private

      attr_reader :reporter, :md_json, :book_owner_email, :default_currency_iso_code, :auto_delete_book

      def create_book!
        owner = User.find_by(email: book_owner_email) || raise(DataError, "No user with E-Mail #{book_owner_email} found")
        metadata = md_json["metadata"] || raise(MissingInformationError.new("metadata"))
        file_name = metadata["file_name"].to_s.presence || raise(MissingInformationError.new("filename (in metadata)"))
        raw_export_date = metadata["export_date"].to_s.presence || raise(MissingInformationError.new("export_date (in metadata)"))
        export_date = Util.parse_string_date(raw_export_date)
        book_name = "#{file_name} (#{export_date.iso8601})"

        existing_book = owner.books.find_by(name: book_name)
        if existing_book
          reporter.log { "Book with name \"#{book_name}\" found" }
          if auto_delete_book
            reporter.log { "Deleting book \"#{book_name}\"" }
            existing_book.destroy!
          else
            raise DataError, "A book with name \"#{book_name}\" already exists."
          end
        end

        default_currency = Currency.find_by(iso_code: default_currency_iso_code) ||
          raise(DataError, "No currency with ISO Code #{default_currency_iso_code} found")

        reporter.log { "Creating book with name \"#{book_name}\"" }
        book = owner.books.create!(name: book_name, default_currency:)
        reporter.log { "Created book with ID #{book.id}" }

        book
      end
    end
  end
end
