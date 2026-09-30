# Run using bin/ci

require "amazing_print"

CI.run do
  @blacklist, @whitelist = ENV.fetch("CI_STEPS", "")
    .split(",")
    .map(&:downcase)
    .partition { _1.start_with?("-") }

  @blacklist.map! { _1.slice(1..) }

  if @whitelist.length > 0
    puts "Only executing those groups:"
    @whitelist.each { puts "  - #{_1}" }
  end

  if @blacklist.length > 0
    puts "Excluding those groups:"
    @blacklist.each { puts "  - #{_1}" }
  end

  def group(name, default: false)
    return if !@whitelist.include?(name) && !@whitelist.include?("all") && !(default && @whitelist.empty?)
    return if @blacklist.include?(name)

    heading "Running group: #{name}"
    yield
  end

  step "Setup", "bin/setup --skip-server"

  group "lint", default: true do
    step "Style: Ruby", "bin/rubocop"
    step "Format: Prettier", "bin/prettier --check ."
    step "Format: ERB", "bin/erb-format --check"
    step "Style: Markdown", "bin/markdownlint"
  end

  group "annotations", default: true do
    step "Annotation: Routes in Controllers", "bin/chusaku --dry-run --exit-with-error-on-annotation --verbose"
    step "Annotation: Models", "bin/annotaterb models --frozen"
    step "Annotation: Routes", "bin/annotaterb routes --frozen"
  end

  group "audit", default: false do
    step "Security: Gem audit", "bin/bundler-audit"
    step "Security: Yarn vulnerability audit", "yarn audit"
    step "Security: Brakeman code analysis", "bin/brakeman --quiet --no-pager --exit-on-warn --exit-on-error"
  end

  group "test", default: true do
    step "Tests: Rails", "bin/rails test"

    step "Tests: Seeds", "env RAILS_ENV=test bin/rails db:seed:replant"
    step "Tests: System", "bin/rails test:system"
  end

  # Optional: set a green GitHub commit status to unblock PR merge.
  # Requires the `gh` CLI and `gh extension install basecamp/gh-signoff`.
  # if success?
  #   step "Signoff: All systems go. Ready for merge and deploy.", "gh signoff"
  # else
  #   failure "Signoff: CI failed. Do not merge or deploy.", "Fix the issues and try again."
  # end
end
