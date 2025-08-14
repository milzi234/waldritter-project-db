class UmbrellaProjectsController < ApplicationController
  before_action :set_umbrella_project, only: %i[show update destroy]

  # GET /api/v1/umbrella_projects
  def index
    @umbrella_projects = UmbrellaProject.all
    render json: @umbrella_projects
  end

  # GET /api/v1/umbrella_projects/1
  def show
    render json: @umbrella_project
  end

  # GET /api/v1/umbrella_projects/1/projects
  def projects
    @umbrella_project = UmbrellaProject.includes(:projects).find(params[:id])
    @projects = @umbrella_project.projects
    render json: @projects
  end

  # POST /api/v1/umbrella_projects/1/projects
  def set_projects
    @umbrella_project = UmbrellaProject.find(params[:id])
    @umbrella_project.projects = Project.where(id: params[:project_ids])
    if @umbrella_project.save
      render json: @umbrella_project.projects, status: :ok
    else
      render json: @umbrella_project.errors, status: :unprocessable_entity
    end
  end

  # POST /api/v1/umbrella_projects/1/upload_image
  def upload_image
    @umbrella_project = UmbrellaProject.find(params[:id])
    @umbrella_project.image.attach(params[:image])
    if @umbrella_project.save
      render json: @umbrella_project, status: :ok
    else
      render json: @umbrella_project.errors, status: :unprocessable_entity
    end
  end

  # POST /api/v1/umbrella_projects
  def create
    @umbrella_project = UmbrellaProject.new(umbrella_project_params)

    if @umbrella_project.save
      render json: @umbrella_project, status: :created, location: @umbrella_project
    else
      render json: @umbrella_project.errors, status: :unprocessable_entity
    end
  end

  # PATCH/PUT /api/v1/umbrella_projects/1
  def update
    if @umbrella_project.update(umbrella_project_params)
      render json: @umbrella_project
    else
      render json: @umbrella_project.errors, status: :unprocessable_entity
    end
  end

  # DELETE /api/v1/umbrella_projects/1
  def destroy
    @umbrella_project.destroy
  end

  private
    # Use callbacks to share common setup or constraints between actions.
    def set_umbrella_project
      @umbrella_project = UmbrellaProject.find(params[:id])
    end

    # Only allow a list of trusted parameters through.
    def umbrella_project_params
      params.require(:umbrella_project).permit(:name, :description)
    end
end