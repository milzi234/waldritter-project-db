class EventsController < ApplicationController
  before_action :set_event, only: %i[ show update destroy occurrences exceptions create_exception delete_exception]
  before_action :set_project

  # GET  /api/v1/events/:event_id/events
  def index
    @events = @project.events.order(:start_date)

    render json: @events
  end

  # GET  /api/v1/events/:event_id/events/1
  def show
    render json: @event
  end

  # POST  /api/v1/events/:event_id/events
  def create
    @event = @project.events.build(event_params)

    if @event.save
      @event.recalculate_occurrences
      render json: @event, status: :created, location: [@project, @event]
    else
      render json: @event.errors, status: :unprocessable_entity
    end
  end

  # PATCH/PUT  /api/v1/events/:event_id/events/1
  def update
    original_start_date = @event.start_date.to_datetime
    if @event.update(event_params)
      if @event.saved_change_to_attribute?(:recurrence_type) || original_start_date.day != @event.start_date.day
        @event.exceptions.destroy_all
      elsif original_start_date.hour != @event.start_date.hour || original_start_date.min != @event.start_date.min
        @event.exceptions.each do |exception|
          exception.start_date = exception.start_date.change(hour: @event.start_date.hour, min: @event.start_date.min)
          exception.save
        end
      end
      @event.recalculate_occurrences
      render json: @event
    else
      render json: @event.errors, status: :unprocessable_entity
    end
  end

  # DELETE  /api/v1/events/:event_id/events/1
  def destroy
    @event.destroy
  end

  # GET /api/v1/events/:event_id/occurrences?start_date=2021-01-01&end_date=2021-12-31&limit=100
  def occurrences
    @occurrences = @event.occurrences.where(start_date: params[:start_date]..params[:end_date]).order(:start_date).limit(params[:limit] || 100)
    render json: @occurrences
  end

  # GET /api/v1/events/:event_id/exceptions?start_date=2021-01-01&end_date=2021-12-31&limit=100
  def exceptions
    @exceptions = @event.exceptions.where(start_date: params[:start_date]..params[:end_date]).order(:start_date).limit(params[:limit] || 100)
    render json: @exceptions
  end

  # DELETE /api/v1/events/:event_id/occurrences/:occurrence_id
  def create_exception 
    @occurrence = @event.occurrences.find(params[:occurrence_id])
    @event.exceptions.create(start_date: @occurrence.start_date, end_date: @occurrence.end_date)
    @event.recalculate_occurrences
    render json: @event
  end

  # DELETE /api/v1/events/:event_id/exceptions/:exception_id
  def delete_exception
    @exception = @event.exceptions.find(params[:exception_id])
    @exception.destroy
    @event.recalculate_occurrences
    render json: @event
  end

  private
    # Use callbacks to share common setup or constraints between actions.
    def set_event
      @event = Event.find(params[:id])
    end

    def set_project
      @project = Project.find(params[:project_id])
    end

    # Only allow a list of trusted parameters through.
    def event_params
      params.require(:event).permit(:start_date, :end_date, :recurrence_type, :project_id)
    end
end
