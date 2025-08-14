require "test_helper"

class TagsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @tag = tags(:one)
  end

  test "should get index" do
    get tags_url, as: :json
    assert_response :success
  end

  test "should create tag" do
    assert_difference("Tag.count") do
      post tags_url, params: { tag: { category_id: @tag.category_id, description: @tag.description, title: @tag.title } }, as: :json
    end

    assert_response :created
  end

  test "should show tag" do
    get tag_url(@tag), as: :json
    assert_response :success
  end

  test "should update tag" do
    patch tag_url(@tag), params: { tag: { category_id: @tag.category_id, description: @tag.description, title: @tag.title } }, as: :json
    assert_response :success
  end

  test "should destroy tag" do
    assert_difference("Tag.count", -1) do
      delete tag_url(@tag), as: :json
    end

    assert_response :no_content
  end
end
