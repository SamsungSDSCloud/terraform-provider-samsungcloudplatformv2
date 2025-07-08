package tag

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/resourcemanager"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common/filter"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"
)

func ResourceSchema() resourceschema.MapAttribute {
	return resourceschema.MapAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "A map of key-value pairs representing tags for the resource.\n" +
			"  - Keys must be a maximum of 128 characters.\n" +
			"  - Values must be a maximum of 256 characters.",
		Validators: []validator.Map{
			mapvalidator.KeysAre(
				stringvalidator.LengthBetween(1, 128),
			),
			mapvalidator.ValueStringsAre(
				stringvalidator.LengthBetween(0, 256),
			),
		},
	}
}

func DataSourceSchema() datasourceschema.MapAttribute {
	return datasourceschema.MapAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "A map of key-value pairs representing tags for the resource.\n" +
			" - Keys must be a maximum of 128 characters.\n" +
			" - Values must be a maximum of 256 characters.",
		Validators: []validator.Map{
			mapvalidator.KeysAre(
				stringvalidator.LengthBetween(1, 128),
			),
			mapvalidator.ValueStringsAre(
				stringvalidator.LengthBetween(0, 256),
			),
		},
	}
}

func GetTagIndices(clients *client.SCPClient, input interface{}, tagElements map[string]attr.Value, id string) []int {
	ctx := context.Background()

	var filters []filter.Filter
	for k, v := range tagElements {
		tagInfo, _ := clients.ResourceManager.GetTagList(ctx, resourcemanager.TagDataSource{
			Key:   types.StringValue(k),
			Value: types.StringValue(strings.Trim(v.String(), `"`)),
		})
		tagContents := tagInfo.Content
		for _, tagContent := range tagContents {
			f := filter.Filter{
				Name: types.StringValue(id),
				Values: []types.String{
					types.StringValue(common.GetIdFromSrn(tagContent.Srn)),
				},
			}
			filters = append(filters, f)
		}
	}

	if len(filters) == 0 {
		return []int{}
	}

	wrapStructs, _ := filter.WrapStructs(input)
	contents := common.ConvertStructToMaps(wrapStructs)
	contents = filter.ApplyFilter(contents, filters)

	var indices []int
	for _, item := range contents {
		index, _ := common.ToInt(item["index"])
		indices = append(indices, index)
	}

	return indices
}

func getOffering(AuthUrl string) (string, error) {
	regexPattern := `http[s]?://([a-zA-Z0-9-]+)\.([a-zA-Z0-9-]+)\.[a-zA-Z0-9-]+\.com`
	re := regexp.MustCompile(regexPattern)

	matches := re.FindStringSubmatch(AuthUrl)
	if len(matches) > 2 {
		return matches[2], nil
	}
	return "", fmt.Errorf("failed to parse offering")
}

func getPrimaryRegion(clients *client.SCPClient) (string, error) {
	iamClient := clients.Iam
	if iamClient.Config.Region != "" {
		return iamClient.Config.Region, nil
	} else if iamClient.Config.DefaultRegion != "" {
		return iamClient.Config.DefaultRegion, nil
	}
	regionList := iamClient.GetRegionList()
	if len(regionList) > 0 {
		return regionList[0], nil
	}
	return "", fmt.Errorf("failed to get region")
}

func GetSRN(clients *client.SCPClient, serviceName string, resourceType string, resourceIdentifier string) (string, error) {
	offering, err := getOffering(clients.Iam.Config.AuthUrl)
	if err != nil {
		return "", err
	}

	accountId, err := clients.Iam.GetAccountId()
	if err != nil {
		return "", err
	}

	region, err := getPrimaryRegion(clients)
	if err != nil {
		return "", err
	}

	srnFormat := "srn:%s::%s:%s::%s:%s/%s"
	srn := fmt.Sprintf(srnFormat, offering, accountId, region, serviceName, resourceType, resourceIdentifier)
	encodedSrn := common.EncodeBase64(srn)

	return encodedSrn, nil
}

func UpdateTags(clients *client.SCPClient, serviceName string, resourceType string, resourceIdentifier string, tagElements map[string]attr.Value) (types.Map, error) {
	srn, err := GetSRN(clients, serviceName, resourceType, resourceIdentifier)
	if err != nil {
		return types.Map{}, err
	}
	_, err = clients.ResourceManager.BulkUpdateResourceTags(srn, tagElements)
	if err != nil {
		return types.Map{}, err
	}

	tagsMap, _ := types.MapValue(types.StringType, tagElements)
	return tagsMap, nil
}

func GetTags(clients *client.SCPClient, serviceName string, resourceType string, resourceIdentifier string) (types.Map, error) {
	srn, err := GetSRN(clients, serviceName, resourceType, resourceIdentifier)
	if err != nil {
		return types.Map{}, err
	}
	resp, err := clients.ResourceManager.GetResourceTags(srn)

	tags := make(map[string]attr.Value)
	for _, tag := range resp.Content.Tags {
		tags[tag.Key] = types.StringValue(*tag.Value.Get())
	}

	if err != nil {
		return types.Map{}, err
	}

	tagsMap, _ := types.MapValue(types.StringType, tags)
	return tagsMap, nil
}
