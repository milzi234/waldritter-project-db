Rails.application.routes.draw do
  post "/graphql", to: "graphql#execute"
  if Rails.env.development?
    mount GraphiQL::Rails::Engine, at: "/graphiql", graphql_path: "/graphql"
  end

  
  scope "/api" do
    scope "/v1" do
      resources :projects do
        member do
          get :tags
          get :occurrences
          post :tags, to: "projects#set_tags"
          post :upload_image, to: "projects#upload_image"
        end
        resources :events do
          member do
            get :occurrences
            get :exceptions
            delete 'occurrences/:occurrence_id', to: "events#create_exception"
            delete 'exceptions/:exception_id', to: "events#delete_exception"
          end
        end
      end

      resources :categories do
        resources :tags
      end

      resources :umbrella_projects do
        member do
          get :projects
          post :projects, to: "umbrella_projects#set_projects"
          post :upload_image, to: "umbrella_projects#upload_image"
        end
      end
      
      # search
      put '/search', to: "search#search"

    end
  end
end
