class BooksController < ApplicationController
  def index
    @books = Book.where(owner: current_user).page(params[:page])
  end
end
