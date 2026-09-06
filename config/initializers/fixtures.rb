# Helpers available to the ERB of every fixture file. `context_class` is the object fixture ERB is
# rendered against, so anything included here can be called directly from a `.yml` fixture.
#
# This lives outside of tests because fixtures are also loaded in development from the CLI.
module FixtureHelpers
  # Hashes a plain-text password the same way Devise does when assigning `User#password=`, so
  # fixtures can keep the password readable:
  #
  #   encrypted_password: <%= encrypt_user_password("joe") %>
  #
  # Reading `stretches` off the model at render time is what makes this work in both environments:
  # cost 1 in test for speed, cost 12 in development. `User` is baked into the name because
  # `stretches` and `pepper` are per-model settings; another authenticatable model needs its own
  # helper rather than a reuse of this one.
  def encrypt_user_password(password)
    Devise::Encryptor.digest(User, password)
  end
end

ActiveSupport.on_load(:active_record) do
  ActiveRecord::FixtureSet.context_class.include FixtureHelpers
end
