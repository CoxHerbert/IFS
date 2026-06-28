import { createApp } from 'vue'
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/antd.css'
import router from './router'
import App from './App.vue'
import './style.css'

createApp(App).use(router).use(Antd).mount('#app')
