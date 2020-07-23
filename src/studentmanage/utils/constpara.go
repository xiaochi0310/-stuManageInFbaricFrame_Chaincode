package utils

/**
 * @Author: WuNaiChi
 * @Date: 2020/6/5 14:02
 * @Desc:
 */

type AdminRole struct {
	RoleId    string    `json:"roleId"`
	Name      string    `json:"name"`
	Authority Authority `json:"authority"`
}

const (
	RoleAdmin   = "role-admin"
	RoleStudent = "role-student"
	//管理员的权限集合
	AdminAuthority = AuthorityCreateMember | AuthorityUpdateMember | AuthorityGetMemberList | AuthorityGetMemberInfo | AuthorityDeleteMember

	StudentDocType = "table-student"
)

type SelectorStructInterface struct {
	Selector map[string]interface{} `json:"selector"`
	UseIndex []string               `json:"use_index,omitempty"`
	Sort     []interface{}          `json:"sort,omitempty"`
}

type Authority int

const (
	AuthorityCreateMember  Authority = 1 << iota //1
	AuthorityUpdateMember                        //2
	AuthorityGetMemberList                       //...
	AuthorityGetMemberInfo
	AuthorityDeleteMember
)

var AuthorityFuncMap = map[Authority]string{
	AuthorityCreateMember:  "createStudentInfo",
	AuthorityUpdateMember:  "updateStudentInfo",
	AuthorityGetMemberList: "getStudentInfoList",
	AuthorityGetMemberInfo: "queryStudentInfo",
	AuthorityDeleteMember:  "deleteStudentInfo",
}

const (
	OK = 200

	ERROR    = 400
	ERROR401 = 401
	ERROR403 = 403
	ERROR500 = 500
)

// 响应的错误信息
const (
	InputParaError = "InputParaError"
	InternalError  = "InternalError"
)

// 学生信息
type StudentInfo struct {
	AcctId string `json:"acctId"` // 学号
	Name   string `json:"name"`   // 姓名
	Sex    string `json:"sex"`    // 性别：map
	Grade  string `json:"grade"`  // 年级: map 123
	Hobby  string `json:"hobby"`  //
}

// 学生查询信息
type QueryStudentList struct {
	PageNo   int    `json:"pageNo,omitempty"`   //从第1页开始
	PageSize int    `json:"pageSize,omitempty"` //每页条数
	Bookmark string `json:"bookmark,omitempty"` //couchdb中的书签号
	AcctId   string `json:"acctId"`
	Name     string `json:"name"` // 模糊查询
	Sex      string `json:"sex"`
	Grade    string `json:"grade"`
	Hobby    string `json:"hobby"` // 模糊查询
}
type ChainOfStudentInfo struct {
	DocType     string      `json:"docType"`
	StudentInfo StudentInfo `json:"studentInfo"`
}

type QueryStudentInfo struct {
	AcctId string `json:"acctId"`
}

type RspQueryStudentInfo struct {
	StudentInfo []StudentInfo `json:"studentInfo"`
	Bookmark    string        `json:"bookmark,omitempty"` //couchdb中的书签号
	Count       int32         `json:"count"`
}

type DeleteStudentInfo struct {
	AcctId string `json:"acctId"`
}
