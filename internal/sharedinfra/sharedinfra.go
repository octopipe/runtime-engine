package sharedinfra

import (
	"context"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/pagination"
)

type SharedInfraStatus struct {
	Plugins    []commonv1alpha1.PluginExecutionStatus `json:"plugins"`
	StartedAt  string                                 `json:"startedAt"`
	FinishedAt string                                 `json:"finishedAt"`
	Status     string                                 `json:"status"`
	Error      string                                 `json:"error"`
}

type SharedInfra struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Status    SharedInfraStatus `json:"status"`
	commonv1alpha1.SharedInfraSpec
}

type UseCase interface {
	List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[SharedInfra], error)
	Create(ctx context.Context, sharedInfra SharedInfra) (SharedInfra, error)
	Update(ctx context.Context, sharedInfra SharedInfra) (SharedInfra, error)
	Get(ctx context.Context, name string, namespace string) (SharedInfra, error)
	Reconcile(ctx context.Context, name string, namespace string) error
	Delete(ctx context.Context, name string, namespace string) error
}

type Repository interface {
	List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (commonv1alpha1.SharedInfraList, error)
	Apply(ctx context.Context, s commonv1alpha1.SharedInfra) (commonv1alpha1.SharedInfra, error)
	Get(ctx context.Context, name string, namespace string) (commonv1alpha1.SharedInfra, error)
	Reconcile(ctx context.Context, name string, namespace string) error
	Delete(ctx context.Context, name string, namespace string) error
}
