// @ts-nocheck
import store from '@/store'
import router from '@/router'

export default {
  // 刷新当前 tab 页签
  refreshPage(obj) {
    const { path, matched } = router.currentRoute.value;
    if (obj === undefined) {
      matched.forEach((m) => {
        if (m.components && m.components.default && m.components.default.name) {
          if (!['Layout', 'ParentView'].includes(m.components.default.name)) {
            obj = { name: m.components.default.name, path: path };
          }
        }
      });
    }
    return store.dispatch('tagsView/delCachedView', obj).then(() => {
      const { path } = obj
      router.replace({
        path: '/redirect' + path
      })
    })
  },
  // 关闭当前 tab 页签并打开新页签
  closeOpenPage(obj) {
    store.dispatch("tagsView/delView", router.currentRoute.value);
    if (obj !== undefined) {
      return router.push(obj);
    }
  },
  // 关闭指定 tab 页签
  closePage(obj) {
    if (obj === undefined) {
      return store.dispatch('tagsView/delView', router.currentRoute.value).then(({ lastPath }) => {
        return router.push(lastPath || '/index');
      });
    }
    return store.dispatch('tagsView/delView', obj);
  },
  // 关闭全部 tab 页签
  closeAllPage() {
    return store.dispatch('tagsView/delAllViews');
  },
  // 关闭左侧 tab 页签
  closeLeftPage(obj) {
    return store.dispatch('tagsView/delLeftTags', obj || router.currentRoute.value);
  },
  // 关闭右侧 tab 页签
  closeRightPage(obj) {
    return store.dispatch('tagsView/delRightTags', obj || router.currentRoute.value);
  },
  // 关闭其他 tab 页签
  closeOtherPage(obj) {
    return store.dispatch('tagsView/delOthersViews', obj || router.currentRoute.value);
  },
  // 打开 tab 页签
  openPage(url) {
    return router.push(url);
  },
  // 修改 tab 页签
  updatePage(obj) {
    return store.dispatch('tagsView/updateVisitedView', obj);
  }
}


