package service

import (
	"context"
	"github.com/magiconair/properties/assert"
	"go.uber.org/mock/gomock"
	"start/webook/internal/domain"
	svcmocks "start/webook/internal/service/mocks"
	"testing"
)

func TestRanking(t *testing.T) {
	testCass := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (ArticleService, InteractiveService)
		wantErr  error
		wantArts []domain.Article
	}{
		{
			name: "成功",
			mock: func(ctrl *gomock.Controller) (ArticleService, InteractiveService) {
				isvc := svcmocks.NewMockInteractiveService(ctrl)
				isvc.EXPECT().TopN(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]domain.Interactive{}, nil)
				asvc := svcmocks.NewMockArticleService(ctrl)
				asvc.EXPECT().TopN(gomock.Any(), gomock.Any()).
					Return([]domain.Article{}, nil)

				return asvc, isvc
			},
			wantArts: []domain.Article{
				{},
				{},
				{},
			},
		},
	}
	for _, tc := range testCass {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			asvc, isvc := tc.mock(ctrl)
			service := NewRankingService(asvc, isvc)
			arts, err := service.TobN(context.Background())
			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, arts, tc.wantArts)
		})
	}
}
