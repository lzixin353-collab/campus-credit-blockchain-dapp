const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("CreditContract", function () {
    let Creditcontract;
    let CreaditContract;
    let owner;// 超级管理员
    let teacher;// 教师
    let adimt;// 管理员
    let student;// 学生
    let randomUser;// 普通用户

    // 每个测试前部署合约
    beforeEach(async function (){
        CreditContract = await ethers.getContractFactory("CreditContract");
        [owner, teacher, admin, student, randomUser] = await ethers.getSigners();

        creditContract = await CreditContract.deploy();
        await creditContract.deployed();

        // 分配角色：给teacher赋予教师角色，admin赋予管理员角色，student赋予学生角色
        await creditContract.grantTeacherRole(teacher.address);
        await creditContract.grantAdminRole(admin.address);
        await creditContract.grantStudentRole(student.address);
    });

    // 测试1：角色分配功能
  it("Should assign roles correctly", async function () {
    expect(await creditContract.hasRole(await creditContract.TEACHER_ROLE(), teacher.address)).to.be.true;
    expect(await creditContract.hasRole(await creditContract.ADMIN_ROLE(), admin.address)).to.be.true;
    expect(await creditContract.hasRole(await creditContract.STUDENT_ROLE(), student.address)).to.be.true;
    expect(await creditContract.hasRole(await creditContract.TEACHER_ROLE(), randomUser.address)).to.be.false;
  });

  // 测试2：教师录入学分
  it("Should allow teacher to record credit", async function () {
    const studentId = "20230001";
    const courseName = "区块链原理";
    const score = 90;

    // 教师录入学分
    await expect(creditContract.connect(teacher).recordCredit(studentId, courseName, score))
      .to.emit(creditContract, "CreditRecorded")
      .withArgs(0, studentId, teacher.address);

    // 校验学分信息
    const credit = await creditContract.credits(0);
    expect(credit.studentId).to.equal(studentId);
    expect(credit.courseName).to.equal(courseName);
    expect(credit.score).to.equal(score);
    expect(credit.isApproved).to.be.false;
  });

  // 测试3：非教师不能录入学分
  it("Should reject non-teacher to record credit", async function () {
    await expect(creditContract.connect(randomUser).recordCredit("20230001", "高数", 80))
      .to.be.revertedWith("RoleContract: not a teacher");
  });

  // 测试4：管理员审核学分
  it("Should allow admin to approve credit", async function () {
    // 先录入学分
    await creditContract.connect(teacher).recordCredit("20230001", "区块链原理", 90);
    // 管理员审核
    await expect(creditContract.connect(admin).approveCredit(0))
      .to.emit(creditContract, "CreditApproved")
      .withArgs(0, admin.address);

    // 校验审核状态
    const credit = await creditContract.credits(0);
    expect(credit.isApproved).to.be.true;
  });

  // 测试5：查询学生学分
  it("Should return student's credits", async function () {
    const studentId = "20230001";
    // 录入2个学分
    await creditContract.connect(teacher).recordCredit(studentId, "区块链原理", 90);
    await creditContract.connect(teacher).recordCredit(studentId, "Web3开发", 85);

    // 查询学生学分
    const credits = await creditContract.getStudentCredits(studentId);
    expect(credits.length).to.equal(2);
    expect(credits[0].courseName).to.equal("区块链原理");
    expect(credits[1].courseName).to.equal("Web3开发");
  });

});