class ProjectsController < ApplicationController
  before_action :set_project, only: %i[ show occurrences update upload_image destroy ]

  # GET /api/v1/projects
  def index
    @projects = Project.all

    render json: @projects
  end

  # GET /api/v1/projects/1
  def show
    if @project.image.attached?
      render json: @project.as_json.merge({ image: url_for(@project.image) })
    else
      render json: @project
    end
  end

  # GET /api/v1/projects/1/occurrences?start_date=2021-01-01&end_date=2021-12-31&limit=100
  def occurrences
    @occurrences = @project.occurrences.where(start_date: params[:start_date]..params[:end_date]).order(:start_date).limit(params[:limit] || 100)
    render json: @occurrences
  end


  # GET /api/v1/projects/1/tags
  def tags
    @project = Project.includes(:tags).find(params[:id])
    @tags = @project.tags
    render json: @tags
  end

  # POST /api/v1/projects/1/tags
  def set_tags
    @project = Project.find(params[:id])
    @project.tags = Tag.where(id: params[:tag_ids])
    render json: @tags
  end

  # POST /api/v1/projects
  def create
    @project = Project.new(project_params)

    if @project.save
      render json: @project, status: :created, location: @project
    else
      render json: @project.errors, status: :unprocessable_entity
    end
  end

  # PATCH/PUT /api/v1/projects/1
  def update
    if @project.update(project_params)
      render json: @project
    else
      render json: @project.errors, status: :unprocessable_entity
    end
  end

  # POST /api/v1/projects/1/upload_image
  def upload_image
    @project.image.attach(params[:image])
    if @project.save
      render json: @project, status: :ok
    else
      render json: @project.errors, status: :unprocessable_entity
    end
  end

  # DELETE /api/v1/projects/1
  def destroy
    @project.destroy
  end

  private
    # Use callbacks to share common setup or constraints between actions.
    def set_project
      @project = Project.find(params[:id])
    end

    # Only allow a list of trusted parameters through.
    def project_params
      params.require(:project).permit(:title, :description, :image)
    end
end
