package nsqlbproducer

import (
	"context"
	"errors"
	"testing"

	"github.com/Scalingo/go-utils/nsqproducer"
	"github.com/Scalingo/go-utils/nsqproducer/nsqproducermock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockedRandSource struct {
	current int
	values  []int
}

func (m *mockedRandSource) Int() int {
	val := m.values[m.current]
	m.current++
	return val
}

type example struct {
	LBProducer   func([]producer) *NsqLBProducer
	ExpectP1Call bool
	ExpectP2Call bool
	P1Error      error
	P2Error      error
	ExpectError  bool
	RandInt      func() int
}

func randLBProducer(order []int) func(producers []producer) *NsqLBProducer {
	return func(producers []producer) *NsqLBProducer {
		return &NsqLBProducer{
			producers: producers,
			randInt:   (&mockedRandSource{current: 0, values: order}).Int,
		}
	}
}

func TestLBPublish(t *testing.T) {
	examples := map[string]example{
		"when all host works": {
			LBProducer:   randLBProducer([]int{1}),
			ExpectP1Call: false,
			ExpectP2Call: true,
			P1Error:      nil,
			P2Error:      nil,
			ExpectError:  false,
		},
		"when a	single host is down": {
			LBProducer:   randLBProducer([]int{1, 0}),
			ExpectP1Call: true,
			ExpectP2Call: true,
			P1Error:      nil,
			P2Error:      errors.New("NOP"),
			ExpectError:  false,
		},
		"when all hosts are down": {
			LBProducer:   randLBProducer([]int{1, 0}),
			ExpectP1Call: true,
			ExpectP2Call: true,
			P1Error:      errors.New("NOP"),
			P2Error:      errors.New("NOP"),
			ExpectError:  true,
		},
		"when using the fallback mode, the first node ": {
			LBProducer: func(producers []producer) *NsqLBProducer {
				return &NsqLBProducer{
					producers: producers,
					randInt:   alwaysZero,
				}
			},
			ExpectP1Call: true,
			ExpectP2Call: false,
			P1Error:      nil,
		},
		"when using the fallback mode and the firs node is failing, it should call the second one": {
			LBProducer: func(producers []producer) *NsqLBProducer {
				return &NsqLBProducer{
					producers: producers,
					randInt:   alwaysZero,
				}
			},
			ExpectP1Call: true,
			ExpectP2Call: true,
			P1Error:      errors.New("FAIL"),
		},
	}

	for name, example := range examples {
		t.Run(name, func(t *testing.T) {
			t.Run("Publish", func(t *testing.T) {
				runPublishExample(t, example, false)
			})
			t.Run("DeferredPublish", func(t *testing.T) {
				runPublishExample(t, example, true)
			})

		})
	}
}

func runPublishExample(t *testing.T, example example, deferred bool) {
	ctx := context.Background()
	message := nsqproducer.NsqMessageSerialize{}
	topic := "topic"
	delay := int64(0)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	p1 := nsqproducermock.NewMockProducer(ctrl)
	p2 := nsqproducermock.NewMockProducer(ctrl)

	if example.ExpectP1Call {
		if deferred {
			p1.EXPECT().DeferredPublish(ctx, topic, delay, message).Return(example.P1Error)
		} else {
			p1.EXPECT().Publish(ctx, topic, message).Return(example.P1Error)
		}
	}

	if example.ExpectP2Call {
		if deferred {
			p2.EXPECT().DeferredPublish(ctx, topic, delay, message).Return(example.P2Error)
		} else {
			p2.EXPECT().Publish(ctx, topic, message).Return(example.P2Error)
		}
	}

	producer := example.LBProducer([]producer{{producer: p1, host: Host{}}, {producer: p2, host: Host{}}})

	var err error
	if deferred {
		err = producer.DeferredPublish(ctx, topic, delay, message)
	} else {
		err = producer.Publish(ctx, topic, message)
	}

	if example.ExpectError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
}
