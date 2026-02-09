# 校园区块链学分存证轻量 DApp

基于区块链的校园学分录入、审核与查询系统，支持学生端、教师端、管理员端三端角色，学分上链存证、链上+后端双重权限校验。

---

## 功能特性

| 端     | 功能说明 |
|--------|----------|
| **学生端** | 钱包登录、我的学分列表、个人信息、审核状态查看 |
| **教师端** | 钱包登录、学分录入（上链+落库）、录入列表、待审核状态 |
| **管理员端** | 钱包登录、学分审核（通过/驳回）、角色管理（链上分配 teacher/admin/student） |

- 登录方式：Metamask 钱包连接 + 后端按地址/数据库角色签发 JWT。
- 学分流程：教师录入学分 → 上链并落库 → 管理员审核 → 学生可查已通过记录。
- 角色与地址：支持 Postman 注册账号 + 绑定钱包；管理员可为地址分配链上角色。

---

## 技术栈

| 模块       | 技术 |
|------------|------|
| 智能合约   | Solidity 0.8.x、OpenZeppelin（RoleContract）、Hardhat |
| 后端       | Go、Gin、MySQL、JWT、go-ethereum |
| 前端       | Vue 3、Vite、Element Plus、Pinia、Vue Router、Axios、Web3.js |
| 产品设计   | 问卷调研、需求拆解、墨刀低保真原型 |

---

## 运行效果

项目已本地跑通，可现场/远程演示；代码与说明见本仓库。

| 页面 | 说明 |
|------|------|
| 登录页 | Metamask 连接后显示「已连接」与短地址，可点击「登录系统」按角色进入 |
| 学生端 | 我的学分、个人信息 |
| 教师端 | 学分录入、录入列表 |
| 管理员端 | 学分审核、角色管理 |
 
将本地运行后的截图放入 `docs/screenshots/`，建议文件名见 [docs/screenshots/README.md](docs/screenshots/README.md)，则下方图片会自动显示。


---

## 项目结构

```
credit-verification/
├── 01-product-design/          # 产品设计
│   ├── 01-research-questionnaire/   # 调研问卷
│   ├── 02-requirement-analysis/     # 需求与角色边界
│   └── 03-mock-prototype/           # 墨刀原型与截图
├── 02-tech-development/         # 技术开发
│   ├── 01-smart-contract/       # 合约（CreditContract + RoleContract）
│   ├── 02-backend/              # Go 后端（Gin + MySQL + 链上交互）
│   └── 03-frontend/             # Vue3 前端
└── README.md
```

---

## 上传 GitHub 前注意

- **不要上传**：`02-backend/config/config.yaml`（含数据库密码、私钥、JWT 密钥），已通过 `.gitignore` 排除；本地请复制 `config.example.yaml` 为 `config.yaml` 后填写。
- **不要上传**：后端编译产物 `campus-credit`、各目录下的 `node_modules`、`.env` 等，已由各层 `.gitignore` 排除。
- **建议上传**：产品设计文档与截图、合约与前后端源码、`init.sql`、`config.example.yaml`、本 README。

---

## 环境要求

- **Node.js** 16+
- **Go** 1.20+
- **MySQL** 8.0+（或 5.7+）
- **Metamask** 浏览器插件
- （可选）**npx hardhat node** 本地链，或 Infura/Alchemy 等 RPC

---

## 本地运行

### 1. 数据库

创建库并执行初始化 SQL：

```bash
mysql -u root -p < 02-tech-development/02-backend/config/init.sql
```

若表已存在，需保证 `users` 表含 `username, password, address, role`，`credits` 表含 `contract_credit_id` 等字段；缺少时可参考 `init.sql` 中建表与注释中的 `ALTER`。

### 2. 智能合约（本地链）

```bash
cd 02-tech-development/01-smart-contract
npm install
npx hardhat node          # 终端一：启动本地链（保持运行）
npx hardhat run scripts/deploy.js --network localhost   # 终端二：部署
```

记下输出的 **CreditContract 地址**，用于后端与前端配置。

### 3. 后端

```bash
cd 02-tech-development/02-backend
# 修改 config/config.yaml：mysql.dsn、ethereum.credit_contract_addr、ethereum.private_key（与 hardhat 账户一致）
go mod tidy
go run main.go
```

默认端口 **8080**。接口前缀：`/api`（如 `/api/user/login`、`/api/credit/record`）。

### 4. 前端

```bash
cd 02-tech-development/03-frontend
npm install
npm run dev
```

默认 **http://localhost:8081**，代理 `/api` 到 `http://localhost:8080`（见 `vite.config.js`）。

### 5. 使用流程简述

1. 浏览器打开 http://localhost:8081 ，点击「连接 Metamask 钱包」并切换到 Hardhat 本地网（链 ID 31337）。
2. **首次使用**：可用 Postman 注册账号（如 admin/teacher），再在「个人信息」或调用 `POST /api/user/bind-address` 绑定当前钱包地址；或由管理员在「角色管理」中为地址分配链上角色。
3. 点击「登录系统」→ 按角色进入学生/教师/管理员首页。
4. **教师**：在「学分录入」填写学生地址/学号、课程名称、成绩(0-100)，提交上链并落库。
5. **管理员**：在「学分审核」对待审核记录进行通过/驳回；在「角色管理」为钱包地址分配角色。
6. **学生**：在「我的学分」查看已录入与审核状态。

---

## 配置说明

### 后端 `02-tech-development/02-backend/config/config.yaml`

- **server.port**：服务端口，默认 8080。
- **mysql.dsn**：数据库连接串，按本机 MySQL 账号密码修改。
- **ethereum.rpc_url**：链 RPC，本地为 `http://127.0.0.1:8545`。
- **ethereum.credit_contract_addr**：部署后的 CreditContract 地址。
- **ethereum.private_key**：后端用于发链上交易的私钥（如 hardhat 默认账户）。
- **jwt.secret / jwt.expire_hours**：登录 Token 配置。

### 前端合约地址

若部署后合约地址变更，需同步修改：

- `02-tech-development/03-frontend/src/utils/web3.js` 中的 `CREDIT_CONTRACT_ADDRESS`（前端直接调合约时使用，当前登录已不依赖前端 getRole）。

---

## 接口一览（后端）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/user/register | 注册（测试用） |
| POST | /api/user/login | 登录（支持 username+password 或 address 钱包登录） |
| GET  | /api/user/info | 当前用户信息（需 Token） |
| POST | /api/user/update | 更新用户信息/地址（需 Token） |
| POST | /api/user/bind-address | 绑定钱包地址到当前账号（需 Token） |
| GET  | /api/credit/list | 学分列表（按角色：学生/教师/管理员） |
| POST | /api/credit/record | 教师录入学分（需 Token） |
| GET  | /api/credit/pending | 管理员待审核列表（需 Token） |
| POST | /api/credit/approve | 管理员审核通过（需 Token） |
| POST | /api/credit/reject | 管理员驳回（需 Token） |
| POST | /api/credit/sync | 学分同步（占位） |
| POST | /api/role/assign | 管理员分配链上角色（需 Token） |
| GET  | /api/role/get | 查询链上角色（需 Token） |

---

## 产品设计资源

- **需求与角色边界**：`01-product-design/02-requirement-analysis/demand-role-analysis.md`
- **墨刀原型**：[校园区块链学分存证 DApp 低保真原型](https://modao.cc/proto/C0Jq60r2ta15q9aISKMefb/sharing?view_mode=read_only&screen=rbpVAR7MgSqPv2yaz)
- **问卷与截图**：`01-product-design/01-research-questionnaire/`、`03-mock-prototype/`

---

## 测试

- **合约单元测试**：
  ```bash
  cd 02-tech-development/01-smart-contract
  npx hardhat test test/unit-tests/credit-test.js
  ```
- 后端接口可用 Postman 按上表逐个测试（登录后 Header 带 `Authorization: Bearer <token>`）。

---

## 后续可做

- 将合约部署至 **Sepolia** 等测试网，并更新后端/前端配置与文档。
- 前端打包部署至 **Vercel / GitHub Pages**，提供公网访问地址。
- 本地 **Fabric 2.0** 联盟链部署，实现「公链 + 联盟链」双部署。
- Gas 与前端加载等优化，并在 README 中补充测试网地址、截图与使用说明。

---

## License

MIT
