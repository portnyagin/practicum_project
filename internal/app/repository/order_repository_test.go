package repository

import (
	"context"
	"fmt"
	"github.com/portnyagin/practicum_project/internal/app/model"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestOrderRepositoryImpl_Save(t *testing.T) {
	tests := []struct {
		name    string
		order   model.Order
		wantErr bool
	}{
		{
			name: "OrderRepository. Save. Case #1",
			order: model.Order{
				Id:         0,
				User_id:    1,
				Num:        "11",
				Status:     "STATUS",
				Upload_at:  time.Now().Truncate(time.Microsecond),
				Updated_at: time.Now().Truncate(time.Microsecond),
			},
			wantErr: false,
		},
	}
	target, _ := NewOrderRepository(postgresHandler, Log)
	initDatabase(context.Background(), postgresHandler)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("save object")
			err := target.Save(context.Background(), tt.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println("get object")
			res, err := target.GetByNum(context.Background(), tt.order.Num)
			if err != nil {
				t.Errorf("GetByNum() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println("check got object")
			assert.Equal(t, tt.order.Status, res.Status, "Compare error (order.Status): expected %s,  got %s", tt.order.Status, res.Status)
			assert.Equal(t, tt.order.User_id, res.User_id, "Compare error (order.User_id): expected %s,  got %s", tt.order.User_id, res.User_id)
			assert.Equal(t, tt.order.Updated_at, res.Updated_at, "Compare error (order.Updated_at): expected %s,  got %s", tt.order.Updated_at.String(), res.Updated_at.String())
			assert.Equal(t, tt.order.Upload_at, res.Upload_at, "Compare error (order.Upload_at): expected %s,  got %s", tt.order.Upload_at, res.Upload_at)
		})
	}
}

func TestOrderRepositoryImpl_UpdateStatus(t *testing.T) {
	tests := []struct {
		name      string
		order     model.Order
		statusNew string
		wantErr   bool
	}{
		{
			name: "OrderRepository. Update. Case #1",
			order: model.Order{
				Id:         2,
				User_id:    2,
				Num:        "21",
				Status:     "STATUS",
				Upload_at:  time.Now().Truncate(time.Microsecond),
				Updated_at: time.Now().Truncate(time.Microsecond),
			},
			statusNew: "STATUS_NEW",
			wantErr:   false,
		},
	}
	initDatabase(context.Background(), postgresHandler)
	target, _ := NewOrderRepository(postgresHandler, Log)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("save object")
			if err := target.Save(context.Background(), tt.order); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println("update saved object")
			if err := target.UpdateStatus(context.Background(), tt.order.User_id, tt.order.Num, tt.statusNew); (err != nil) != tt.wantErr {
				t.Errorf("UpdateStatus() error = %v, wantErr %v", err, tt.wantErr)
			}

			fmt.Println("get object")
			res, err := target.GetByNum(context.Background(), tt.order.Num)
			if err != nil {
				t.Errorf("GetByNum() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println("check got object")
			assert.Equal(t, tt.statusNew, res.Status, "Compare error (order.status): expected %s,  got %s", tt.statusNew, res.Status)
		})
	}
}

func TestOrderRepositoryImpl_FindByUser(t *testing.T) {
	tests := []struct {
		name      string
		userID    int
		objCount  int
		objStatus string
		create    bool
		wantErr   bool
	}{
		{
			name:      "OrderRepository. FindByUser. Case #1",
			userID:    31,
			objCount:  2,
			objStatus: "NEW",
			wantErr:   false,
		},
		{
			name:      "OrderRepository. FindByUser. Case #2",
			userID:    32,
			objCount:  0,
			objStatus: "NEW",
			wantErr:   false,
		},
	}
	initDatabase(context.Background(), postgresHandler)
	target, _ := NewOrderRepository(postgresHandler, Log)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("save objects")
			for i := 0; i < tt.objCount; i++ {
				var obj model.Order
				obj.Num = strconv.Itoa(i)
				obj.User_id = tt.userID
				obj.Status = tt.objStatus
				timeLabel := time.Now().Truncate(time.Microsecond)
				obj.Updated_at = timeLabel
				obj.Upload_at = timeLabel

				if err := target.Save(context.Background(), obj); (err != nil) != tt.wantErr {
					t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			fmt.Println("get saved objects")
			resArr, err := target.FindByUser(context.Background(), tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			fmt.Println("check got objects")
			assert.Equal(t, tt.objCount, len(resArr), "Compare error (order.status): expected %s,  got %s", tt.objCount, len(resArr))
		})
	}
}
