package controller

import (
	"encoding/json"
	"fmt"
	"gin-server-skeleton/middleware"
	"gin-server-skeleton/model"
	"gin-server-skeleton/service"
	"gin-server-skeleton/util"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

// controller struct
type UserApiController struct {
	apiVersion string
	Service    *service.UserService
}

// get controller
func (uc *UserApiController) getCtl() *UserApiController {
	var svc *service.UserService
	return &UserApiController{"v1", svc}
}

// do login
func (uc *UserApiController) DoLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	var us *service.UserService
	user, errFind := us.FindUserByEmail(email)

	if user == nil || errFind != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  -1,
			"error": "user not find or find err !",
		})
	} else {
		errPasswd := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if user.Email == email && errPasswd == nil {
			token, err := middleware.GenerateToken(user.ID)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":  -1,
					"error": err.Error(),
				})
			}
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "success",
				"token":   token,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":  -1,
				"error": "email or password error !",
			})
		}
	}
}

// get all users
func (uc *UserApiController) GetAllUsers(c *gin.Context) {
	var currentPageInt, pageSizeInt = util.CURRENT_PAGE, util.PAGE_SIZE
	var totalRows int64
	pageSizeInt = viper.GetInt("PAGE_SIZE")
	currentPage, cpExist := c.GetQuery("currentpage")
	if cpExist {
		currentPageInt, _ = strconv.Atoi(currentPage)
	}

	pageSize, psExist := c.GetQuery("pagesize")
	if psExist {
		pageSizeInt, _ = strconv.Atoi(pageSize)
	}

	// data option setting
	dataOrder, dataOrderExist := c.GetQuery("dataOrder")
	if !dataOrderExist {
		dataOrder = "id desc"
	}

	dataSelect, dataSelectExist := c.GetQuery("dataSelect")
	if !dataSelectExist {
		dataSelect = ""
	}

	dataWhereMap := map[string]interface{}{}
	dataWhere, dataWhereExist := c.GetQuery("dataWhere")
	if dataWhereExist {
		err := json.Unmarshal([]byte(dataWhere), &dataWhereMap)
		if err != nil {
			util.SendError(c, err.Error())
			return
		}
	}

	dataLimitInt := 0
	dataLimit, dataLimitExist := c.GetQuery("dataLimit")
	if dataLimitExist {
		dataLimitInt, _ = strconv.Atoi(dataLimit)
	}

	daoOpt := model.DAOOption{
		Select: dataSelect,
		Order:  dataOrder,
		Where:  dataWhereMap, //map[string]interface{}{},
		Limit:  dataLimitInt,
	}

	users, err := uc.getCtl().Service.FindAllUserByPagesWithKeys(
		map[string]interface{}{},
		map[string]interface{}{},
		currentPageInt,
		pageSizeInt,
		&totalRows,
		daoOpt)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  -1,
			"error": err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"users": users,
	})
}

// get user
func (uc *UserApiController) GetUser(c *gin.Context) {
	userId, ok := c.Get("userId");
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code":  -1,
			"error": "no user",
		})
		return
	}

	idUint64, errConv := strconv.ParseUint(fmt.Sprintf("%v", userId), 10, 64)
	if errConv != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  -1,
			"error": "errConv",
		})
		return
	}

	user, err := uc.getCtl().Service.FindUserById(idUint64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  -1,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"user": user,
	})
}
