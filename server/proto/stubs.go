package proto

import (
	"context"

	"google.golang.org/grpc"
)

// This file is a TEMPORARY STUB to allow compilation without running protoc.
// It matches the expected output of proper protobuf generation.

type Tenant struct {
	Id               string `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	Domain           string `json:"domain,omitempty"`
	SubscriptionTier string `json:"subscription_tier,omitempty"`
	IsActive         bool   `json:"is_active,omitempty"`
	CreatedAt        string `json:"created_at,omitempty"`
	UpdatedAt        string `json:"updated_at,omitempty"`
}

type GetTenantRequest struct {
	TenantId string `json:"tenant_id,omitempty"`
}

type GetTenantResponse struct {
	Tenant *Tenant `json:"tenant,omitempty"`
}

type ListTenantsRequest struct {
	Page     int32 `json:"page,omitempty"`
	PageSize int32 `json:"page_size,omitempty"`
}

type ListTenantsResponse struct {
	Tenants []*Tenant `json:"tenants,omitempty"`
	Total   int32     `json:"total,omitempty"`
}

type CreateTenantRequest struct {
	Name             string `json:"name,omitempty"`
	Domain           string `json:"domain,omitempty"`
	SubscriptionTier string `json:"subscription_tier,omitempty"`
}

type CreateTenantResponse struct {
	Tenant *Tenant `json:"tenant,omitempty"`
}

type UpdateTenantRequest struct {
	TenantId         string `json:"tenant_id,omitempty"`
	Name             string `json:"name,omitempty"`
	Domain           string `json:"domain,omitempty"`
	SubscriptionTier string `json:"subscription_tier,omitempty"`
}

type UpdateTenantResponse struct {
	Tenant *Tenant `json:"tenant,omitempty"`
}

type DeleteTenantRequest struct {
	TenantId string `json:"tenant_id,omitempty"`
}

type DeleteTenantResponse struct {
	Success bool `json:"success,omitempty"`
}

type AddUserToTenantRequest struct {
	TenantId string `json:"tenant_id,omitempty"`
	UserId   string `json:"user_id,omitempty"`
	Role     string `json:"role,omitempty"`
}

type AddUserToTenantResponse struct {
	Success bool `json:"success,omitempty"`
}

type RemoveUserFromTenantRequest struct {
	TenantId string `json:"tenant_id,omitempty"`
	UserId   string `json:"user_id,omitempty"`
}

type RemoveUserFromTenantResponse struct {
	Success bool `json:"success,omitempty"`
}

// TenantServiceClient is the client API for TenantService.
type TenantServiceClient interface {
	GetTenant(ctx context.Context, in *GetTenantRequest, opts ...grpc.CallOption) (*GetTenantResponse, error)
	ListTenants(ctx context.Context, in *ListTenantsRequest, opts ...grpc.CallOption) (*ListTenantsResponse, error)
	CreateTenant(ctx context.Context, in *CreateTenantRequest, opts ...grpc.CallOption) (*CreateTenantResponse, error)
	UpdateTenant(ctx context.Context, in *UpdateTenantRequest, opts ...grpc.CallOption) (*UpdateTenantResponse, error)
	DeleteTenant(ctx context.Context, in *DeleteTenantRequest, opts ...grpc.CallOption) (*DeleteTenantResponse, error)
	AddUserToTenant(ctx context.Context, in *AddUserToTenantRequest, opts ...grpc.CallOption) (*AddUserToTenantResponse, error)
	RemoveUserFromTenant(ctx context.Context, in *RemoveUserFromTenantRequest, opts ...grpc.CallOption) (*RemoveUserFromTenantResponse, error)
}

type tenantServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTenantServiceClient(cc grpc.ClientConnInterface) TenantServiceClient {
	return &tenantServiceClient{cc}
}

func (c *tenantServiceClient) GetTenant(ctx context.Context, in *GetTenantRequest, opts ...grpc.CallOption) (*GetTenantResponse, error) {
	out := new(GetTenantResponse)
	err := c.cc.Invoke(ctx, "/tenant.TenantService/GetTenant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantServiceClient) ListTenants(ctx context.Context, in *ListTenantsRequest, opts ...grpc.CallOption) (*ListTenantsResponse, error) {
	out := new(ListTenantsResponse)
	err := c.cc.Invoke(ctx, "/tenant.TenantService/ListTenants", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantServiceClient) CreateTenant(ctx context.Context, in *CreateTenantRequest, opts ...grpc.CallOption) (*CreateTenantResponse, error) {
	out := new(CreateTenantResponse)
	err := c.cc.Invoke(ctx, "/tenant.TenantService/CreateTenant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantServiceClient) UpdateTenant(ctx context.Context, in *UpdateTenantRequest, opts ...grpc.CallOption) (*UpdateTenantResponse, error) {
	out := new(UpdateTenantResponse)
	err := c.cc.Invoke(ctx, "/tenant.TenantService/UpdateTenant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantServiceClient) DeleteTenant(ctx context.Context, in *DeleteTenantRequest, opts ...grpc.CallOption) (*DeleteTenantResponse, error) {
	out := new(DeleteTenantResponse)
	err := c.cc.Invoke(ctx, "/tenant.TenantService/DeleteTenant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantServiceClient) AddUserToTenant(ctx context.Context, in *AddUserToTenantRequest, opts ...grpc.CallOption) (*AddUserToTenantResponse, error) {
	out := new(AddUserToTenantResponse)
	err := c.cc.Invoke(ctx, "/tenant.TenantService/AddUserToTenant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantServiceClient) RemoveUserFromTenant(ctx context.Context, in *RemoveUserFromTenantRequest, opts ...grpc.CallOption) (*RemoveUserFromTenantResponse, error) {
	out := new(RemoveUserFromTenantResponse)
	err := c.cc.Invoke(ctx, "/tenant.TenantService/RemoveUserFromTenant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TenantServiceServer is the server API for TenantService.
type TenantServiceServer interface {
	GetTenant(context.Context, *GetTenantRequest) (*GetTenantResponse, error)
	ListTenants(context.Context, *ListTenantsRequest) (*ListTenantsResponse, error)
	CreateTenant(context.Context, *CreateTenantRequest) (*CreateTenantResponse, error)
	UpdateTenant(context.Context, *UpdateTenantRequest) (*UpdateTenantResponse, error)
	DeleteTenant(context.Context, *DeleteTenantRequest) (*DeleteTenantResponse, error)
	AddUserToTenant(context.Context, *AddUserToTenantRequest) (*AddUserToTenantResponse, error)
	RemoveUserFromTenant(context.Context, *RemoveUserFromTenantRequest) (*RemoveUserFromTenantResponse, error)
	mustEmbedUnimplementedTenantServiceServer()
}

// UnimplementedTenantServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTenantServiceServer struct{}

func (UnimplementedTenantServiceServer) GetTenant(context.Context, *GetTenantRequest) (*GetTenantResponse, error) {
	return nil, nil
}
func (UnimplementedTenantServiceServer) ListTenants(context.Context, *ListTenantsRequest) (*ListTenantsResponse, error) {
	return nil, nil
}
func (UnimplementedTenantServiceServer) CreateTenant(context.Context, *CreateTenantRequest) (*CreateTenantResponse, error) {
	return nil, nil
}
func (UnimplementedTenantServiceServer) UpdateTenant(context.Context, *UpdateTenantRequest) (*UpdateTenantResponse, error) {
	return nil, nil
}
func (UnimplementedTenantServiceServer) DeleteTenant(context.Context, *DeleteTenantRequest) (*DeleteTenantResponse, error) {
	return nil, nil
}
func (UnimplementedTenantServiceServer) AddUserToTenant(context.Context, *AddUserToTenantRequest) (*AddUserToTenantResponse, error) {
	return nil, nil
}
func (UnimplementedTenantServiceServer) RemoveUserFromTenant(context.Context, *RemoveUserFromTenantRequest) (*RemoveUserFromTenantResponse, error) {
	return nil, nil
}
func (UnimplementedTenantServiceServer) mustEmbedUnimplementedTenantServiceServer() {}

func RegisterTenantServiceServer(s grpc.ServiceRegistrar, srv TenantServiceServer) {
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: "tenant.TenantService",
		HandlerType: (*TenantServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{MethodName: "GetTenant", Handler: nil},
			{MethodName: "ListTenants", Handler: nil},
			{MethodName: "CreateTenant", Handler: nil},
			{MethodName: "UpdateTenant", Handler: nil},
			{MethodName: "DeleteTenant", Handler: nil},
			{MethodName: "AddUserToTenant", Handler: nil},
			{MethodName: "RemoveUserFromTenant", Handler: nil},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "tenant.proto",
	}, srv)
}
