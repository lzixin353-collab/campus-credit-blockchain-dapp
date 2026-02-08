// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

import "./RoleContract.sol";

// 学分存证核心合约
contract CreditContract is RoleContract {
    // 学分结构体：存储核心信息
    struct Credit {
        uint256 creditId; // 学分唯一ID
        string studentId; // 学生学号（链下关联，便于查询）
        string courseName; // 课程名称
        uint8 score; // 学分分数（0-100）
        address teacherAddress; // 录入的教师地址
        bool isApproved; // 是否通过管理员审核
        bool exists; // 新增：用布尔值标记是否存在，替代时间戳
    }

    // 存储所有学分信息（creditId => Credit）
    mapping(uint256 => Credit) public credits;
    // 学生学号 => 学分ID列表（便于学生查询自己的所有学分）
    mapping(string => uint256[]) public studentCreditIds;
    // 自增的学分ID
    uint256 public nextCreditId;

    // 事件：学分录入
    event CreditRecorded(
        uint256 indexed creditId, 
        string indexed studentId, 
        string courseName,
        uint8 score,
        address indexed teacherAddress
    );
    // 事件：学分审核通过
    event CreditApproved(
        uint256 indexed creditId, 
        address indexed adminAddress
    );

    // 教师录入学分（仅教师角色可操作）
    function recordCredit(
        string calldata studentId,
        string calldata courseName,
        uint8 score
    ) external onlyTeacher {
        // 校验分数范围（0-100）
        require(score <= 100, "CreditContract: invalid score(0-100)");
        // 校验学生学号非空（简单校验）
        require(bytes(studentId).length > 0, "CreditContract: empty studentId");
        require(bytes(courseName).length > 0, "CreditContract: courseName empty");

        uint256 creditId = nextCreditId;
        credits[creditId] = Credit({
            creditId: creditId,
            studentId: studentId,
            courseName: courseName,
            score: score,
            teacherAddress: msg.sender,
            isApproved: false,
            exists: true 
        });
        nextCreditId++;

        // 将学分ID关联到学生学号
        studentCreditIds[studentId].push(creditId);

        emit CreditRecorded(creditId, studentId, courseName, score, msg.sender);
    }

    // 管理员审核学分（仅管理员角色可操作）
    function approveCredit(uint256 creditId) external onlyAdmin {
        // 校验学分存在
        require(credits[creditId].exists, "CreditContract: credit not exist");
        // 校验未审核过
        require(!credits[creditId].isApproved, "CreditContract: credit already approved");

        credits[creditId].isApproved = true;
        emit CreditApproved(creditId, msg.sender);
    }

    // 查询学生所有学分（任何人可查，但学生只能看自己的，前端做控制）
    function getStudentCredits(string calldata studentId) external view returns (Credit[] memory) {
        require(bytes(studentId).length > 0, "CreditContract: studentId empty");

        uint256[] memory ids = studentCreditIds[studentId];
        uint256 validCount = 0;

        // 第一步：统计有效学分数量（exists=true）
        for (uint256 i = 0; i < ids.length; i++) {
            if (credits[ids[i]].exists) {
                validCount++;
            }
        }

        // 第二步：填充有效学分数据
        Credit[] memory result = new Credit[](validCount);
        uint256 index = 0;
        for (uint256 i = 0; i < ids.length; i++) {
            if (credits[ids[i]].exists) {
                result[index] = credits[ids[i]];
                index++;
            }
        }

        return result;
    }

    // 查询单个学分详情
    function getCreditById(uint256 creditId) external view returns (Credit memory) {
        require(credits[creditId].exists, "CreditContract: credit not exist");
        return credits[creditId];
    }
}