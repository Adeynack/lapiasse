# typed: strict

module Import
  module Moneydance
    class Importer < T::Struct
      extend T::Sig

      const :reporter, Import::Reporter
      const :json_content, String
      const :book_owner_email, String
      const :default_currency_iso_code, String
      const :auto_delete_book, T::Boolean, default: false

      sig { void }
      def import!
        md_json = T.let(JSON.parse(json_content), JsonHash)
        create_book!(md_json:)
      end

      private

      sig { params(md_json: JsonHash).returns(Book) }
      def create_book!(md_json:)
        owner = User.find_by(email: book_owner_email) ||
          raise(StandardError, "No user with E-Mail #{book_owner_email} found")

        file_name = T.let(md_json["metadata"]["export_date"], String)
          .presence || raise(StandardError, "No filename found in Moneydance Export")
        export_date = Import::Moneydance::Util.parse_string_date(md_json["metadata"]["export_date"])
        book_name = "#{file_name} (#{export_date.iso8601})"

        existing_book = owner.books.find_by(name: book_name)
        if existing_book
          reporter.log { "Book with name \"#{book_name}\" found" }
          if auto_delete_book
            reporter.log { "Deleting book \"#{book_name}\"" }
            existing_book.destroy!
          else
            raise StandardError, "A book with name \"#{book_name}\" already exists."
          end
        end

        default_currency = Currency.find_by(iso_code: default_currency_iso_code) ||
          raise(StandardError, "No currency with ISO Code #{default_currency_iso_code} found")

        reporter.log { "Creating book with name \"#{book_name}\"" }
        owner.books.create!(name: book_name, default_currency:)
      end
    end
  end
end
