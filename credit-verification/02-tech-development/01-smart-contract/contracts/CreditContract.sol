// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

// 极简学分合约：内置权限，无继承，无依赖
contract CreditContract {
    // 核心状态变量
    address public owner; // 部署者=超级管理员
    mapping(address => bool) public isTeacher; // 教师权限映射
    mapping(address => bool) public isAdmin; // 管理员权限映射

    // 学分结构体（保持你的原有结构）
    struct Credit {
        uint256 creditId;
        string studentId;
        string courseName;
        uint8 score;
        address teacherAddress;
        bool isApproved;
        bool exists;
    }

    // 学分存储（保持你的原有映射）
    mapping(uint256 => Credit) public credits;
    mapping(string => uint256[]) public studentCreditIds;
    uint256 public nextCreditId;

    // 事件（保持原有）
    event CreditRecorded(
        uint256 indexed creditId, 
        string indexed studentId, 
        string courseName,
        uint8 score,
        address indexed teacherAddress
    );
    event CreditApproved(
        uint256 indexed creditId, 
        address indexed adminAddress
    );
    event RoleAssigned(address indexed user, string indexed role);

    // 构造函数：部署者默认拥有所有权限
    constructor() {
        owner = msg.sender;
        isTeacher[msg.sender] = true;
        isAdmin[msg.sender] = true;
    }

    // 基础权限修饰符（极简，无依赖）
    modifier onlyOwner() {
        require(msg.sender == owner, "CreditContract: only owner");
        _;
    }
    modifier onlyTeacher() {
        require(isTeacher[msg.sender], "CreditContract: not a teacher");
        _;
    }
    modifier onlyAdmin() {
        require(isAdmin[msg.sender], "CreditContract: not a admin");
        _;
    }

    // 分配角色（仅Owner可操作，适配你的原有接口）
    function assignRole(address user, string calldata role) external onlyOwner {
        require(bytes(role).length > 0, "Role cannot be empty");
        if (keccak256(bytes(role)) == keccak256(bytes("teacher"))) {
            isTeacher[user] = true;
            emit RoleAssigned(user, "teacher");
        } else if (keccak256(bytes(role)) == keccak256(bytes("admin"))) {
            isAdmin[user] = true;
            emit RoleAssigned(user, "admin");
        } else if (keccak256(bytes(role)) == keccak256(bytes("student"))) {
            // 学生角色仅标记，无特殊权限
            emit RoleAssigned(user, "student");
        }
    }

    // 查询角色（适配你的原有接口）
    function getRole(address user) external view returns (string memory) {
        if (isTeacher[user]) return "teacher";
        if (isAdmin[user]) return "admin";
        return "";
    }

    // 录入学分（完全保留你的业务逻辑，仅权限修饰符极简）
    function recordCredit(
        string calldata studentId,
        string calldata courseName,
        uint8 score
    ) external onlyTeacher {
        require(score <= 100, "CreditContract: invalid score(0-100)");
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
        studentCreditIds[studentId].push(creditId);

        emit CreditRecorded(creditId, studentId, courseName, score, msg.sender);
    }

    // 审核学分（保留原有逻辑）
    function approveCredit(uint256 creditId) external onlyAdmin {
        require(credits[creditId].exists, "CreditContract: credit not exist");
        require(!credits[creditId].isApproved, "CreditContract: credit already approved");
        credits[creditId].isApproved = true;
        emit CreditApproved(creditId, msg.sender);
    }

    // 查询学生学分（保留原有逻辑）
    function getStudentCredits(string calldata studentId) external view returns (Credit[] memory) {
        require(bytes(studentId).length > 0, "CreditContract: studentId empty");
        uint256[] memory ids = studentCreditIds[studentId];
        uint256 validCount = 0;

        for (uint256 i = 0; i < ids.length; i++) {
            if (credits[ids[i]].exists) validCount++;
        }

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

    // 查询单个学分（保留原有逻辑）
    function getCreditById(uint256 creditId) external view returns (Credit memory) {
        require(credits[creditId].exists, "CreditContract: credit not exist");
        return credits[creditId];
    }
}