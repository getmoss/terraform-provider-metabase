# ===============
#     Group
# ===============
// resource "metabase_permission_group" "example" {
//     name = "created-from-resource"
// }

data "metabase_permission_group" "read_example" {
  group_id = 1
}

output "read_example_name" {
  value = data.metabase_permission_group.read_example
}

# Imported with `TF_LOG=debug tf import -var-file=local.tfvars metabase_permission_group.import_test 2`
resource "metabase_permission_group" "import_test" {
  name = "Administrators"
}
# ===============
#     User
# ===============
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

# Created, then deleted
#resource "metabase_user" "example" {
#  email = "wed-12@example.com"
#  first_name = "John"
#  last_name = "Doe"
#}

# Imported with `TF_LOG=debug tf import -var-file=local.tfvars metabase_user.imported 21`
resource "metabase_user" "imported" {
  email = "moritz.kronberger@getmoss.com"
  first_name = "Moritz"
  last_name = "Kronberger"
}

output "imported_user" {
  value = resource.metabase_user.imported
}