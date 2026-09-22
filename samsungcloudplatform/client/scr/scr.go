package scr

import (
	"context"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scr11 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/scr/1.1"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *scr11.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: scr11.NewAPIClient(config),
	}
}

// ============================================================
// Container Registry APIs
// ============================================================

func (c *Client) CheckRegistryNameDuplication(ctx context.Context, name string) (*scr11.CheckDuplicationResponse, error) {
	req := c.sdkClient.ScrV1ContainerRegistryApisAPI.CheckRegistryNameDuplication(ctx)
	req = req.Name(name)
	resp, _, err := req.Execute()
	return resp, err
}

func (c *Client) CreateRegistry(ctx context.Context, req scr11.ContainerRegistryCreateRequest) (*scr11.ContainerRegistryCreateResponse, error) {
	r := c.sdkClient.ScrV1ContainerRegistryApisAPI.CreateRegistry(ctx)
	r = r.ContainerRegistryCreateRequest(req)
	resp, _, err := r.Execute()
	return resp, err
}

func (c *Client) DeleteRegistry(ctx context.Context, registryId string) error {
	req := c.sdkClient.ScrV1ContainerRegistryApisAPI.DeleteRegistry(ctx, registryId)
	_, err := req.Execute()
	return err
}

func (c *Client) ShowRegistry(ctx context.Context, registryId string) (*scr11.ContainerRegistryShowResponse, error) {
	req := c.sdkClient.ScrV1ContainerRegistryApisAPI.ShowRegistry(ctx, registryId)
	resp, _, err := req.Execute()
	return resp, err
}

func (c *Client) ListRegistries(ctx context.Context, name string) (*scr11.ContainerRegistryListResponse, error) {
	req := c.sdkClient.ScrV1ContainerRegistryApisAPI.ListRegistries(ctx)
	if name != "" {
		req = req.Name(name)
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (c *Client) ListRepositories(ctx context.Context, registryId string, name string, sort string, page int32, size int32) (*scr11.RepositoryListResponse, error) {
	req := c.sdkClient.ScrV1ContainerRegistryApisAPI.ListRepositories(ctx, registryId)
	if name != "" {
		req = req.Name(name)
	}
	if sort != "" {
		req = req.Sort(sort)
	}
	if page > 0 {
		req = req.Page(page)
	}
	if size > 0 {
		req = req.Size(size)
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (c *Client) ListConnectableResources(ctx context.Context) (*scr11.ConnectableResourceListResponse, error) {
	req := c.sdkClient.ScrV1ContainerRegistryApisAPI.ListConnectableResources(ctx)
	resp, _, err := req.Execute()
	return resp, err
}

func (c *Client) UpdatePrivateAcl(ctx context.Context, registryId string, req scr11.PrivateAclSetRequest) error {
	r := c.sdkClient.ScrV1ContainerRegistryApisAPI.UpdatePrivateAcl(ctx, registryId)
	r = r.PrivateAclSetRequest(req)
	_, err := r.Execute()
	return err
}

func (c *Client) UpdatePublicAcl(ctx context.Context, registryId string, req scr11.PublicAclSetRequest) error {
	r := c.sdkClient.ScrV1ContainerRegistryApisAPI.UpdatePublicAcl(ctx, registryId)
	r = r.PublicAclSetRequest(req)
	_, err := r.Execute()
	return err
}

func (c *Client) UpdatePublicEndpointEnabled(ctx context.Context, registryId string, req scr11.PublicEndpointEnabledSetRequest) error {
	r := c.sdkClient.ScrV1ContainerRegistryApisAPI.UpdatePublicEndpointEnabled(ctx, registryId)
	r = r.PublicEndpointEnabledSetRequest(req)
	_, err := r.Execute()
	return err
}

// ============================================================
// Repository APIs
// ============================================================

func (c *Client) CheckRepositoryNameDuplication(ctx context.Context, registryId, name string) (*scr11.CheckDuplicationResponse, error) {
	req := c.sdkClient.ScrV1RepositoriesApiAPI.CheckRepositoryNameDuplication(ctx)
	req = req.RegistryId(registryId)
	req = req.Name(name)
	resp, _, err := req.Execute()
	return resp, err
}

func (c *Client) CreateRepository(ctx context.Context, req scr11.RepositoryCreateRequestV11) (*scr11.RepositoryCreateResponse, error) {
	r := c.sdkClient.ScrV1RepositoriesApiAPI.CreateRepository(ctx)
	r = r.RepositoryCreateRequestV11(req)
	resp, _, err := r.Execute()
	return resp, err
}

func (c *Client) DeleteRepository(ctx context.Context, repositoryId string) error {
	req := c.sdkClient.ScrV1RepositoriesApiAPI.DeleteRepository(ctx, repositoryId)
	_, err := req.Execute()
	return err
}

func (c *Client) ShowRepository(ctx context.Context, repositoryId string) (*scr11.RepositoryShowResponseV11, error) {
	req := c.sdkClient.ScrV1RepositoriesApiAPI.ShowRepository(ctx, repositoryId)
	resp, _, err := req.Execute()
	return resp, err
}

func (c *Client) ListImages(ctx context.Context, repositoryId string, name string, sort string, page int32, size int32) (*scr11.ImageListResponse, error) {
	req := c.sdkClient.ScrV1RepositoriesApiAPI.ListImages(ctx, repositoryId)
	if name != "" {
		req = req.Name(name)
	}
	if sort != "" {
		req = req.Sort(sort)
	}
	if page >= 0 {
		req = req.Page(page)
	}
	if size > 0 {
		req = req.Size(size)
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (c *Client) UpdateRepositoryDescription(ctx context.Context, repositoryId string, req scr11.RepositorySetRequest) error {
	r := c.sdkClient.ScrV1RepositoriesApiAPI.UpdateRepositoryDescription(ctx, repositoryId)
	r = r.RepositorySetRequest(req)
	_, err := r.Execute()
	return err
}

func (c *Client) UpdateRepositoryLifecyclePolicy(ctx context.Context, repositoryId string, req scr11.LifecyclePolicySetRequestV11) error {
	r := c.sdkClient.ScrV1RepositoriesApiAPI.UpdateRepositoryLifecyclePolicy(ctx, repositoryId)
	r = r.LifecyclePolicySetRequestV11(req)
	_, err := r.Execute()
	return err
}

func (c *Client) UpdateRepositoryLockPolicy(ctx context.Context, repositoryId string, req scr11.LockPolicySetRequest) error {
	r := c.sdkClient.ScrV1RepositoriesApiAPI.UpdateRepositoryLockPolicy(ctx, repositoryId)
	r = r.LockPolicySetRequest(req)
	_, err := r.Execute()
	return err
}

func (c *Client) UpdateRepositoryPullPolicy(ctx context.Context, repositoryId string, req scr11.PullPolicySetRequest) error {
	r := c.sdkClient.ScrV1RepositoriesApiAPI.UpdateRepositoryPullPolicy(ctx, repositoryId)
	r = r.PullPolicySetRequest(req)
	_, err := r.Execute()
	return err
}

func (c *Client) UpdateRepositoryScanPolicy(ctx context.Context, repositoryId string, req scr11.ScanPolicySetRequest) error {
	r := c.sdkClient.ScrV1RepositoriesApiAPI.UpdateRepositoryScanPolicy(ctx, repositoryId)
	r = r.ScanPolicySetRequest(req)
	_, err := r.Execute()
	return err
}

// ============================================================
// Image APIs
// ============================================================

func (c *Client) ShowImage(ctx context.Context, imageId string) (*scr11.ImageShowResponseV11, error) {
	req := c.sdkClient.ScrV1ImagesApiAPI.ShowImage(ctx, imageId)
	resp, _, err := req.Execute()
	return resp, err
}

func (c *Client) DeleteImage(ctx context.Context, imageId string) error {
	req := c.sdkClient.ScrV1ImagesApiAPI.DeleteImage(ctx, imageId)
	_, err := req.Execute()
	return err
}

func (c *Client) UpdateImageDescription(ctx context.Context, imageId string, req scr11.ImageSetRequest) error {
	r := c.sdkClient.ScrV1ImagesApiAPI.UpdateImageDescription(ctx, imageId)
	r = r.ImageSetRequest(req)
	_, err := r.Execute()
	return err
}

func (c *Client) UpdateImageLifecyclePolicy(ctx context.Context, imageId string, req scr11.LifecyclePolicySetRequestV11) error {
	r := c.sdkClient.ScrV1ImagesApiAPI.UpdateImageLifecyclePolicy(ctx, imageId)
	r = r.LifecyclePolicySetRequestV11(req)
	_, err := r.Execute()
	return err
}

func (c *Client) UpdateImageLockPolicy(ctx context.Context, imageId string, req scr11.LockPolicySetRequest) error {
	r := c.sdkClient.ScrV1ImagesApiAPI.UpdateImageLockPolicy(ctx, imageId)
	r = r.LockPolicySetRequest(req)
	_, err := r.Execute()
	return err
}

func (c *Client) UpdateImagePullPolicy(ctx context.Context, imageId string, req scr11.PullPolicySetRequest) error {
	r := c.sdkClient.ScrV1ImagesApiAPI.UpdateImagePullPolicy(ctx, imageId)
	r = r.PullPolicySetRequest(req)
	_, err := r.Execute()
	return err
}

func (c *Client) UpdateImageScanPolicy(ctx context.Context, imageId string, req scr11.ScanPolicySetRequest) error {
	r := c.sdkClient.ScrV1ImagesApiAPI.UpdateImageScanPolicy(ctx, imageId)
	r = r.ScanPolicySetRequest(req)
	_, err := r.Execute()
	return err
}

// ============================================================
// Tag APIs
// ============================================================

func (c *Client) ShowTags(ctx context.Context, tagsId string) (*scr11.TagsShowResponse, error) {
	req := c.sdkClient.ScrV1TagsApiAPI.ShowTags(ctx, tagsId)
	resp, _, err := req.Execute()
	return resp, err
}

type TagPackagesQueryOptions struct {
	OsLanguage  string
	PackageName string
	Sort        string
	Page        *int32
	Size        *int32
}

func (c *Client) ShowTagsPackages(ctx context.Context, tagsId string, opts TagPackagesQueryOptions) (*scr11.PackageListResponse, error) {
	req := c.sdkClient.ScrV1TagsApiAPI.ShowTagsPackages(ctx, tagsId)
	if opts.OsLanguage != "" {
		req = req.OsLanguage(opts.OsLanguage)
	}
	if opts.PackageName != "" {
		req = req.PackageName(opts.PackageName)
	}
	if opts.Sort != "" {
		req = req.Sort(opts.Sort)
	}
	if opts.Page != nil {
		req = req.Page(*opts.Page)
	}
	if opts.Size != nil {
		req = req.Size(*opts.Size)
	}
	resp, _, err := req.Execute()
	return resp, err
}

type TagSecretsQueryOptions struct {
	Sort     string
	Page     *int32
	Size     *int32
	FileName string
}

func (c *Client) ShowTagsSecrets(ctx context.Context, tagsId string, opts TagSecretsQueryOptions) (*scr11.SecretListResponse, error) {
	req := c.sdkClient.ScrV1TagsApiAPI.ShowTagsSecrets(ctx, tagsId)
	if opts.Sort != "" {
		req = req.Sort(opts.Sort)
	}
	if opts.Page != nil {
		req = req.Page(*opts.Page)
	}
	if opts.Size != nil {
		req = req.Size(*opts.Size)
	}
	if opts.FileName != "" {
		req = req.FileName(opts.FileName)
	}
	resp, _, err := req.Execute()
	return resp, err
}

type TagVulnerabilityQueryOptions struct {
	OsLanguage             string
	PackageName            string
	Sort                   string
	Page                   *int32
	Size                   *int32
	UpdateVersionAvailable *bool
	Severity               string
	Category               string
}

func (c *Client) ShowTagsVulnerabilities(ctx context.Context, tagsId string, opts TagVulnerabilityQueryOptions) (*scr11.VulnerabilityListResponse, error) {
	req := c.sdkClient.ScrV1TagsApiAPI.ShowTagsVulnerabilities(ctx, tagsId)
	if opts.OsLanguage != "" {
		req = req.OsLanguage(opts.OsLanguage)
	}
	if opts.PackageName != "" {
		req = req.PackageName(opts.PackageName)
	}
	if opts.Sort != "" {
		req = req.Sort(opts.Sort)
	}
	if opts.Page != nil {
		req = req.Page(*opts.Page)
	}
	if opts.Size != nil {
		req = req.Size(*opts.Size)
	}
	if opts.UpdateVersionAvailable != nil {
		req = req.UpdateVersionAvailable(*opts.UpdateVersionAvailable)
	}
	if opts.Severity != "" {
		req = req.Severity(scr11.Severity{String: &opts.Severity})
	}
	if opts.Category != "" {
		req = req.Category(opts.Category)
	}
	resp, _, err := req.Execute()
	return resp, err
}

func (c *Client) ListTagses(ctx context.Context, imageId string, sort string, page *int32, size *int32, referenceTags string) (*scr11.TagsListResponse, error) {
	req := c.sdkClient.ScrV1ImagesApiAPI.ListTagses(ctx, imageId)
	if sort != "" {
		req = req.Sort(sort)
	}
	if page != nil {
		req = req.Page(*page)
	}
	if size != nil {
		req = req.Size(*size)
	}
	if referenceTags != "" {
		req = req.ReferenceTags(referenceTags)
	}
	resp, _, err := req.Execute()
	return resp, err
}
