module Import
  class ConsoleReporter
    def initialize
      @active_progress_bar = nil
    end

    def log
      if @active_progress_bar
        @active_progress_bar.puts yield
      else
        puts yield
      end
    end

    def with_progress(title:, total:)
      @active_progress_bar = ProgressBar.new(total, :count, :pourcentage, :elapsed, :rate, :eta)
      yield progress_reporter
    ensure
      @active_progress_bar = nil
    end

    def increment_progress
      @active_progress_bar&.increment! 1
    end
  end
end
