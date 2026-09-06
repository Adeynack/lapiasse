# Run using bin/ci

CI.run do
  @steps = ENV.fetch("CI_STEPS", "").split(",").map(&:downcase)

  def group(name, default: false)
    return if !@steps.include?(name) && !@steps.include?("all") && !(default && @steps.empty?)

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

  group "audit" do
    step "Security: Gem audit", "bin/bundler-audit"
    step "Security: Yarn vulnerability audit", "yarn audit"
    step "Security: Importmap vulnerability audit", "bin/importmap audit"
    step "Security: Brakeman code analysis", "bin/brakeman --quiet --no-pager --exit-on-warn --exit-on-error"
  end

  group "test", default: true do
    step "RSpec", "bin/rspec"
  end

  # Optional: set a green GitHub commit status to unblock PR merge.
  # Requires the `gh` CLI and `gh extension install basecamp/gh-signoff`.
  # if success?
  #   step "Signoff: All systems go. Ready for merge and deploy.", "gh signoff"
  # else
  #   failure "Signoff: CI failed. Do not merge or deploy.", "Fix the issues and try again."
  # end
end
