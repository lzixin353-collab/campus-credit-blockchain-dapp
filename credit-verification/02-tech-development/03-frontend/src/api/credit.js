import request from './request'

// 用户登录（钱包地址）
export const login = (address) => {
  return request({
    url: '/user/login',
    method: 'post',
    data: { address }
  })
}

// 同步链上学分到后端
export const syncCredit = () => {
  return request({
    url: '/credit/sync',
    method: 'post'
  })
}

// 查询学分列表（按角色：教师看自己录入的）
export const getCreditList = (params) => {
  return request({
    url: '/credit/list',
    method: 'get',
    params
  })
}

// 教师录入学分（后端上链并落库）
export const recordCredit = (data) => {
  return request({
    url: '/credit/record',
    method: 'post',
    data: {
      student_address: data.student_address,
      course_name: data.course_name,
      score: data.score
    }
  })
}

// 管理员：待审核学分列表
export const getCreditPending = () => {
  return request({ url: '/credit/pending', method: 'get' })
}

// 管理员：审核通过（credit_id 为数据库主键 id）
export const approveCredit = (creditId) => {
  return request({ url: '/credit/approve', method: 'post', data: { credit_id: creditId } })
}

// 管理员：驳回
export const rejectCredit = (creditId) => {
  return request({ url: '/credit/reject', method: 'post', data: { credit_id: creditId } })
}

// 管理员：分配角色（后端调合约）
export const assignRole = (user_address, role) => {
  return request({
    url: '/role/assign',
    method: 'post',
    data: { user_address, role }
  })
}