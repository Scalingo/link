package ip

import (
	"context"
	"testing"

	"github.com/Scalingo/link/v2/models"
	"github.com/Scalingo/link/v2/network/networkmock"
	"github.com/golang/mock/gomock"
	"github.com/looplab/fsm"
)

func TestUpdateIP(t *testing.T) {
	examples := map[string]struct {
		oldIP                  models.IP
		newIP                  models.IP
		state                  string
		expectNetworkInterface func(mock *networkmock.MockNetworkInterface)
	}{
		"if the old and new ip has the no no_network flag": {
			oldIP: models.IP{NoNetwork: true},
			newIP: models.IP{NoNetwork: true},
		},
		"if the old and new ip doesn't have the no_network flag": {
			oldIP: models.IP{NoNetwork: false},
			newIP: models.IP{NoNetwork: false},
		},
		"if the user enable the no_network flag": {
			oldIP: models.IP{NoNetwork: false},
			newIP: models.IP{IP: "10.10.10.10/32", NoNetwork: true},
			expectNetworkInterface: func(mock *networkmock.MockNetworkInterface) {
				mock.EXPECT().RemoveIP("10.10.10.10/32")
			},
		},
		"if the user disable the no_network flag and we are not master": {
			oldIP: models.IP{NoNetwork: true},
			newIP: models.IP{IP: "10.10.10.10/32", NoNetwork: false},
			state: STANDBY,
		},
		"if the user disable the no_network flag and we are master": {
			oldIP: models.IP{NoNetwork: true},
			newIP: models.IP{IP: "10.10.10.10/32", NoNetwork: false},
			state: ACTIVATED,
			expectNetworkInterface: func(mock *networkmock.MockNetworkInterface) {
				mock.EXPECT().EnsureIP("10.10.10.10/32")
			},
		},
	}

	for name, example := range examples {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockNetworkInterface := networkmock.NewMockNetworkInterface(ctrl)

			if example.expectNetworkInterface != nil {
				example.expectNetworkInterface(mockNetworkInterface)
			}

			manager := manager{
				ip:               example.oldIP,
				stateMachine:     fsm.NewFSM(example.state, fsm.Events{}, fsm.Callbacks{}),
				networkInterface: mockNetworkInterface,
			}

			manager.UpdateIP(context.Background(), example.newIP)
		})
	}
}
