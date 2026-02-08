// scripts/deploy.js - 适配 CreditContract 0参数构造函数
const hre = require("hardhat");

async function main() {
  console.log("开始部署合约到本地 Hardhat 节点...");

  // 1. 部署 RoleContract（权限管理合约）
  const RoleContract = await hre.ethers.getContractFactory("RoleContract");
  const roleContract = await RoleContract.deploy();
  await roleContract.deployed(); // v5 等待部署完成
  console.log(`✅ RoleContract 部署完成，地址: ${roleContract.address}`);

  // 2. 部署 CreditContract（学分合约，无构造函数参数）
  const CreditContract = await hre.ethers.getContractFactory("CreditContract");
  // 关键修改：移除 roleContract.address 参数，适配0参数构造函数
  const creditContract = await CreditContract.deploy(); 
  await creditContract.deployed();
  console.log(`✅ CreditContract 部署完成，地址: ${creditContract.address}`);

  // 部署总结
  console.log("\n📌 本地部署总结：");
  console.log(`- RoleContract 地址: ${roleContract.address}`);
  console.log(`- CreditContract 地址: ${creditContract.address}`);
  console.log(`- 本地 RPC 地址: http://127.0.0.1:8545`);
}

// 执行部署并捕获错误
main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("❌ 部署失败：", error);
    process.exit(1);
  });