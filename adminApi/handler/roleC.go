package handler

import (
	"github.com/gin-gonic/gin"
	"joe-micro/adminApi/model"
	"joe-micro/lib/log"
	"joe-micro/lib/orm"
)

type RoleController struct{} //角色管理控制器

// 获取角色列表
// @Summary 获取角色列表
// @Tags   role   角色管理
// @Accept  json
// @Produce  json
// @Param   page         query    int       false     "页码,默认为1"
// @Param   num          query    int       false     "返回条数,默认为10"
// @Param   sort         query    string    false     "排序字段,默认为createdAt"
// @Param   key          query    string    false     "搜索关键字"
// @Param   orderType    query    string    false     "排序规则,默认为DESC"
// @Param   beginAt      query    string    false     "开始时间"
// @Param   endAt        query    string    false     "结束时间"
// @Success 200 {array}   model.Role 	           "角色列表"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /user_mana/list [get]
func (RoleController) List(c *gin.Context) {
	reqData := ListReq{}
	err := reqData.getListQuery(c)
	if err != nil {
		log.Warn(err)
		resBadRequest(c, err.Error())
		return
	}
	parent_id := c.Query("parent_id")
	var whereOrder []orm.PageWhere
	if reqData.Key != "" {
		v := "%" + reqData.Key + "%"
		var arr []interface{}
		arr = append(arr, v)
		arr = append(arr, v)
		whereOrder = append(whereOrder, orm.PageWhere{Where: "username like ? or real_name like ?", Value: arr})
	}
	if reqData.BeginAt != "" {
		var arr []interface{}
		arr = append(arr, reqData.BeginAt)
		whereOrder = append(whereOrder, orm.PageWhere{Where: "createAt > ? ", Value: arr})
	}
	if reqData.EndAt != "" {
		var arr []interface{}
		arr = append(arr, reqData.EndAt)
		whereOrder = append(whereOrder, orm.PageWhere{Where: "createAt < ? ", Value: arr})
	}
	if parent_id != "" {
		var arr []interface{}
		arr = append(arr, parent_id)
		whereOrder = append(whereOrder, orm.PageWhere{Where: "parent_id = ?", Value: arr})
	}
	list := make([]model.Role, 0)
	var indexPage orm.IndexPage
	indexPage.Page = reqData.Page
	indexPage.Num = reqData.Num
	err = orm.GetPage(&model.Role{}, &model.Role{}, &list, &indexPage, reqData.Sort, whereOrder...)
	if err != nil {
		log.Error(err)
		resErrSrv(c)
		return
	}
	resSuccessPage(c, indexPage, list)
}

// 获取角色详情
// @Summary 获取角色详情
// @Tags   role   角色管理
// @Accept  json
// @Produce  json
// @Param   id         query    string    true     "角色id"
// @Success 200 {object}   model.AdminUser 	    "{code:0,msg:ok}"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /user_mana/detail [get]
func (RoleController) Detail(c *gin.Context) {
	id, exist := c.GetQuery("id")
	if !exist {
		resBadRequest(c, "uid不能为空")
		return
	}
	var role model.Role
	where := model.Role{}
	where.ID = id
	_, err := orm.First(&where, &role)
	if err != nil {
		log.Error(err)
		resErrSrv(c)
		return
	}
	resSuccess(c, role)
}

type updateRoleReq struct {
	ID       string `gorm:"column:id" json:"id"  binding:"required"` //用户uid
	Memo     string `gorm:"column:memo;" json:"memo" `               // 备注
	Name     string `gorm:"column:name;" json:"name" `               // 名称
	Sequence uint64 `gorm:"column:sequence;" json:"sequence" `       // 排序值
}

// 更新角色信息
// @Summary 更新角色信息
// @Tags   role   角色管理
// @Accept  json
// @Produce  json
// @Param   body        body   handler.updateRoleReq     true     "角色信息"
// @Success 200 {object}   handler.ResponseModel 	"角色详情"
// @Failure 400  {object} handler.ResponseModel "{code:1,msg:无效的请求参数}"
// @Failure 500 {object} handler.ResponseModel  "{code:-1,msg:服务器故障}"
// @Security MustToken
// @Router /user_mana/update [post]
func (RoleController) Update(c *gin.Context) {
	reqData := updateRoleReq{}
	err := c.ShouldBind(&reqData)
	if err != nil {
		resBadRequest(c, err.Error())
		return
	}
	where := model.Role{}
	where.ID = reqData.ID
	err = orm.Updates(&where, &reqData)
	if err != nil {
		log.Error(err)
		resErrSrv(c)
		return
	}
	resSuccessMsg(c)
}


// 删除角色

func (RoleController) Delete(c *gin.Context){

}