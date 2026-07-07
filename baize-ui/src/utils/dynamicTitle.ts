// @ts-nocheck
import store from '@/store'
import defaultSettings from '@/settings'

/**
 * 鍔ㄦ€佷慨鏀规爣棰?
 */
export function useDynamicTitle() {
  if (store.state.settings.dynamicTitle) {
    document.title = store.state.settings.title + ' - ' + defaultSettings.title;
  } else {
    document.title = defaultSettings.title;
  }
}

