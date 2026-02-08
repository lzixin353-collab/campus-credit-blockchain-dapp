# Credit Verification DApp
基于以太坊智能合约+Go后端的学分验证DApp

## 项目结构
- 01-smart-contract：智能合约（RoleContract、CreditContract）
- 02-backend：Go后端服务（Gin框架+MySQL+以太坊交互）

## 快速启动
### 1. 智能合约部署
```bash
cd 02-tech-development/01-smart-contract
npm install
npx hardhat node
npx hardhat run scripts/deploy.js --network localhost
