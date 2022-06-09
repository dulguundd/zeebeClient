package application

import "github.com/camunda/zeebe/clients/go/v8/pkg/pb"

type DeployResourceRequest struct {
	path string
}

type DeployResourceResponse struct {
	key         int64
	deployments []*pb.Deployment
}

type DeployInstaceRequest struct {
	bpmnProcessId string
	variables     []Variable
}

type Variable struct {
	name  string
	value string
}

type DeployInstanceResponse struct {
	processDefinitionKey int64
	bpmnProcessId        string
	version              int32
	processInstanceKey   int64
}
