package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/portnyagin/practicum_project/internal/app/dto"
	"github.com/portnyagin/practicum_project/internal/app/model"
	"github.com/portnyagin/practicum_project/internal/app/service/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrderService_Save(t *testing.T) {

	type args struct {
		order   *dto.Order
		repoErr error
	}
	type wants struct {
		wantErr bool
		error   error
		anyErr  bool
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "OrderService. Save. Case 1. Empty Order",
			args: args{
				order: nil,
			},
			wants: wants{
				wantErr: true,
				error:   dto.ErrBadParam,
				anyErr:  false,
			},
		},
		{
			name: "OrderService. Save. Case 2. Empty user",
			args: args{
				order: &dto.Order{},
			},
			wants: wants{
				wantErr: true,
				error:   dto.ErrBadParam,
				anyErr:  false,
			},
		},
		{
			name: "OrderService. Save. Case 3. Order registered early",
			args: args{
				order: &dto.Order{Num: "Case 3",
					UserID: 3,
				},
			},
			wants: wants{
				wantErr: true,
				error:   dto.ErrOrderRegistered,
				anyErr:  false,
			},
		},
		{
			name: "OrderService. Save. Case 4. Order registered by another user",
			args: args{
				order: &dto.Order{Num: "Case 4",
					UserID: 4,
				},
			},
			wants: wants{
				wantErr: true,
				error:   dto.ErrOrderRegisteredByAnotherUser,
				anyErr:  false,
			},
		},
		{
			name: "OrderService. Save. Case 5. Unexpected error",
			args: args{
				order: &dto.Order{Num: "Case 5",
					UserID: 5,
				},
			},
			wants: wants{
				wantErr: true,
				anyErr:  true,
			},
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := context.Background()
	orderRepository := mocks.NewMockOrderRepository(mockCtrl)
	target := NewOrderService(orderRepository, log, false)
	orderRepository.EXPECT().GetByNum(ctx, "Case 3").Return(&model.Order{Num: "Case 3", UserID: 3}, nil)
	orderRepository.EXPECT().GetByNum(ctx, "Case 4").Return(&model.Order{Num: "Case 4", UserID: 30}, nil)
	orderRepository.EXPECT().GetByNum(ctx, "Case 5").Return(nil, errors.New("any error"))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := target.Save(ctx, tt.args.order)
			if tt.wants.wantErr {
				if !tt.wants.anyErr {
					assert.ErrorIs(t, err, tt.wants.error, "Expected error is %v, got %v", tt.wants.error, err)
				}
			} else {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wants.wantErr)
			}
		})
	}
}

func TestOrderService_Save2(t *testing.T) {

	type args struct {
		order   *dto.Order
		repoErr error
	}
	type wants struct {
		wantErr bool
		error   error
		anyErr  bool
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "OrderService. Save. Case 6. Successful",
			args: args{
				order: &dto.Order{Num: "Case 6",
					UserID: 6,
				},
			},
			wants: wants{
				wantErr: false,
				anyErr:  true,
			},
		},
		{
			name: "OrderService. Save. Case 7. Unexpected error",
			args: args{
				order: &dto.Order{Num: "Case 7",
					UserID: 7,
				},
			},
			wants: wants{
				wantErr: true,
				anyErr:  true,
			},
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := context.Background()
	orderRepository := mocks.NewMockOrderRepository(mockCtrl)
	target := NewOrderService(orderRepository, log, false)

	orderRepository.EXPECT().GetByNum(ctx, gomock.Any()).Return(nil, &model.NoRowFound).AnyTimes()
	orderRepository.EXPECT().Save(ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, arg *model.Order) error {
			if arg.Num == "Case 6" {
				return nil
			}
			if arg.Num == "Case 7" {
				return errors.New("any error")
			}
			return nil
		},
	).AnyTimes()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := target.Save(ctx, tt.args.order)
			if tt.wants.wantErr {
				if !tt.wants.anyErr {
					assert.ErrorIs(t, err, tt.wants.error, "Expected error is %v, got %v", tt.wants.error, err)
				}
			} else if err != nil {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wants.wantErr)
			}
		})
	}
}
