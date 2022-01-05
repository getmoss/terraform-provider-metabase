// resource "metabase_permission_group" "example" {
//     name = "created-from-resource"
// }

data "metabase_permission_group" "read_example" {
    group_id = 1
}

output "read_example_name" {
    value = "${data.metabase_permission_group.read_example.name}"
}