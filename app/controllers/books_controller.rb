class BooksController < ApplicationController
  def index
    @books = Book.where(owner: current_user)
  end
end
