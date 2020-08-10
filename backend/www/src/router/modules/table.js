/** When your routing table is too long, you can split it into small modules **/

import Layout from '@/layout'

const tableRouter = {
  path: '/table',
  component: Layout,
  redirect: '/table/complex-table',
  name: 'Table',
  meta: {
    title: 'Table',
    icon: 'table'
  },
  children: [
    // {
    //   path: 'dynamic-table',
    //   component: () => import('@/views/table/dynamic-table/index'),
    //   name: 'DynamicTable',
    //   meta: { title: 'Dynamic Table' }
    // },
    // {
    //   path: 'drag-table',
    //   component: () => import('@/views/table/drag-table'),
    //   name: 'DragTable',
    //   meta: { title: 'Drag Table' }
    // },
    // {
    //   path: 'inline-edit-table',
    //   component: () => import('@/views/table/inline-edit-table'),
    //   name: 'InlineEditTable',
    //   meta: { title: 'Inline Edit' }
    // },
    {
      path: 'complex-table',
      component: () => import('@/views/table/complex-table'),
      name: '拉流管理',
      meta: { title: '拉流管理' }
    },
    {
      path: 'play/flv/:id(\\d+)',
      component: () => import('@/views/table/flv'),
      name: 'flv播放视频',
      meta: { title: 'flv播放视频', noCache: true, activeMenu: '/table/complex-table' },
      hidden: true
    },
    {
      path: 'play/m3u8/:id(\\d+)',
      component: () => import('@/views/table/m3u8'),
      name: 'm3u8播放视频',
      meta: { title: 'm3u8播放视频', noCache: true, activeMenu: '/table/complex-table' },
      hidden: true
    }
  ]
}
export default tableRouter
