resource "opensearch_composable_index_template" "foo" {
  name = "foo-template"
  body = <<EOF
{
  "index_patterns": ["foo-data-stream*"],
  "data_stream": {}
}
EOF
}

resource "opensearch_data_stream" "foo" {
  name       = "foo-data-stream"
  depends_on = [opensearch_composable_index_template.foo]
}
