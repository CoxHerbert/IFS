import { createApp } from 'vue'
import Antd from 'ant-design-vue'
import VXETable from 'vxe-table'
import 'ant-design-vue/dist/reset.css'
import 'vxe-table/lib/style.css'
import router from './router'
import App from './App.vue'
import './style.css'

createApp(App).use(router).use(Antd).use(VXETable).mount('#app')
