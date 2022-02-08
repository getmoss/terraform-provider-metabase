resource "metabase_permission_group" "admins" {
  name = "Administrators"
}

resource "metabase_user" "example" {
  email      = "john.doe@example.com"
  first_name = "John"
  last_name  = "Doe"
}

resource "metabase_membership" "example" {
  group_id = metabase_permission_group.admins.group_id
  user_id  = metabase_user.example.user_id
}