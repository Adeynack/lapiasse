# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# This file is the source Rails uses to define your schema when running `bin/rails
# db:schema:load`. When creating a new database, `bin/rails db:schema:load` tends to
# be faster and is potentially less error prone than running all of your
# migrations from scratch. Old migrations may fail to apply correctly if those
# migrations use external dependencies or application code.
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema[8.1].define(version: 2026_08_15_164946) do
  create_table "books", force: :cascade do |t|
    t.datetime "created_at", null: false
    t.integer "default_currency_id", null: false
    t.string "name", null: false
    t.integer "owner_id", null: false
    t.datetime "updated_at", null: false
    t.index ["default_currency_id"], name: "index_books_on_default_currency_id"
    t.index ["owner_id"], name: "index_books_on_owner_id"
  end

  create_table "currencies", force: :cascade do |t|
    t.datetime "created_at", null: false
    t.string "iso_code", limit: 3
    t.string "name", null: false
    t.integer "subunit_to_unit", default: 100, null: false
    t.string "symbol"
    t.boolean "symbol_first", default: false, null: false
    t.datetime "updated_at", null: false
    t.index ["iso_code"], name: "index_currencies_on_iso_code", unique: true
  end

  create_table "users", force: :cascade do |t|
    t.datetime "created_at", null: false
    t.string "display_name", null: false
    t.string "email", null: false
    t.string "encrypted_password", null: false
    t.datetime "updated_at", null: false
    t.index ["email"], name: "index_users_on_email", unique: true
  end

  add_foreign_key "books", "currencies", column: "default_currency_id"
  add_foreign_key "books", "users", column: "owner_id"
end
