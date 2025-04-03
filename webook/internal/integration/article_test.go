package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"start/webook/internal/domain"
	"start/webook/internal/e"
	"start/webook/internal/integration/startup"
	"start/webook/internal/repository/dao/article"
	"testing"
)

type ArticleTestSuite struct {
	suite.Suite
	db     *gorm.DB
	server *gin.Engine
}

func (suite *ArticleTestSuite) SetupSuite() {
	suite.server = gin.Default()
	suite.db = startup.InitTestDB()
	//hdl := startup.InitArticleTest(suite.db)
	//suite.server.Use(func(ctx *gin.Context) {
	//	ctx.Set("claims", &jwtx.UserClaims{
	//		Id: 1,
	//	})
	//})
	//hdl.RegisterRouter(suite.server)

}

func (suite *ArticleTestSuite) TearDownSubTest() {
	err := suite.db.Exec("TRUNCATE TABLE articles").Error
	require.NoError(suite.T(), err)
	err = suite.db.Exec("TRUNCATE TABLE article_publishes").Error
	require.NoError(suite.T(), err)
}

func (suite *ArticleTestSuite) TestEdit() {
	t := suite.T()
	type Req struct {
		Id      int
		Title   string
		Content string
	}
	testCase := []struct {
		name       string
		before     func(t *testing.T)
		after      func(t *testing.T)
		wantResult Result[int]
		wantCode   int
		req        Req
	}{
		{
			name: "保存新的帖子成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art article.Article
				err := suite.db.First(&art).Error
				require.NoError(t, err)
				assert.Equal(t, art.Id, 1)
				assert.Equal(t, art.AuthorId, 1)
				assert.Equal(t, art.Title, "新的标题")
				assert.Equal(t, art.Content, "新的内容")
				assert.Equal(t, art.Status, uint(domain.ArticleStatusUnpublished))
			},
			wantResult: Result[int]{
				Code: 200,
				Data: 1,
			},
			wantCode: 200,
			req: Req{
				Title:   "新的标题",
				Content: "新的内容",
			},
		},
		{
			name: "更新",
			before: func(t *testing.T) {
				art := article.Article{
					Id:       1,
					Title:    "旧的标题",
					Content:  "旧的内容",
					AuthorId: 1,
					Ctime:    123,
					Utime:    234,
				}
				err := suite.db.Create(art).Error
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				var art article.Article
				err := suite.db.First(&art).Error
				require.NoError(t, err)
				assert.Equal(t, art.Id, 1)
				assert.Equal(t, art.AuthorId, 1)
				assert.Equal(t, art.Title, "更改的标题")
				assert.Equal(t, art.Content, "更改的内容")
				assert.Equal(t, art.Status, uint(domain.ArticleStatusUnpublished))
				assert.True(t, art.Utime > 234)
			},
			wantResult: Result[int]{
				Code: 200,
				Data: 1,
			},
			wantCode: 200,
			req: Req{
				Id:      1,
				Title:   "更改的标题",
				Content: "更改的内容",
			},
		},
		{
			name: "更新别人的帖子",
			before: func(t *testing.T) {
				art := article.Article{
					Id:       1,
					Title:    "旧的标题",
					Content:  "旧的内容",
					AuthorId: 10,
					Status:   uint(domain.ArticleStatusUnpublished),
					Ctime:    123,
					Utime:    234,
				}
				err := suite.db.Create(art).Error
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				var art article.Article
				err := suite.db.Where("id = ?", 1).First(&art).Error
				require.NoError(t, err)
				assert.Equal(t, art.Id, 1)
				assert.Equal(t, art.AuthorId, 10)
				assert.Equal(t, art.Title, "旧的标题")
				assert.Equal(t, art.Content, "旧的内容")
				assert.Equal(t, art.Status, uint(domain.ArticleStatusUnpublished))
				assert.True(t, art.Utime == 234)
			},
			wantResult: Result[int]{
				Code: e.ServerErr.ToInt(),
				Msg:  e.ServerErr.String(),
			},
			wantCode: 200,
			req: Req{
				Id:      1,
				Title:   "标题",
				Content: "内容",
			},
		},
	}
	for _, tc := range testCase {
		suite.Run(tc.name, func() {
			tc.before(t)
			data, err := json.Marshal(tc.req)
			// 不能有 error
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost,
				"/article/edit", bytes.NewReader(data))
			assert.NoError(t, err)
			req.Header.Set("Content-Type",
				"application/json")
			recorder := httptest.NewRecorder()
			suite.server.ServeHTTP(recorder, req)
			code := recorder.Code
			assert.Equal(t, tc.wantCode, code)
			if code != http.StatusOK {
				return
			}
			// 反序列化为结果
			// 利用泛型来限定结果必须是 int64
			var result Result[int]
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, result)
			tc.after(t)
		})
	}
}
func (suite *ArticleTestSuite) TestPublish() {
	t := suite.T()
	type Req struct {
		Id      int
		Title   string
		Content string
	}
	testCase := []struct {
		name       string
		before     func(t *testing.T)
		after      func(t *testing.T)
		wantResult Result[int]
		wantCode   int
		req        Req
	}{
		{
			name: "新建帖子并发表",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art article.Article
				var artPublish article.Article
				err := suite.db.Where("id = ?", 1).First(&art).Error
				require.NoError(t, err)
				err = suite.db.Where("id = ?", 1).First(&artPublish).Error
				require.NoError(t, err)
				assert.Equal(t, "新的标题", art.Title)
				assert.Equal(t, "新的内容", art.Content)
				assert.Equal(t, 1, art.AuthorId)
				assert.Equal(t, uint(domain.ArticleStatusPublished), art.Status)
				require.Equal(t, art, article.Article(artPublish))
			},
			wantResult: Result[int]{
				Code: 200,
				Data: 1,
			},
			wantCode: 200,
			req: Req{
				Title:   "新的标题",
				Content: "新的内容",
			},
		},
		{
			name: "更新帖子并发表",
			before: func(t *testing.T) {
				art := article.Article{
					Id:       9,
					Title:    "旧的标题",
					Content:  "旧的内容",
					AuthorId: 1,
					Ctime:    123,
					Utime:    234,
				}
				err := suite.db.Create(art).Error
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				var art article.Article
				err := suite.db.Where("id = ?", 9).First(&art).Error
				var artP article.ArticlePublish
				err = suite.db.Where("id = ?", 9).First(&artP).Error
				require.NoError(t, err)
				assert.Equal(t, art.Id, 9)
				assert.Equal(t, art.AuthorId, 1)
				assert.Equal(t, art.Title, "更改的标题")
				assert.Equal(t, art.Content, "更改的内容")
				assert.Equal(t, art.Status, uint(domain.ArticleStatusPublished))
				assert.Equal(t, art.Title, artP.Title)
				assert.Equal(t, art.Content, artP.Content)
				assert.Equal(t, art.Status, artP.Status)
			},
			wantResult: Result[int]{
				Code: 200,
				Data: 9,
			},
			wantCode: 200,
			req: Req{
				Id:      9,
				Title:   "更改的标题",
				Content: "更改的内容",
			},
		},
		{
			name: "更新别人发表的帖子,发表失败",
			before: func(t *testing.T) {
				art := article.Article{
					Id:       9,
					Title:    "旧的标题",
					Content:  "旧的内容",
					AuthorId: 10,
					Status:   uint(domain.ArticleStatusPublished),
					Ctime:    123,
					Utime:    234,
				}
				err := suite.db.Create(art).Error
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				var art article.Article
				err := suite.db.Where("id = ?", 9).First(&art).Error
				require.NoError(t, err)
				assert.Equal(t, art.Id, 9)
				assert.Equal(t, art.AuthorId, 10)
				assert.Equal(t, art.Title, "旧的标题")
				assert.Equal(t, art.Content, "旧的内容")
				assert.Equal(t, art.Status, uint(domain.ArticleStatusPublished))
			},
			wantResult: Result[int]{

				Code: e.ServerErr.ToInt(),
				Msg:  e.ServerErr.String(),
			},
			wantCode: 200,
			req: Req{
				Id:      9,
				Title:   "标题",
				Content: "内容",
			},
		},
	}
	for _, tc := range testCase {
		suite.Run(tc.name, func() {
			tc.before(t)
			data, err := json.Marshal(tc.req)
			// 不能有 error
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost,
				"/article/publish", bytes.NewReader(data))
			assert.NoError(t, err)
			req.Header.Set("Content-Type",
				"application/json")
			recorder := httptest.NewRecorder()
			suite.server.ServeHTTP(recorder, req)
			code := recorder.Code
			assert.Equal(t, tc.wantCode, code)
			if code != http.StatusOK {
				return
			}
			// 反序列化为结果
			// 利用泛型来限定结果必须是 int64
			var result Result[int]
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, result)
			tc.after(t)
		})
	}
}
func (suite *ArticleTestSuite) TestList() {
	t := suite.T()
	type Req struct {
		Id      int
		Title   string
		Content string
	}
	testCase := []struct {
		name       string
		before     func(t *testing.T)
		after      func(t *testing.T)
		wantResult Result[int]
		wantCode   int
		req        Req
	}{}
	for _, tc := range testCase {
		suite.Run(tc.name, func() {
			tc.before(t)
			data, err := json.Marshal(tc.req)
			// 不能有 error
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost,
				"/article/list", bytes.NewReader(data))
			assert.NoError(t, err)
			req.Header.Set("Content-Type",
				"application/json")
			recorder := httptest.NewRecorder()
			suite.server.ServeHTTP(recorder, req)
			code := recorder.Code
			assert.Equal(t, tc.wantCode, code)
			if code != http.StatusOK {
				return
			}
			// 反序列化为结果
			// 利用泛型来限定结果必须是 int64
			var result Result[int]
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, result)
			tc.after(t)
		})
	}
}

func TestArticleTestSuite(t *testing.T) {
	suite.Run(t, new(ArticleTestSuite))
}
