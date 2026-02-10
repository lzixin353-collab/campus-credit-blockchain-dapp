const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("CreditContract", function () {
  let creditContract;
  let owner;
  let teacher;
  let admin;
  let student;
  let randomUser;

  beforeEach(async function () {
    const CreditContract = await ethers.getContractFactory("CreditContract");
    [owner, teacher, admin, student, randomUser] = await ethers.getSigners();
    creditContract = await CreditContract.deploy();
    await creditContract.deployed();

    // 当前合约只有 assignRole(string)，用 assignRole 分配角色
    await creditContract.assignRole(teacher.address, "teacher");
    await creditContract.assignRole(admin.address, "admin");
    // 学生无映射，仅事件；查询 getRole 返回 ""
  });

  it("Should assign roles correctly", async function () {
    expect(await creditContract.getRole(teacher.address)).to.equal("teacher");
    expect(await creditContract.getRole(admin.address)).to.equal("admin");
    expect(await creditContract.getRole(student.address)).to.equal("");
    expect(await creditContract.getRole(randomUser.address)).to.equal("");
  });

  it("Should allow teacher to record credit", async function () {
    const studentId = "20230001";
    const courseName = "区块链原理";
    const score = 90;

    await expect(
      creditContract.connect(teacher).recordCredit(studentId, courseName, score)
    ).to.emit(creditContract, "CreditRecorded");

    const credit = await creditContract.credits(0);
    expect(credit.studentId).to.equal(studentId);
    expect(credit.courseName).to.equal(courseName);
    expect(credit.score).to.equal(score);
    expect(credit.isApproved).to.be.false;
  });

  it("Should reject non-teacher to record credit", async function () {
    await expect(
      creditContract.connect(randomUser).recordCredit("20230001", "高数", 80)
    ).to.be.revertedWith("CreditContract: not a teacher");
  });

  it("Should allow admin to approve credit", async function () {
    await creditContract.connect(teacher).recordCredit("20230001", "区块链原理", 90);
    await expect(creditContract.connect(admin).approveCredit(0)).to.emit(
      creditContract,
      "CreditApproved"
    );

    const credit = await creditContract.credits(0);
    expect(credit.isApproved).to.be.true;
  });

  it("Should return student credits by studentId", async function () {
    const studentId = "20230001";
    await creditContract.connect(teacher).recordCredit(studentId, "区块链原理", 90);
    await creditContract.connect(teacher).recordCredit(studentId, "Web3开发", 85);

    const credits = await creditContract.getStudentCredits(studentId);
    expect(credits.length).to.equal(2);
    expect(credits[0].courseName).to.equal("区块链原理");
    expect(credits[1].courseName).to.equal("Web3开发");
  });
});
