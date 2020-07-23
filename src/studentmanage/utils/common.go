package utils

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

/**
 * @Author: WuNaiChi
 * @Date: 2020/6/5 15:00
 * @Desc:
 */

func WriteInfoToChain(stub shim.ChaincodeStubInterface, id string, data interface{}) error {
	prodBytes, err := json.Marshal(data)
	if err != nil {
		//logs.Error("FailWritrToChain:err%v", err)
		return err
	}

	err = stub.PutState(id, prodBytes)
	if err != nil {
		//logs.Error("FailPutState:err%v", err)
		return err
	}
	return nil
}

func GetCreateStudentParam(s string) (StudentInfo, error) {
	res := StudentInfo{}
	err := json.Unmarshal([]byte(s), &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func Error(msg string) peer.Response {
	return peer.Response{
		Status:  ERROR,
		Message: msg,
	}
}
func SUCCESS(payload []byte) peer.Response {
	return peer.Response{
		Status:  OK,
		Payload: payload,
	}
}

// 获取查询学生列表的参数
func GetQueryStudentParam(s string) (QueryStudentList, error) {
	res := QueryStudentList{}
	err := json.Unmarshal([]byte(s), &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

// 获取查询学生详情的参数
func GetQueryStudentInfoParam(s string) (QueryStudentInfo, error) {
	res := QueryStudentInfo{}
	err := json.Unmarshal([]byte(s), &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func GetUpdateStudentParam(s string) (StudentInfo, error) {
	res := StudentInfo{}
	err := json.Unmarshal([]byte(s), &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func GetDeleteStudentInfoParam(s string) (DeleteStudentInfo, error) {
	res := DeleteStudentInfo{}
	err := json.Unmarshal([]byte(s), &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

// 设置查询语句
func GetQueryStuListSelectString(stuInfo QueryStudentList, docType string) (string, error) {
	selector := SelectorStructInterface{
		Selector: map[string]interface{}{"docType": docType},
	}

	// 根据学号查询
	if stuInfo.AcctId != "" {
		selector.Selector["studentInfo.acctId"] = stuInfo.AcctId
	}
	// 根据姓名（支持模糊查询）
	if stuInfo.Name != "" {
		selector.Selector["studentInfo.name"] = map[string]string{
			"$regex": fmt.Sprintf("%s.*", stuInfo.Name),
		}
	}
	// 性别
	if stuInfo.Sex != "" {
		selector.Selector["studentInfo.sex"] = stuInfo.Sex
	}
	// 年级
	if stuInfo.Grade != "" {
		selector.Selector["studentInfo.grade"] = stuInfo.Grade
	}
	// 爱好（支持模糊查询）
	if stuInfo.Hobby != "" {
		selector.Selector["studentInfo.hobby"] = map[string]string{
			"$regex": fmt.Sprintf(".*%s.*", stuInfo.Hobby),
		}
	}
	selector.Sort = append(selector.Sort, map[string]string{"studentInfo.grade": "desc"})

	selectorBytes, err := json.Marshal(selector)
	if err != nil {
		return "", err
	}
	return string(selectorBytes), nil
}

func GetCounchdbIter(bookmark, selectorString string, pageNo, pageSize int, stub shim.ChaincodeStubInterface) (shim.StateQueryIteratorInterface, *peer.QueryResponseMetadata, error) {
	var stateQueryIteration shim.StateQueryIteratorInterface
	var queryResponseMetadata *peer.QueryResponseMetadata
	var err error
	if bookmark != "" {
		stateQueryIteration, queryResponseMetadata, err = stub.GetQueryResultWithPagination(selectorString, int32(pageSize), bookmark)
		if err != nil {
			fmt.Println("GetQueryResultWithPagination1!!!", err)
			return stateQueryIteration, queryResponseMetadata, err
		}
		// TODO：这里为什么是元数据为空的时候返回,不为空的时候啥都不做吗
		if queryResponseMetadata == nil {
			return stateQueryIteration, queryResponseMetadata, nil
		}
	} else {
		fmt.Println("GetQueryResultWithPagination1OK!!!!")
		bookmark := ""
		for i := 1; i <= pageNo; i++ {
			stateQueryIteration, queryResponseMetadata, err = stub.GetQueryResultWithPagination(selectorString, int32(pageSize), bookmark)
			if err != nil {
				return stateQueryIteration, queryResponseMetadata, err
			}
			if queryResponseMetadata == nil {
				return stateQueryIteration, queryResponseMetadata, nil
			}
			bookmark = queryResponseMetadata.Bookmark
		}
	}
	fmt.Println("GetQueryResultWithPagination02OK!!!!")
	// TODO:这个的作用是什么
	_, respMeta, err := stub.GetQueryResultWithPagination(selectorString, -1, "")
	if err != nil {
		fmt.Println("GetQueryResultWithPagination3!!!!", err)
		return stateQueryIteration, respMeta, err
	}
	fmt.Println("GetQueryResultWithPagination3OK!!!!", err)
	if respMeta == nil {
		return stateQueryIteration, respMeta, nil
	}
	//设置该表有多少条记录
	fmt.Println("okkk!!")
	queryResponseMetadata.FetchedRecordsCount = respMeta.FetchedRecordsCount
	return stateQueryIteration, queryResponseMetadata, nil
}

// K-V查询
func IfExist(stub shim.ChaincodeStubInterface, AcctId string) ([]byte, bool, error) {
	var stateByte []byte
	var err error
	stateByte, err = stub.GetState(AcctId)
	if err != nil {
		return stateByte, false, err
	}
	if stateByte == nil {
		return stateByte, false, nil
	}
	return stateByte, true, nil
}
