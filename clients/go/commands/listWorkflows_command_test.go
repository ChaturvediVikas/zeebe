package commands

import (
	"github.com/golang/mock/gomock"
	"github.com/zeebe-io/zeebe/clients/go/mock_pb"
	"github.com/zeebe-io/zeebe/clients/go/pb"
	"testing"
)

func TestListWorkflowsCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	request := &pb.ListWorkflowsRequest{}
	stub := &pb.ListWorkflowsResponse{}

	client.EXPECT().ListWorkflows(gomock.Any(), &rpcMsg{msg: request}).Return(stub, nil)

	command := NewListWorkflowsCommand(client)

	response, err := command.Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestListWorkflowsCommandWithBpmnProcessId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	request := &pb.ListWorkflowsRequest{
		BpmnProcessId: "foo",
	}
	stub := &pb.ListWorkflowsResponse{}

	client.EXPECT().ListWorkflows(gomock.Any(), &rpcMsg{msg: request}).Return(stub, nil)

	command := NewListWorkflowsCommand(client)

	response, err := command.BpmnProcessId("foo").Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}