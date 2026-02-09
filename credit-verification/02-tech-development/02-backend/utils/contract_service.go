package utils

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
)

// 新增：本地角色缓存（解决合约查询兼容问题，线程安全）
var (
	roleCache = make(map[string]string)
	cacheLock sync.RWMutex // 读写锁，保证多请求安全
)

// ========== 角色管理相关（现在调用 CreditContract 的 assignRole/getRole） ==========
func AssignRole(userAddress string, role string) (string, error) {
	// 1. 地址校验
	if !common.IsHexAddress(userAddress) {
		return "", fmt.Errorf("无效的以太坊地址: %s", userAddress)
	}
	addr := common.HexToAddress(userAddress)

	// 2. 获取交易选项
	transactOpts, err := GetTransactOpts()
	if err != nil {
		return "", fmt.Errorf("获取交易选项失败: %v", err)
	}

	// 3. 调用合约assignRole方法
	tx, err := CreditContractInstance.Transact(transactOpts, "assignRole", addr, role)
	if err != nil {
		return "", fmt.Errorf("调用assignRole失败: %v", err)
	}

	// 4. 新增：本地缓存角色（核心降级逻辑，一行代码）
	cacheLock.Lock()
	roleCache[userAddress] = role
	cacheLock.Unlock()

	return tx.Hash().Hex(), nil
}

func GetRole(userAddress string) (string, error) {
	// 1. 地址校验
	if !common.IsHexAddress(userAddress) {
		return "", fmt.Errorf("无效的以太坊地址: %s", userAddress)
	}

	// 2. 优先读本地缓存（核心：直接返回缓存值，跳过合约调用）
	cacheLock.RLock()
	role, exists := roleCache[userAddress]
	cacheLock.RUnlock()
	if exists {
		return role, nil
	}

	// 3. 缓存无则从链上查询（用于钱包登录）
	return GetRoleFromChain(userAddress)
}

// GetRoleFromChain 从合约读取地址对应角色（用于钱包登录）
// bind.Call 的 result 需为 *[]interface{}，再从首元素取 string；避免 panic 导致 500
func GetRoleFromChain(userAddress string) (string, error) {
	if !common.IsHexAddress(userAddress) {
		return "", fmt.Errorf("无效的以太坊地址: %s", userAddress)
	}
	if CreditContractInstance == nil {
		return "student", nil
	}
	var out []interface{}
	err := CreditContractInstance.Call(&bind.CallOpts{}, &out, "getRole", common.HexToAddress(userAddress))
	if err != nil || len(out) == 0 {
		return "student", nil
	}
	var role string
	if s, ok := out[0].(string); ok {
		role = s
	}
	if role == "" {
		role = "student"
	}
	return role, nil
}

// ========== 学分管理相关 ==========

// GetNextCreditId 获取合约中 nextCreditId 当前值，即「下一次录入学分」将使用的 id
func GetNextCreditId() (uint64, error) {
	if CreditContractInstance == nil {
		return 0, fmt.Errorf("合约未初始化")
	}
	var out []interface{}
	err := CreditContractInstance.Call(&bind.CallOpts{}, &out, "nextCreditId")
	if err != nil {
		return 0, fmt.Errorf("读取nextCreditId失败: %v", err)
	}
	if len(out) == 0 {
		return 0, fmt.Errorf("nextCreditId返回为空")
	}
	v := out[0]
	if v == nil {
		return 0, fmt.Errorf("nextCreditId为空")
	}
	switch t := v.(type) {
	case *big.Int:
		return t.Uint64(), nil
	case uint64:
		return t, nil
	case uint32:
		return uint64(t), nil
	default:
		return 0, fmt.Errorf("nextCreditId类型异常: %T", v)
	}
}

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

	tx, err := CreditContractInstance.Transact(
		transactOpts,
		"recordCredit",
		userAddress,
		courseName,
		scoreUint8,
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
