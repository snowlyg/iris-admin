import request from '@/utils/request'

export function fetchList(query) {
  return request({
    url: '/vue-element-admin/article/list',
    method: 'get',
    params: query
  })
}

export function fetchArticle(id) {
  return request({
    url: `/vue-element-admin/article/${id}`,
    method: 'get'
  })
}

export function startArticle(id) {
  return request({
    url: `/vue-element-admin/article/start/${id}`,
    method: 'get'
  })
}
export function stopArticle(id) {
  return request({
    url: `/vue-element-admin/article/stop/${id}`,
    method: 'get'
  })
}

export function fetchPv(pv) {
  return request({
    url: '/vue-element-admin/article/pv',
    method: 'get',
    params: { pv }
  })
}

export function createArticle(data) {
  return request({
    url: '/vue-element-admin/article/create',
    method: 'post',
    data
  })
}

export function updateArticle(data, id) {
  return request({
    url: `/vue-element-admin/article/${id}`,
    method: 'post',
    data
  })
}

export function deleteArticle(id) {
  return request({
    url: `/vue-element-admin/article/${id}`,
    method: 'delete'
  })
}
