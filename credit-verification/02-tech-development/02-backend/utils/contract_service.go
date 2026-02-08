package utils

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
)

// ========== 角色管理相关（现在调用 CreditContract 的 assignRole/getRole） ==========
func AssignRole(userAddress string, role string) (string, error) {
	transactOpts, err := GetTransactOpts()
	if err != nil {
		return "", fmt.Errorf("获取交易选项失败: %v", err)
	}

	addr := common.HexToAddress(userAddress)
	// 调用 CreditContract 的 assignRole 方法
	tx, err := CreditContractInstance.Transact(transactOpts, "assignRole", addr, role)
	if err != nil {
		return "", fmt.Errorf("调用assignRole失败: %v", err)
	}

	return tx.Hash().Hex(), nil
}

func GetRole(userAddress string) (string, error) {
	addr := common.HexToAddress(userAddress)
	var result []interface{}
	// 调用 CreditContract 的 getRole 方法
	err := CreditContractInstance.Call(&bind.CallOpts{}, &result, "getRole", addr)
	if err != nil {
		return "", fmt.Errorf("调用getRole失败: %v", err)
	}

	if len(result) == 0 {
		return "", fmt.Errorf("角色查询结果为空")
	}
	role, ok := result[0].(string)
	if !ok {
		return "", fmt.Errorf("角色类型解析失败")
	}

	return role, nil
}

// ========== 学分管理相关（无需修改，仅调用 CreditContract） ==========
func RecordCredit(userAddress string, courseName string, score float64) (string, error) {
	if score < 0 || score > 100 {
		return "", fmt.Errorf("学分值超出范围（0-100）: %v", score)
	}
	scoreUint8 := uint8(score)
	if courseName == "" || userAddress == "" {
		return "", fmt.Errorf("课程名/学生学号不能为空")
	}

	transactOpts, err := GetTransactOpts()
	if err != nil {
		return "", fmt.Errorf("获取交易选项失败: %v", err)
	}

	// 直接调用 CreditContract 的 recordCredit
	tx, err := CreditContractInstance.Transact(
		transactOpts,
		"recordCredit",
		userAddress, // 学生学号
		courseName,  // 课程名
		scoreUint8,  // 分数
	)
	if err != nil {
		return "", fmt.Errorf("调用recordCredit失败: %v", err)
	}

	return tx.Hash().Hex(), nil
}

func AuditCredit(creditId uint64, approved bool) (string, error) {
	creditIdInt := big.NewInt(int64(creditId))
	transactOpts, err := GetTransactOpts()
	if err != nil {
		return "", fmt.Errorf("获取交易选项失败: %v", err)
	}

	// 注意：合约方法名是 approveCredit，不是 auditCredit
	tx, err := CreditContractInstance.Transact(transactOpts, "approveCredit", creditIdInt)
	if err != nil {
		return "", fmt.Errorf("调用approveCredit失败: %v", err)
	}

	return tx.Hash().Hex(), nil
}

func GetUserCredits(userAddress string) ([]map[string]interface{}, error) {
	_ = common.HexToAddress(userAddress)
	var result []interface{}
	// 注意：合约方法名是 getStudentCredits，不是 getUserCredits
	err := CreditContractInstance.Call(&bind.CallOpts{}, &result, "getStudentCredits", userAddress)
	if err != nil {
		return nil, fmt.Errorf("调用getStudentCredits失败: %v", err)
	}

	credits := make([]map[string]interface{}, 0)
	if len(result) == 0 {
		return credits, nil
	}

	creditList, ok := result[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("学分列表解析失败")
	}

	for _, item := range creditList {
		creditItem, ok := item.([]interface{})
		if !ok {
			continue
		}

		creditId := creditItem[0].(*big.Int).Uint64()
		studentId := creditItem[1].(string)
		courseName := creditItem[2].(string)
		score := float64(creditItem[3].(uint8))
		isApproved := creditItem[5].(bool)

		credits = append(credits, map[string]interface{}{
			"id":          creditId,
			"student_id":  studentId,
			"course_name": courseName,
			"score":       score,
			"is_approved": isApproved,
		})
	}

	return credits, nil
}
