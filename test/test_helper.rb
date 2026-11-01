ENV["RAILS_ENV"] ||= "test"
require_relative "../config/environment"
require "rails/test_help"
require "minitest/spec"

module ActiveSupport
  class TestCase
    extend Minitest::Spec::DSL # gives TestCase `it`, `let`, `describe`, …

    # Run tests in parallel with specified workers
    parallelize(workers: :number_of_processors)

    # Setup all fixtures in test/fixtures/*.yml for all tests in alphabetical order.
    fixtures :all

    # Add more helper methods to be used by all tests here...
  end
end

# Make top-level `describe` blocks subclass ActiveSupport::TestCase
Minitest::Spec.register_spec_type(ActiveSupport::TestCase) { true }
