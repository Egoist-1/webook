package service

import (
	"context"
	"github.com/magiconair/properties/assert"
	"go.uber.org/mock/gomock"
	"testing"
	service2 "webook/_internal/article/_internal/service"
	domain2 "webook/_internal/article/internal/domain"
	"webook/_internal/interactive/internal/domain"
	svcmocks "webook/_internal/internal/service/mocks"
)

func TestRanking(t *testing.T) {
	testCass := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (service2.ArticleService, InteractiveService)
		wantErr  error
		wantArts []domain2.Article
	}{
		{
			name: "成功",
			mock: func(ctrl *gomock.Controller) (service2.ArticleService, InteractiveService) {
				isvc := svcmocks.NewMockInteractiveService(ctrl)
				isvc.EXPECT().TopN(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]domain.Interactive{}, nil)
				asvc := svcmocks.NewMockArticleService(ctrl)
				asvc.EXPECT().TopN(gomock.Any(), gomock.Any()).
					Return([]domain2.Article{}, nil)

				return asvc, isvc
			},
			wantArts: []domain2.Article{
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
