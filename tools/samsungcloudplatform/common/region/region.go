package region

import (
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func ResourceSchema() resourceschema.StringAttribute {
	return resourceschema.StringAttribute{
		Description: "Region",
		Optional:    true,
	}
}

func DataSourceSchema() datasourceschema.StringAttribute {
	return datasourceschema.StringAttribute{
		Description: "Region",
		Optional:    true,
	}
}
