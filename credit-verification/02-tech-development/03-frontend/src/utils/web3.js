import Web3 from 'web3'
// 只导入CreditContract的ABI（删除RoleContract导入）
import { default as creditContractJson } from '@/assets/abi/credit_contract.json'

// 1. 提取ABI
const creditContractABI = creditContractJson.abi || []

// 2. 仅保留CreditContract地址
const CREDIT_CONTRACT_ADDRESS = '0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512' // 你的本地部署地址

// Hardhat 本地网 chainId
const HARDHAT_CHAIN_ID = 31337
const HARDHAT_CHAIN_ID_HEX = '0x7A69'
const HARDHAT_RPC = 'http://127.0.0.1:8545'

// 3. 初始化Web3（适配Metamask+本地测试网）
export const initWeb3 = async () => {
  if (!window.ethereum) {
    alert('请安装Metamask钱包！')
    return null
  }
  try {
    await window.ethereum.request({ method: 'eth_requestAccounts' })
    const web3 = new Web3(window.ethereum)
    let chainId = await web3.eth.getChainId()
    if (chainId !== HARDHAT_CHAIN_ID) {
      try {
        await window.ethereum.request({
          method: 'wallet_switchEthereumChain',
          params: [{ chainId: HARDHAT_CHAIN_ID_HEX }]
        })
      } catch (switchErr) {
        // 4902 = 该链未添加到 Metamask，先添加再切换
        if (switchErr.code === 4902 || switchErr.message?.includes('not been added')) {
          await window.ethereum.request({
            method: 'wallet_addEthereumChain',
            params: [{
              chainId: HARDHAT_CHAIN_ID_HEX,
              chainName: 'Hardhat Local',
              nativeCurrency: { name: 'ETH', symbol: 'ETH', decimals: 18 },
              rpcUrls: [HARDHAT_RPC]
            }]
          })
          await window.ethereum.request({
            method: 'wallet_switchEthereumChain',
            params: [{ chainId: HARDHAT_CHAIN_ID_HEX }]
          })
        } else {
          alert('请切换到 Hardhat 本地测试网（链 ID 31337），或确保已启动 npx hardhat node。')
          return null
        }
      }
    }
    return web3
  } catch (error) {
    console.error('Web3初始化失败：', error)
    alert('连接钱包失败，请安装Metamask并授权！')
    return null
  }
}

// 4. 获取当前钱包地址
export const getCurrentAddress = async (web3) => {
  if (!web3) return ''
  const accounts = await web3.eth.getAccounts()
  return accounts[0] || ''
}

// 5. 获取用户角色（调用CreditContract的getRole方法，合约返回字符串 "teacher"/"admin"/""）
export const getUserRole = async (web3, address) => {
  if (!web3 || !address) return ''
  try {
    const creditContract = new web3.eth.Contract(creditContractABI, CREDIT_CONTRACT_ADDRESS)
    const role = await creditContract.methods.getRole(address).call()
    // 合约直接返回 "teacher" | "admin" | ""（空为学生）
    if (role === 'teacher') return 'teacher'
    if (role === 'admin') return 'admin'
    return 'student'
  } catch (error) {
    console.error('获取角色失败：', error)
    return 'student'
  }
}

// 6. 分配角色（调用CreditContract的assignRole方法）
export const assignRole = async (web3, adminAddress, userAddress, role) => {
  if (!web3 || !adminAddress || !userAddress) return false
  try {
    const creditContract = new web3.eth.Contract(creditContractABI, CREDIT_CONTRACT_ADDRESS)
    // role参数转换：student=0, teacher=1, admin=2
    await creditContract.methods.assignRole(userAddress, role)
      .send({ from: adminAddress })
    return true
  } catch (error) {
    console.error('分配角色失败：', error)
    return false
  }
}

// 7. 获取学生学分
export const getStudentCredit = async (web3, studentAddress) => {
  if (!web3 || !studentAddress) return []
  try {
    const creditContract = new web3.eth.Contract(creditContractABI, CREDIT_CONTRACT_ADDRESS)
    const credits = await creditContract.methods.getCreditByStudent(studentAddress).call()
    // 格式化数据（适配你的合约返回结构）
    return credits.map(item => ({
      creditId: item.id || '',
      courseName: item.courseName || '',
      creditScore: web3.utils.fromWei(item.creditScore || '0', 'ether'),
      teacherAddress: item.teacherAddress || '',
      status: item.status || '0', // 0=待审核，1=已通过，2=已驳回
      timestamp: item.timestamp ? new Date(item.timestamp * 1000).toLocaleString() : ''
    }))
  } catch (error) {
    console.error('获取学分失败：', error)
    return []
  }
}

// 8. 教师录入学分
export const teacherInputCredit = async (web3, teacherAddress, studentAddress, courseName, creditScore) => {
  if (!web3 || !teacherAddress || !studentAddress) return false
  try {
    const creditContract = new web3.eth.Contract(creditContractABI, CREDIT_CONTRACT_ADDRESS)
    await creditContract.methods.inputCredit(studentAddress, courseName, web3.utils.toWei(creditScore, 'ether'))
      .send({ from: teacherAddress })
    return true
  } catch (error) {
    console.error('录入学分失败：', error)
    return false
  }
}

// 9. 管理员审核学分
export const adminAuditCredit = async (web3, adminAddress, creditId, isApproved) => {
  if (!web3 || !adminAddress || !creditId) return false
  try {
    const creditContract = new web3.eth.Contract(creditContractABI, CREDIT_CONTRACT_ADDRESS)
    // isApproved：true=通过（status=1），false=驳回（status=2）
    await creditContract.methods.auditCredit(creditId, isApproved)
      .send({ from: adminAddress })
    return true
  } catch (error) {
    console.error('审核学分失败：', error)
    return false
  }
}