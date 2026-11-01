module Import
  class TestReporter
    def error_entries
      @error_entries ||= []
    end

    def error(&block)
      error_entries << yield
    end

    def log_entries
      @log_entries ||= []
    end

    def log(&block)
      log_entries << yield
    end

    def verbose_entries
      @verbose_entries ||= []
    end

    def verbose(&block)
      verbose_entries << yield
    end

    def with_progress(title:, total:, &block)
      yield(progress_reporter)
    end

    private

    def progress_reporter
      @progress_reporter ||= TestProgressReporter.new
    end
  end
end
