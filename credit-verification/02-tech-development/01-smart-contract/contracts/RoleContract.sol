// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

// 权限管理合约：定义不同角色的权限
contract RoleContract is AccessControl,Ownable{
    // 定义角色标识符
    bytes32 public constant STUDENT_ROLE = keccak256("STUDENT_ROLE");
    bytes32 public constant TEACHER_ROLE = keccak256("TEACHER_ROLE");
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");

    // 构造函数：部署者默认是超级管理员（Owner），并赋予ADMIN_ROLE
    constructor() {
        _grantRole(ADMIN_ROLE, msg.sender);
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);// 默认管理员角色，用于分配其他角色
    }

    // 存储用户角色：address → role（admin/teacher/student）
    mapping(address => string) private _userRoles;

    // 分配角色的方法（重点：方法名assignRole，确保存在）
    function assignRole(address user, string calldata role) external onlyOwner {
        require(bytes(role).length > 0, "Role cannot be empty");
        _userRoles[user] = role;

        // 关键：根据字符串角色，分配对应的AccessControl角色
        if (keccak256(bytes(role)) == keccak256(bytes("teacher"))) {
            grantRole(TEACHER_ROLE, user);
        } else if (keccak256(bytes(role)) == keccak256(bytes("admin"))) {
            grantRole(ADMIN_ROLE, user);
        } else if (keccak256(bytes(role)) == keccak256(bytes("student"))) {
            grantRole(STUDENT_ROLE, user);
        }
    }

    // 查询角色的方法
    function getRole(address user) external view returns (string memory) {
        return _userRoles[user];
    }

    // 超级管理员（Owner）分配教师角色
    function grantTeacherRole(address account) external onlyOwner {
        grantRole(TEACHER_ROLE, account);
    }

    // 超级管理员分配管理员角色
    function grantAdminRole(address account) external onlyOwner {
        grantRole(ADMIN_ROLE, account);
    }

    // 超级管理员分配学生角色
    function grantStudentRole(address account) external onlyOwner {
        grantRole(STUDENT_ROLE, account);
    }

    // 撤销角色（仅超级管理员课操作）
    function revokeAnyRole(bytes32 role, address account) external onlyOwner {
        revokeRole(role, account);
    }

    // 校验角色的修饰器
    modifier onlyTeacher() {
        require(hasRole(TEACHER_ROLE, msg.sender),"RoleContract: not a teacher");
        _;
    }

    modifier onlyAdmin() {
        require(hasRole(ADMIN_ROLE, msg.sender),"RoleContract: not a admin");
        _;
    }

    modifier onlyStudent() {
        require(hasRole(STUDENT_ROLE, msg.sender),"RoleContract: not a student");
        _;
    }
    
}