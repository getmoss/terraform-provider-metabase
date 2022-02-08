data "metabase_user" "user_by_email" {
  email = "john.doe@example.com"
}

data "metabase_user" "user_by_id" {
  user_id = 1
}