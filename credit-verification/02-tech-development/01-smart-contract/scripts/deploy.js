const hre = require("hardhat");

async function main() {
  console.log("开始部署极简 CreditContract（无继承）...");

  // 仅部署 CreditContract（内置所有权限逻辑）
  const CreditContract = await hre.ethers.getContractFactory("CreditContract");
  const creditContract = await CreditContract.deploy(); 
  await creditContract.deployed();
  console.log(`✅ CreditContract 部署完成，地址: ${creditContract.address}`);

  // 部署总结
  console.log("\n📌 本地部署总结：");
  console.log(`- CreditContract 地址: ${creditContract.address}`);
  console.log(`- 本地 RPC 地址: http://127.0.0.1:8545`);
  console.log(`- 部署者地址（默认教师/管理员）: ${(await hre.ethers.getSigners())[0].address}`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("❌ 部署失败：", error);
    process.exit(1);
  });