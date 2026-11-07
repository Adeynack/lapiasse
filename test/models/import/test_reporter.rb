module Import
  class TestReporter
    # Reporter Interface

    def log
      log_entries << yield
    end

    def with_progress(title:, total:)
      yield
    end

    # Debugger Part

    def log_entries
      @log_entries ||= []
    end
  end
end
