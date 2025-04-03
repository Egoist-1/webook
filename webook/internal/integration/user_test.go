package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"start/webook/internal/integration/startup"
	"start/webook/internal/repository/dao"
	"start/webook/internal/web"
	"testing"
	"time"
)

type UserTestSuite struct {
	suite.Suite
	db     *gorm.DB
	server *gin.Engine
}

func (suite *UserTestSuite) SetupSuite() {
	fmt.Println("SetupSuite")
	suite.server = gin.Default()
	suite.db = startup.InitTestDB()
	//hdl := startup.InitUserTest(dao.NewUserDao(suite.db), startup.InitRedis())
	//suite.server.Use(func(ctx *gin.Context) {
	//	ctx.Set("claims", &jwtx.UserClaims{
	//		Id: 1,
	//	})
	//})
	//hdl.RegisterRouter(suite.server)

}

func (suite *UserTestSuite) SetupSubTest() {
	err := suite.db.Exec("TRUNCATE TABLE users").Error
	require.NoError(suite.T(), err)
}

func (suite *UserTestSuite) TestSignup() {
	type ReqSignup struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	t := suite.T()
	testCase := []struct {
		name       string
		before     func(t *testing.T)
		after      func(t *testing.T)
		wantResult Result[int]
		wantCode   int
		req        ReqSignup
	}{
		{
			name: "注册成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var u dao.User
				err := suite.db.First(&u).Error
				require.NoError(t, err)
				require.Equal(t, "chelly", u.Name)
				require.Equal(t, "999999@qq.com", u.Email.String)

			},
			wantResult: Result[int]{
				Msg:  "注册成功",
				Code: 200,
				Data: 1,
			},
			wantCode: 200,
			req: ReqSignup{
				Name:            "chelly",
				Email:           "999999@qq.com",
				Password:        "qweasd",
				ConfirmPassword: "qweasd",
			},
		},
		{
			name: "2",
		},
	}
	for _, tc := range testCase {
		suite.Run(tc.name, func() {
			tc.before(t)
			data, err := json.Marshal(tc.req)
			// 不能有 error
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost,
				"/user/signup", bytes.NewReader(data))
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
func (suite *UserTestSuite) TestEdit() {
	now := time.Now().Unix()
	type Req struct {
		Name    string
		AboutMe string
		Ctime   int64
	}
	t := suite.T()
	testCase := []struct {
		name       string
		before     func(t *testing.T)
		after      func(t *testing.T)
		wantResult Result[web.ProfileVO]
		wantCode   int
		req        Req
	}{
		{
			name: "修改成功",
			before: func(t *testing.T) {
				u := dao.User{
					Name: "chelly2",
					Email: sql.NullString{
						String: "9630@qq.com",
						Valid:  true,
					},
					AboutMe: "旧的简介",
					CTime:   now,
					UTime:   now,
				}
				err := suite.db.Create(u).Error
				require.NoError(t, err)
			},
			after: func(t *testing.T) {

			},
			wantResult: Result[web.ProfileVO]{
				Msg:  "修改成功",
				Code: 200,
				Data: web.ProfileVO{
					Name:    "chelly",
					Email:   "9630@qq.com",
					AboutMe: "新的简介",
					CTime:   now,
				},
			},
			wantCode: 200,
			req: Req{
				Name:    "chelly",
				AboutMe: "新的简介",
			},
		},
		{
			name: "2",
		},
	}
	for _, tc := range testCase {
		suite.Run(tc.name, func() {
			tc.before(t)
			data, err := json.Marshal(tc.req)
			// 不能有 error
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost,
				"/user/edit", bytes.NewReader(data))
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
			var result Result[web.ProfileVO]
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, result)
			tc.after(t)
		})
	}
}

func (suite *UserTestSuite) TestProfile() {
	now := time.Now().Unix()
	t := suite.T()
	testCase := []struct {
		name       string
		before     func(t *testing.T)
		after      func(t *testing.T)
		wantResult Result[web.ProfileVO]
		wantCode   int
	}{
		{
			name: "简介",
			before: func(t *testing.T) {
				u := dao.User{
					Id:   1,
					Name: "chelly",
					Email: sql.NullString{
						String: "9630@qq.com",
						Valid:  true,
					},
					AboutMe: "这是我的简介",
					CTime:   now,
					UTime:   now,
				}
				err := suite.db.Create(u).Error
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				var u dao.User
				err := suite.db.First(&u).Error
				require.NoError(t, err)
				require.Equal(t, "chelly", u.Name)
				require.Equal(t, "9630@qq.com", u.Email.String)
				require.Equal(t, "这是我的简介", u.AboutMe)
				require.Equal(t, now, u.CTime)

			},
			wantResult: Result[web.ProfileVO]{
				Msg:  "简介",
				Code: 200,
				Data: web.ProfileVO{
					Name:    "chelly",
					Email:   "9630@qq.com",
					AboutMe: "这是我的简介",
					CTime:   now,
				},
			},
			wantCode: 200,
		},
	}
	for _, tc := range testCase {
		suite.Run(tc.name, func() {
			tc.before(t)
			//data, err := json.Marshal("")
			// 不能有 error
			//assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost,
				"/user/profile", bytes.NewReader(nil))
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
			var result Result[web.ProfileVO]
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, result)
			tc.after(t)
		})
	}
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
