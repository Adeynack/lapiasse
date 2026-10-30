# typed: strict

module Import
  ##
  # The interface to reporting some progress in the import process.
  module Reporter
    extend T::Sig
    extend T::Helpers

    interface!

    ##
    # Report an error.
    #
    # e.g.: Write to the err output
    # e.g.: Log an 'error' entry
    #
    # @param block Block returning the text to report.
    sig { abstract.params(block: T.proc.returns(String)).void }
    def error(&block)
    end

    ##
    # Report a main message.
    #
    # e.g.: Write to the output
    # e.g.: Log an 'info' entry
    #
    # @param block Block returning the text to report.
    sig { abstract.params(block: T.proc.returns(String)).void }
    def log(&block)
    end

    ##
    # Report some detail. Is only actually reported when in verbose mode.
    #
    # e.g.: Write to the output
    # e.g.: Log a 'debug' entry.
    #
    # @param block Block returning the text to report.
    sig { abstract.params(block: T.proc.returns(String)).void }
    def verbose(&block)
    end

    ##
    # Perform a series of action that reports progress.
    #
    # @param title Title of the action.
    # @param total Total number of sub-steps.
    # @param block The action itself, that reports its progress.
    sig do
      abstract.params(
        title: String,
        total: Integer,
        block: T.proc.params(progress: ProgressReporter).void,
      ).void
    end
    def with_progress(title:, total:, &block)
    end
  end

  ##
  # Allows a sub-section of the import to report its progress.
  module ProgressReporter
    extend T::Sig
    extend T::Helpers

    interface!

    ##
    # Signal that one step was made in the actual progress.
    sig { abstract.void }
    def increment!
    end
  end
end
