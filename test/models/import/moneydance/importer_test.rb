require "test_helper"

class Import::Moneydance::ImporterTest < ActiveSupport::TestCase
  test "creates a book for the owner with the default currency" do
    assert_difference -> { Book.count }, 1 do
      import
    end

    book = Book.last

    assert_equal(
      [Import::Moneydance::Importer::BOOK_NAME, users(:joe), currencies(:eur)],
      [book.name, book.owner, book.default_currency],
    )
  end

  test "looks up the owner and currency case-insensitively" do
    import(book_owner_email: "JOE@example.com", default_currency: "EUR")

    assert_equal users(:joe), Book.last.owner
  end

  test "raises when the book already exists and auto-delete is off" do
    import

    assert_no_difference -> { Book.count } do
      assert_raises(Import::Moneydance::Importer::BookAlreadyExistsError) { import }
    end
  end

  test "replaces the existing book when auto-delete is on" do
    import
    previous_book = Book.last

    assert_no_difference -> { Book.count } do
      import(auto_delete_book: true)
    end
    assert_not Book.exists?(previous_book.id)
  end

  test "does not touch books of other owners" do
    import(book_owner_email: "mary@example.com")

    assert_difference -> { Book.count }, 1 do
      import
    end
  end

  test "raises on unknown owner" do
    assert_raises(ActiveRecord::RecordNotFound) { import(book_owner_email: "nobody@example.com") }
  end

  test "raises on unknown currency" do
    assert_raises(ActiveRecord::RecordNotFound) { import(default_currency: "XXX") }
  end

  test "raises on invalid JSON" do
    assert_raises(JSON::ParserError) { import(json_content: "not json") }
  end

  private

  def import(**args)
    Import::Moneydance::Importer.new.import(
      logger: Rails.logger,
      json_content: "{}",
      book_owner_email: "joe@example.com",
      default_currency: "eur",
      **args,
    )
  end
end
