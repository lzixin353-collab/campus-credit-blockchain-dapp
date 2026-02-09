# 校园区块链学分存证DApp - 轻量后端

## 项目介绍
基于Go+Gin+以太坊（Hardhat）实现的轻量后端，包含用户管理、合约角色分配/查询核心功能。

## 环境要求
1. Go 1.18+
2. Node.js 16+（Hardhat合约部署）
3. MySQL 8.0+
4. Hardhat本地节点

## 快速运行
### 1. 合约部署
```bash
cd 02-tech-development/01-smart-contract
npm install
npx hardhat node  # 启动本地节点
npx hardhat run scripts/deploy.js --network localhost  # 部署合约，记录地址