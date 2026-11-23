class ApplicationRecord < ActiveRecord::Base
  primary_abstract_class

  scope :random, -> { order("RANDOM()") }
end
