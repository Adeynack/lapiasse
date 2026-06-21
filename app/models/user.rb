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
class User < ApplicationRecord
  devise :database_authenticatable

  validates :email,
    presence: true,
    format: {with: /\A(.+)@(.+)\z/, message: "has invalid format"},
    uniqueness: {case_sensitive: false},
    length: {minimum: 4, maximum: 254}
  validates :display_name, presence: true
  validates :encrypted_password, presence: true
end
