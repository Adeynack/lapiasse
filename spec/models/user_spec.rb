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
require 'rails_helper'

RSpec.describe User do
  def user_with_defaults(**attr)
    User.new(
      email: "user@example.com",
      password: "password",
      display_name: "Max Mustermann",
      **attr,
    )
  end

  describe "is valid" do
    it "with an email and password" do
      expect(user_with_defaults).to be_valid
    end
  end

  describe "is invalid" do
    it "with an empty email" do
      expect(user_with_defaults(email: "")).not_to be_valid
    end

    it "with an invalid email" do
      expect(user_with_defaults(email: "abc@")).not_to be_valid
    end

    it "with an empty display name" do
      expect(user_with_defaults(display_name: "")).not_to be_valid
    end

    it "with an empty password" do
      expect(user_with_defaults(password: "")).not_to be_valid
    end
  end
end
