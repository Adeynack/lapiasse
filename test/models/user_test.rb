# == Schema Information
#
# Table name: users
#
#  id                 :integer          not null, primary key
#  display_name       :string           not null
#  email              :string           not null
#  encrypted_password :string           not null
#  created_at         :datetime         not null
#  updated_at         :datetime         not null
#
# Indexes
#
#  index_users_on_email  (email) UNIQUE
#
require "test_helper"

class UserTest < ActiveSupport::TestCase
  test "is valid with an email and password" do
    assert user_with_defaults.valid?
  end

  test "is invalid with an empty email" do
    assert_not user_with_defaults(email: "").valid?
  end

  test "is invalid with an invalid email" do
    assert_not user_with_defaults(email: "abc@").valid?
  end

  test "is invalid with an empty display name" do
    assert_not user_with_defaults(display_name: "").valid?
  end

  test "is invalid with an empty password" do
    assert_not user_with_defaults(password: "").valid?
  end

  private

  def user_with_defaults(**attributes)
    User.new(
      email: "user@example.com",
      password: "password",
      display_name: "Max Mustermann",
      **attributes,
    )
  end
end
