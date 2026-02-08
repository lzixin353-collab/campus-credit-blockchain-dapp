require("@nomiclabs/hardhat-waffle");
require("@nomiclabs/hardhat-etherscan");
require("dotenv").config();

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: {
    version: "0.8.21", // 和我们写的合约版本一致
    settings: {
      optimizer: {
        enabled: true,
        runs: 200 // 简单的Gas优化
      }
    }
  },
  // 测试网配置（后续部署用，现在先留空也可以）
  networks: {
    // 本地测试节点
    localhost: {
      url: "http://127.0.0.1:8545",
    },
    sepolia: {
      url: `https://sepolia.infura.io/v3/${process.env.INFURA_API_KEY}`,
      accounts: [process.env.PRIVATE_KEY],
      gas: 2100000,
      gasPrice: 8000000000,
    },
  },
  // etherscan: {
  //   apiKey: process.env.ETHERSCAN_API_KEY || "",
  // },
};