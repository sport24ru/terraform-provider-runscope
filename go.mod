module github.com/terraform-providers/terraform-provider-runscope

go 1.15

replace github.com/ewilde/go-runscope v0.0.0-20201105143228-6d0bdf051797 => github.com/sport24ru/go-runscope v0.0.0-20210330231158-2b3d4da063ed

require (
	github.com/ewilde/go-runscope v0.0.0-20201105143228-6d0bdf051797
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.4.3
)
