require "optparse"

namespace :import do
  def usage_moneydance(reason = nil)
    puts

    if reason.present?
      puts reason
    end

    puts <<~EOF

      Import data from a Moneydance JSON export

      Mandatory parameters:
        - Path to the JSON export file
        - E-Mail address of the owner of the new book
        - The ISO code (e.g.: EUR, USD, CAD) of the default currency of the book to create

      Optional options are:
        - `delete`: drops the book if it already exists

      Example:

        rake "import:moneydance[~/Documents/Moneydance.json, joe@example.com, eur, delete]"
    EOF
    1
  end

  desc "Import data from a Moneydance JSON export"
  task :moneydance, [:source, :owner_email, :currency_iso_code] => [:environment] do |t, args|
    source_path = args.source.to_s.squish
    next usage_moneydance("No source path is given (1st parameter)") if source_path.blank?

    book_owner_email = args.owner_email.to_s.squish
    next usage_moneydance("No owner's E-Mail is given (2nd parameter)") if book_owner_email.blank?

    default_currency_iso_code = args.currency_iso_code.to_s.squish
    next usage_moneydance("No owner's currency ISO code is given (3rd parameter)") if default_currency_iso_code.blank?

    extras = args.extras.map(&:squish)
    auto_delete_book = extras.delete("delete").present?

    next usage_moneydance("Unrecognized parameters: #{args.extras.join(", ")}") unless extras.empty?

    reporter = Import::ConsoleReporter.new
    md_json = JSON.parse(File.read(source_path))

    importer = Import::Moneydance::Importer.new(reporter:, md_json:, book_owner_email:, default_currency_iso_code:, auto_delete_book:)
    importer.import!
  end
end
