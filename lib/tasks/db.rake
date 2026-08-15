# frozen_string_literal: true

namespace :db do
  namespace :reset do
    desc "Perform a full reset of the database, loads seeds & fixtures, and re-annotate the models"
    task full: [:environment] do
      success = true
      [
        ["Drop Database", :rake, "db:drop"],
        ["Delete schema cache (schema.rb & structure.sql)", :sh, "rm -v -f db/schema.rb db/structure.sql"],
        ["Create Database", :rake, "db:create"],
        ["Migrate Database", :rake, "db:migrate"], # automatically annotate models
        ["Load Fixtures into Database", :rake, "db:fixtures:load"],
        ["Seed Database for development", :rake, "db:seed:dev"]
      ].each do |title, type, command|
        puts "🔹 " + title + (success ? "" : " => 🟡 Skipped because of previous error")
        case type
        when :rake
          puts "bin/rake #{command}"
          Rake::Task[command].invoke if success
        when :sh
          # `sh` outputs the command already
          sh command if success
        end
      rescue => e
        puts "🔴 " + e.message
        success = false
      end
    end
  end

  namespace :seed do
    desc "Load the seed data from db/seeds.rb, and then db/seeds_dev.rb"
    task dev: "db:seed" do
      require_relative "../../db/seeds_dev"
    end
  end
end
