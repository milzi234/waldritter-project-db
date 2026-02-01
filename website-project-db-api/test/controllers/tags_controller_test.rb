require "test_helper"

class TagsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @category = categories(:one)
    @tag = tags(:one)
  end

  test "should get index" do
    get category_tags_url(@category), as: :json
    assert_response :success
  end

  test "should create tag" do
    assert_difference("Tag.count") do
      post category_tags_url(@category), params: { tag: { description: @tag.description, title: @tag.title } }, as: :json
    end

    assert_response :created
  end

  test "should show tag" do
    get category_tag_url(@category, @tag), as: :json
    assert_response :success
  end

  test "should update tag" do
    patch category_tag_url(@category, @tag), params: { tag: { description: @tag.description, title: @tag.title } }, as: :json
    assert_response :success
  end

  test "should destroy tag" do
    assert_difference("Tag.count", -1) do
      delete category_tag_url(@category, @tag), as: :json
    end

    assert_response :no_content
  end
end
