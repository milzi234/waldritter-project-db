class TagsController < ApplicationController
  before_action :set_tag, only: %i[ show update destroy ]
  before_action :set_category

  # GET  /api/v1/categories/:category_id/tags
  def index
    @tags = @category.tags

    render json: @tags
  end

  # GET  /api/v1/categories/:category_id/tags/1
  def show
    render json: @tag
  end

  # POST  /api/v1/categories/:category_id/tags
  def create
    @tag = @category.tags.build(tag_params)

    if @tag.save
      render json: @tag, status: :created, location: [@category, @tag]
    else
      render json: @tag.errors, status: :unprocessable_entity
    end
  end

  # PATCH/PUT  /api/v1/categories/:category_id/tags/1
  def update
    if @tag.update(tag_params)
      render json: @tag
    else
      render json: @tag.errors, status: :unprocessable_entity
    end
  end

  # DELETE  /api/v1/categories/:category_id/tags/1
  def destroy
    @tag.destroy
  end

  private
    # Use callbacks to share common setup or constraints between actions.
    def set_tag
      @tag = Tag.find(params[:id])
    end

    def set_category
      @category = Category.find(params[:category_id])
    end

    # Only allow a list of trusted parameters through.
    def tag_params
      params.require(:tag).permit(:title, :description, :category_id)
    end
end
