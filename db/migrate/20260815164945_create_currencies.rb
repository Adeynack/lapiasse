class CreateCurrencies < ActiveRecord::Migration[8.1]
  def change
    create_table :currencies do |t|
      t.timestamps
      t.string :iso_code, limit: 3, index: {unique: true}, comment: "Null if it is a custom currency"
      t.string :name, null: false
      t.string :symbol
      t.boolean :symbol_first, null: false, default: false
      t.integer :subunit_to_unit, null: false, default: 100
    end
  end
end
