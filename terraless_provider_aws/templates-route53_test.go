package terraless_provider_aws

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTemplatesFunctions_Route53AliasRecordFor(t *testing.T) {
	// given
	buffer := bytes.Buffer{}

	// when
	buffer = Route53AliasRecordFor("my-domain.com", "my-zone-id", buffer)

	// then
	expected := `resource "aws_route53_record" "terraless-cloudfront-target-my-domain-com" {
  name = "my-domain.com"
  type = "A"
  zone_id = "my-zone-id"

  alias {
    evaluate_target_health = false
    name = "${aws_cloudfront_distribution.terraless-default.domain_name}"
    zone_id = "${aws_cloudfront_distribution.terraless-default.hosted_zone_id}"
  }
}`
	assert.Contains(t, buffer.String(), expected)
}
