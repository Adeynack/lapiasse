class HomeController < ApplicationController
  # @route GET / (root)
  def index
    @use_phlex = ActiveModel::Type::Boolean.new.cast(params[:use_phlex])
    if @use_phlex
      render Views::Home::Index.new(current_user:)
    else
      # Use the default Rails mechanism and render the .html.erb view
    end
  end
end
