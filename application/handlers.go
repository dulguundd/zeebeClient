package application

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/camunda/zeebe/clients/go/v8/pkg/pb"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
	"net/http"
)

type Handler struct {
	zeebeClient *zbc.ClientConfig
}

func (h Handler) TopologyInfo(w http.ResponseWriter, r *http.Request) {
	zbClient, err := zbc.NewClient(h.zeebeClient)

	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	topology, err := zbClient.NewTopologyCommand().Send(ctx)
	if err != nil {
		panic(err)
	}

	for _, broker := range topology.Brokers {
		fmt.Println("Broker", broker.Host, ":", broker.Port)
		for _, partition := range broker.Partitions {
			fmt.Println("  Partition", partition.PartitionId, ":", roleToString(partition.Role))
		}
	}
}

func (h Handler) DeployResource(w http.ResponseWriter, r *http.Request) {
	var req DeployResourceRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		zbClient, err := zbc.NewClient(h.zeebeClient)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, err)
		} else {
			ctx := context.Background()
			response, err := zbClient.NewDeployResourceCommand().AddResourceFile("order-process-4.bpmn").Send(ctx)
			if err != nil {
				writeResponse(w, http.StatusInternalServerError, err)
			} else {
				apiResponse := DeployResourceResponse{
					key:         response.GetKey(),
					deployments: response.GetDeployments(),
				}
				writeResponse(w, http.StatusOK, apiResponse)
			}
		}
	}
}

func (h Handler) DeployInstance(w http.ResponseWriter, r *http.Request) {
	var req DeployInstaceRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		zbClient, err := zbc.NewClient(h.zeebeClient)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, err)
		} else {
			variables := make(map[string]interface{})
			for i := range req.variables {
				variables[req.variables[i].name] = req.variables[i].value
			}

			request, err := zbClient.NewCreateInstanceCommand().BPMNProcessId(req.bpmnProcessId).LatestVersion().VariablesFromMap(variables)
			if err != nil {
				writeResponse(w, http.StatusInternalServerError, err)
			} else {
				ctx := context.Background()
				result, err := request.Send(ctx)
				if err != nil {
					writeResponse(w, http.StatusInternalServerError, err)
				} else {
					apiResponse := DeployInstanceResponse{
						processDefinitionKey: result.GetProcessDefinitionKey(),
						bpmnProcessId:        result.GetBpmnProcessId(),
						version:              result.GetVersion(),
						processInstanceKey:   result.GetProcessInstanceKey(),
					}
					writeResponse(w, http.StatusOK, apiResponse)
				}
			}
		}
	}
}

func newHandler(zeebeConfig ZeebeConfig) *zbc.ClientConfig {
	return &zbc.ClientConfig{
		GatewayAddress:         zeebeConfig.zeebeAddress,
		UsePlaintextConnection: true,
	}
}

func roleToString(role pb.Partition_PartitionBrokerRole) string {
	switch role {
	case pb.Partition_LEADER:
		return "Leader"
	case pb.Partition_FOLLOWER:
		return "Follower"
	default:
		return "Unknown"
	}
}
