// model/credit_db.go 学分表 CRUD（与 init.sql 中 credits 表结构一致，使用 sql.DB）
package model

import (
	"database/sql"
	"time"

	"campus-credit-backend/utils"
)

// CreditRow 学分表行（与 credits 表一一对应）
type CreditRow struct {
	Id               int64          `json:"id"`
	ContractCreditId sql.NullInt64  `json:"contract_credit_id"`
	StudentAddress   string         `json:"student_address"`
	TeacherAddress   string         `json:"teacher_address"`
	CourseName       string         `json:"course_name"`
	Score            float64        `json:"score"`
	Status           string         `json:"status"` // pending / approved / rejected
	TxHash           sql.NullString `json:"tx_hash"`
	AuditAdmin       sql.NullString `json:"audit_admin"`
	AuditTime        sql.NullTime   `json:"audit_time"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// CreateCredit 插入一条学分记录（录入学分后调用）
func CreateCredit(studentAddress, teacherAddress, courseName string, score float64, status, txHash string, contractCreditId int64) (int64, error) {
	res, err := utils.DB.Exec(
		`INSERT INTO credits (contract_credit_id, student_address, teacher_address, course_name, score, status, tx_hash) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		contractCreditId, studentAddress, teacherAddress, courseName, score, status, txHash,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetCreditsByStudentAddress 按学生地址查询学分列表
func GetCreditsByStudentAddress(studentAddress string) ([]CreditRow, error) {
	rows, err := utils.DB.Query(
		`SELECT id, contract_credit_id, student_address, teacher_address, course_name, score, status, tx_hash, audit_admin, audit_time, created_at, updated_at 
		 FROM credits WHERE student_address = ? ORDER BY created_at DESC`,
		studentAddress,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCreditRows(rows)
}

// GetCreditsByTeacherAddress 按教师地址查询其录入的学分列表
func GetCreditsByTeacherAddress(teacherAddress string) ([]CreditRow, error) {
	rows, err := utils.DB.Query(
		`SELECT id, contract_credit_id, student_address, teacher_address, course_name, score, status, tx_hash, audit_admin, audit_time, created_at, updated_at 
		 FROM credits WHERE teacher_address = ? ORDER BY created_at DESC`,
		teacherAddress,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCreditRows(rows)
}

// GetAllCredits 管理员：查询全部学分
func GetAllCredits() ([]CreditRow, error) {
	rows, err := utils.DB.Query(
		`SELECT id, contract_credit_id, student_address, teacher_address, course_name, score, status, tx_hash, audit_admin, audit_time, created_at, updated_at 
		 FROM credits ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCreditRows(rows)
}

// GetPendingCredits 待审核学分列表（管理员用，仅含已有关链上ID的记录）
func GetPendingCredits() ([]CreditRow, error) {
	rows, err := utils.DB.Query(
		`SELECT id, contract_credit_id, student_address, teacher_address, course_name, score, status, tx_hash, audit_admin, audit_time, created_at, updated_at 
		 FROM credits WHERE status = 'pending' AND contract_credit_id > 0 ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCreditRows(rows)
}

// UpdateCreditStatus 更新审核状态
func UpdateCreditStatus(id int64, status, auditAdmin string) error {
	_, err := utils.DB.Exec(
		`UPDATE credits SET status = ?, audit_admin = ?, audit_time = NOW() WHERE id = ?`,
		status, auditAdmin, id,
	)
	return err
}

// GetCreditById 按主键查一条
func GetCreditById(id int64) (*CreditRow, error) {
	var row CreditRow
	err := utils.DB.QueryRow(
		`SELECT id, contract_credit_id, student_address, teacher_address, course_name, score, status, tx_hash, audit_admin, audit_time, created_at, updated_at 
		 FROM credits WHERE id = ?`,
		id,
	).Scan(
		&row.Id, &row.ContractCreditId, &row.StudentAddress, &row.TeacherAddress, &row.CourseName, &row.Score,
		&row.Status, &row.TxHash, &row.AuditAdmin, &row.AuditTime, &row.CreatedAt, &row.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func scanCreditRows(rows *sql.Rows) ([]CreditRow, error) {
	var list []CreditRow
	for rows.Next() {
		var row CreditRow
		err := rows.Scan(
			&row.Id, &row.ContractCreditId, &row.StudentAddress, &row.TeacherAddress, &row.CourseName, &row.Score,
			&row.Status, &row.TxHash, &row.AuditAdmin, &row.AuditTime, &row.CreatedAt, &row.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, row)
	}
	return list, rows.Err()
}
