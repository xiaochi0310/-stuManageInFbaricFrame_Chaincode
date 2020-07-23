package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"studentmanage/utils"
)

/**
 * @Author: WuNaiChi
 * @Date: 2020/6/4 14:13
 * @Desc:
 */
type STUChaincode struct{}

// todo:META-INF json文件的内容怎么写

// 创建学生信息
func (t *STUChaincode) createStudentsInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	// 获取参数
	//logs.Info("StartCreate!!")
	res, err := utils.GetCreateStudentParam(args[0])
	if err != nil {
		return utils.Error(err.Error())
	}
	// 查询链上有没有
	_, exit, err := utils.IfExist(stub, res.AcctId)
	if err != nil {
		return utils.Error(err.Error())
	}
	// 存在就不创建
	if exit {
		return utils.Error("don't have the student info")
	}
	// 没有则创建
	chainInStudentInfo := utils.ChainOfStudentInfo{
		DocType:     utils.StudentDocType,
		StudentInfo: res,
	}
	// 信息写到链上
	err = utils.WriteInfoToChain(stub, res.AcctId, chainInStudentInfo)
	if err != nil {
		//logs.Error("failed to create")
		return utils.Error(err.Error())
	}
	return utils.SUCCESS(nil)
}

// 更新学生信息
func (t *STUChaincode) updateStudentInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	//logs.Info("StartUpdate!!")
	rsp, err := utils.GetUpdateStudentParam(args[0])
	if err != nil {
		return utils.Error(err.Error())
	}
	// 查询链上有没有
	stateByte, exit, err := utils.IfExist(stub, rsp.AcctId)
	if err != nil {
		return utils.Error(err.Error())
	}
	if !exit {
		return utils.Error("don't have the student info")
	}
	studentInfo := new(utils.ChainOfStudentInfo)
	err = json.Unmarshal(stateByte, studentInfo)
	if err != nil {
		return utils.Error(err.Error())
	}
	studentInfo.StudentInfo.Name = rsp.Name
	studentInfo.StudentInfo.Grade = rsp.Grade
	studentInfo.StudentInfo.Hobby = rsp.Hobby

	err = utils.WriteInfoToChain(stub, rsp.AcctId, studentInfo)
	if err != nil {
		//logs.Info("FailedStartUpdate!!")
		return utils.Error(err.Error())
	}
	return utils.SUCCESS([]byte{})
}

// 查询学生信息列表
func (t *STUChaincode) queryStudentList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	// 获取参数
	//logs.Info("StartQuery!!")
	fmt.Println("StartQuery!!!")
	rsp, err := utils.GetQueryStudentParam(args[0])
	if err != nil {
		fmt.Println("GetQueryStudentParam!!!", err)
		return utils.Error(err.Error())
	}
	// 参数校验sdk来做
	// 配置select语句
	fmt.Println("GetQueryStudentParam!!!")
	selectString, err := utils.GetQueryStuListSelectString(rsp, utils.StudentDocType)
	if err != nil {
		fmt.Println("GetQueryStuListSelectString!!!", err)
		return utils.Error(err.Error())
	}
	fmt.Println("GetQueryStuListSelectString!!!", selectString)
	StateQueryIterator, QueryResponseMetadata, err := utils.GetCounchdbIter(rsp.Bookmark, selectString, rsp.PageNo, rsp.PageSize, stub)
	if err != nil {
		fmt.Println("GetCounchdbIter!!!", err)
		return utils.Error(err.Error())
	}
	fmt.Println("GetCounchdbIter!!!")
	defer StateQueryIterator.Close() // todo:这个是为了什么
	// todo : 这里为什么是成功
	if QueryResponseMetadata == nil {
		fmt.Println("QueryResponseMetadata is nil!!!")
		return utils.SUCCESS([]byte{})
	}
	fmt.Println("QueryResponseMetadata!!!")
	var res = utils.RspQueryStudentInfo{
		Bookmark: QueryResponseMetadata.Bookmark,
		Count:    QueryResponseMetadata.FetchedRecordsCount,
	}
	// 获取迭代器的指针
	for StateQueryIterator.HasNext() {
		next, err := StateQueryIterator.Next()
		if err != nil {
			fmt.Println("StateQueryIterator!!!", err)
			return utils.Error(err.Error())
		}
		fmt.Println("StateQueryIterator!!!")
		studentInfo := utils.ChainOfStudentInfo{}
		err = json.Unmarshal(next.Value, &studentInfo)
		if err != nil {
			fmt.Println("Unmarshal!!!", err)
			return utils.Error(err.Error())
		}
		res.StudentInfo = append(res.StudentInfo, studentInfo.StudentInfo)
	}
	fmt.Println("Unmarshal!!!")
	if len(res.StudentInfo) == 0 {
		return utils.SUCCESS([]byte{})
	}

	repStudentInfo, err := json.Marshal(res)
	if err != nil {
		//logs.Info("FailedQuery!!")
		fmt.Println("Marshal!!!", err)
		return utils.Error(err.Error())
	}
	fmt.Println("Marshal!!!")
	return utils.SUCCESS(repStudentInfo)
}

// 查询学生详情
func (t *STUChaincode) queryStudentInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	// 获取参数
	//logs.Info("StartQueryInfo!!")
	fmt.Println("StartQueryInfo!!!")
	rsp, err := utils.GetQueryStudentInfoParam(args[0])
	if err != nil {
		fmt.Println("GetQueryStudentInfoParam!!!", err)
		return utils.Error(err.Error())
	}
	// 参数校验sdk来做
	// 使用k-v查询
	stateByte, _, err := utils.IfExist(stub, rsp.AcctId)
	if err != nil {
		fmt.Println("IfExist!!!", err)
		return utils.Error(err.Error())
	}

	if stateByte == nil {
		return utils.SUCCESS([]byte{})
	}

	var studentInfo = utils.ChainOfStudentInfo{}
	err = json.Unmarshal(stateByte, &studentInfo)
	if err != nil {
		fmt.Println("Unmarshal!!!", err)
		return utils.Error(err.Error())
	}

	tmp := utils.StudentInfo{
		AcctId: studentInfo.StudentInfo.AcctId,
		Name:   studentInfo.StudentInfo.Name,
		Sex:    studentInfo.StudentInfo.Sex,
		Grade:  studentInfo.StudentInfo.Grade,
		Hobby:  studentInfo.StudentInfo.Hobby,
	}

	repStudentInfo, err := json.Marshal(tmp)
	if err != nil {
		//logs.Info("FailedQueryInfo!!")
		fmt.Println("Marshal!!!", err)
		return utils.Error(err.Error())
	}
	return utils.SUCCESS(repStudentInfo)
}

// 删除学生信息
func (t *STUChaincode) deleteStudentInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	// 获取参数
	//logs.Info("StartDeleteInfo!!")
	rsp, err := utils.GetDeleteStudentInfoParam(args[0])
	if err != nil {
		return utils.Error(err.Error())
	}
	stateByte, exist, err := utils.IfExist(stub, rsp.AcctId)
	if err != nil || !exist {
		return utils.Error(err.Error())
	}
	if stateByte == nil {
		return utils.SUCCESS([]byte{})
	}
	err = stub.DelState(rsp.AcctId)
	if err != nil {
		//logs.Info("FailedDeleteInfo!!")
		return utils.Error(err.Error())
	}
	return utils.SUCCESS([]byte{})
}

func (t *STUChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	role := &utils.AdminRole{
		RoleId:    utils.RoleAdmin,
		Name:      "管理员",
		Authority: utils.AdminAuthority,
	}
	err := utils.WriteInfoToChain(stub, role.RoleId, role)
	if err != nil {
		return utils.Error(err.Error())
	}

	return utils.SUCCESS(nil)
}
func (t *STUChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	//panic恢复
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	fmt.Println("Invoke!!!!!")
	//检查调用的函数是否存在
	function, args := stub.GetFunctionAndParameters()

	fmt.Println("func, args: ", function, args)
	//校验入参长度
	if len(args) != 1 {
		fmt.Println("func, args: ", args, len(args))
		return utils.Error(utils.InputParaError)
	}

	switch function {
	//角色
	case utils.AuthorityFuncMap[utils.AuthorityCreateMember]: // 上链学生信息（Admin）
		return t.createStudentsInfo(stub, args)
	case utils.AuthorityFuncMap[utils.AuthorityUpdateMember]: // 更新学生信息（Admin）
		return t.updateStudentInfo(stub, args)
	case utils.AuthorityFuncMap[utils.AuthorityGetMemberList]: // 查新学生信息（Admin）
		return t.queryStudentList(stub, args)
	case utils.AuthorityFuncMap[utils.AuthorityGetMemberInfo]: // 查询学生信息列表（Admin）
		return t.queryStudentInfo(stub, args)
	case utils.AuthorityFuncMap[utils.AuthorityDeleteMember]: // 删除学生信息（Admin）
		return t.deleteStudentInfo(stub, args)

	default:
		return utils.Error(utils.InternalError)
	}

}

func main() {
	err := shim.Start(new(STUChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
