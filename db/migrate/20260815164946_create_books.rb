class CreateBooks < ActiveRecord::Migration[8.1]
  def change
    create_table :books do |t|
      t.timestamps
      t.string :name, null: false
      t.references :owner, foreign_key: {to_table: :users}, null: false
      t.references :default_currency, foreign_key: {to_table: :currencies}, null: false
    end
  end
end
