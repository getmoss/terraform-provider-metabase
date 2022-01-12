# Group
// resource "metabase_permission_group" "example" {
//     name = "created-from-resource"
// }

data "metabase_permission_group" "read_example" {
  group_id = 1
}

output "read_example_name" {
  value = data.metabase_permission_group.read_example
}

resource "metabase_permission_group" "import_test" {
  name = "Administrators"
}

# User
data "metabase_user" "user_by_email" {
  email = "raj@getmoss.com"
}

output "user_by_email" {
  value = data.metabase_user.user_by_email
}

data "metabase_user" "user_by_id" {
  user_id = 1
}

output "user_by_id" {
  value = data.metabase_user.user_by_id
}