# == Schema Information
#
# Table name: users
#
#  id                 :integer          not null, primary key
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
  describe "is valid" do
    it "with an email and password" do
      user = User.new(email: "user@example.com", password: "password")
      expect(user).to be_valid
    end
  end

  describe "is invalid" do
    it "with an empty email" do
      user = User.new(email: "", password: "password")
      expect(user).not_to be_valid
    end
  end
end
