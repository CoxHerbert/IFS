<template>
  <main class="page">
    <section class="hero">
      <div class="hero-copy">
        <a-tag color="blue">联系我们</a-tag>
        <h2>告诉我们你的运输需求，销售顾问会尽快联系你</h2>
        <p>填写航线、货物信息和联系方式后，系统会生成线索编号，方便后续报价和跟进。</p>
      </div>
    </section>

    <a-row :gutter="[18, 18]" class="section">
      <a-col :xs="24" :lg="9">
        <a-card class="panel contact-info" :bordered="false">
          <h3>联系信息</h3>
          <div class="info-list">
            <div class="info-item">
              <PhoneOutlined />
              <span>{{ contactConfig.contact.phone }}</span>
            </div>
            <div class="info-item">
              <MailOutlined />
              <span>{{ contactConfig.contact.email }}</span>
            </div>
            <div class="info-item">
              <EnvironmentOutlined />
              <span>{{ contactConfig.contact.address }}</span>
            </div>
            <div class="info-item">
              <ClockCircleOutlined />
              <span>{{ contactConfig.contact.businessHours }}</span>
            </div>
          </div>

          <a-divider />

          <div class="promise">
            <h4>响应承诺</h4>
            <p>{{ contactConfig.contact.responsePromise }}</p>
          </div>
        </a-card>
      </a-col>

      <a-col :xs="24" :lg="15">
        <a-card class="panel" :bordered="false">
          <a-form ref="formRef" layout="vertical" :model="formState" :rules="rules" @finish="handleSubmit">
            <a-row :gutter="[16, 0]">
              <a-col :xs="24" :md="12">
                <a-form-item label="联系人" name="contactName">
                  <a-input v-model:value="formState.contactName" placeholder="请输入联系人姓名" />
                </a-form-item>
              </a-col>
              <a-col :xs="24" :md="12">
                <a-form-item label="公司名称" name="companyName">
                  <a-input v-model:value="formState.companyName" placeholder="请输入公司名称" />
                </a-form-item>
              </a-col>
              <a-col :xs="24" :md="12">
                <a-form-item label="联系电话" name="phone">
                  <a-input v-model:value="formState.phone" placeholder="手机号 / 电话 / 微信" />
                </a-form-item>
              </a-col>
              <a-col :xs="24" :md="12">
                <a-form-item label="邮箱" name="email">
                  <a-input v-model:value="formState.email" placeholder="用于接收报价单" />
                </a-form-item>
              </a-col>
              <a-col :xs="24" :md="12">
                <a-form-item label="目标航线" name="route">
                  <a-input v-model:value="formState.route" placeholder="例如：上海 - 洛杉矶" />
                </a-form-item>
              </a-col>
              <a-col :xs="24" :md="12">
                <a-form-item label="货物信息" name="cargoInfo">
                  <a-input v-model:value="formState.cargoInfo" placeholder="品名 / 件数 / 体积 / 重量" />
                </a-form-item>
              </a-col>
              <a-col :span="24">
                <a-form-item label="需求说明" name="message">
                  <a-textarea
                    v-model:value="formState.message"
                    :rows="5"
                    placeholder="请描述运输方式、时效、是否报关、是否门到门等需求"
                  />
                </a-form-item>
              </a-col>
            </a-row>

            <div class="form-actions">
              <a-button type="primary" html-type="submit" size="large" :loading="submitting">
                提交需求
              </a-button>
              <a-button size="large" @click="resetForm">清空</a-button>
            </div>
          </a-form>
        </a-card>
      </a-col>
    </a-row>

    <a-card class="panel location-panel" :bordered="false">
      <div class="location-heading">
        <div>
          <span class="location-kicker">位置信息</span>
          <h3>欢迎到访</h3>
          <p><EnvironmentOutlined /> {{ mapAddress }}</p>
        </div>
        <a
          class="map-link"
          :href="amapMarkerUrl"
          target="_blank"
          rel="noopener noreferrer"
        >
          在高德地图中查看
        </a>
      </div>
      <div class="map-shell">
        <iframe
          class="map-fallback"
          :src="fallbackMapUrl"
          title="公司位置地图"
          loading="lazy"
          referrerpolicy="no-referrer-when-downgrade"
        ></iframe>
        <div
          ref="mapContainer"
          class="amap-map"
          :class="{ 'is-ready': mapReady }"
          aria-label="公司位置地图"
        ></div>
        <span v-if="mapError && !mapReady" class="map-fallback-tip">地图已切换至备用服务</span>
      </div>
    </a-card>
  </main>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import type { FormInstance } from 'ant-design-vue'
import {
  ClockCircleOutlined,
  EnvironmentOutlined,
  MailOutlined,
  PhoneOutlined
} from '@ant-design/icons-vue'
import { submitContact, type ContactPayload } from '@/api/portal/contact'
import contactConfig from '@/config/contact.json'

const formRef = ref<FormInstance>()
const submitting = ref(false)
const mapContainer = ref<HTMLDivElement>()
const mapError = ref('')
const mapReady = ref(false)
const mapAddress = contactConfig.location.address
const mapLongitude = contactConfig.location.longitude
const mapLatitude = contactConfig.location.latitude
const amapMarkerUrl = `https://uri.amap.com/marker?position=${mapLongitude},${mapLatitude}&name=${encodeURIComponent(mapAddress)}&coordinate=gaode&callnative=0`
const fallbackMapUrl = `https://www.openstreetmap.org/export/embed.html?bbox=${mapLongitude - 0.015}%2C${mapLatitude - 0.01}%2C${mapLongitude + 0.015}%2C${mapLatitude + 0.01}&layer=mapnik&marker=${mapLatitude}%2C${mapLongitude}`
let map: any

function loadAmap(key: string, securityJsCode: string) {
  return new Promise<void>((resolve, reject) => {
    if ((window as any).AMap?.Map) {
      resolve()
      return
    }

    ;(window as any)._AMapSecurityConfig = { securityJsCode }

    const existingScript = document.querySelector<HTMLScriptElement>('#amap-js-api')
    if (existingScript) {
      existingScript.addEventListener('load', () => resolve(), { once: true })
      existingScript.addEventListener('error', () => reject(new Error('高德地图服务加载失败')), { once: true })
      return
    }

    const script = document.createElement('script')
    script.id = 'amap-js-api'
    script.src = `https://webapi.amap.com/maps?v=2.0&key=${encodeURIComponent(key)}&plugin=AMap.ToolBar`
    script.async = true
    script.onload = () => resolve()
    script.onerror = () => reject(new Error('高德地图服务加载失败'))
    document.head.appendChild(script)
  })
}

async function initMap() {
  const key = import.meta.env.VITE_AMAP_KEY?.trim()
  const securityJsCode = import.meta.env.VITE_AMAP_SECURITY_JS_CODE?.trim()
  if (!key || !securityJsCode) {
    mapError.value = '请配置高德地图 Key 和安全密钥后查看交互地图'
    return
  }

  try {
    await loadAmap(key, securityJsCode)
    if (!mapContainer.value) return

    const AMap = (window as any).AMap
    const position = [mapLongitude, mapLatitude]
    map = new AMap.Map(mapContainer.value, {
      center: position,
      zoom: contactConfig.location.zoom,
      viewMode: '2D',
    })
    map.on('complete', () => {
      mapReady.value = true
      mapError.value = ''
    })
    map.add(new AMap.Marker({ position, title: mapAddress }))
    map.addControl(new AMap.ToolBar({ position: 'RT' }))
  } catch (error) {
    mapError.value = error instanceof Error ? error.message : '高德地图服务加载失败'
  }
}

onMounted(initMap)
onBeforeUnmount(() => {
  map?.destroy?.()
  map = undefined
})

const initialForm: ContactPayload = {
  contactName: '',
  companyName: '',
  phone: '',
  email: '',
  route: '',
  cargoInfo: '',
  message: '',
  source: 'portal-contact'
}

const formState = reactive<ContactPayload>({ ...initialForm })

const rules = {
  contactName: [{ required: true, message: '请输入联系人', trigger: 'blur' }],
  phone: [
    {
      validator: async () => {
        if (!formState.phone && !formState.email) {
          return Promise.reject(new Error('请至少填写电话或邮箱'))
        }
        return Promise.resolve()
      },
      trigger: 'blur'
    }
  ],
  email: [{ type: 'email', message: '邮箱格式不正确', trigger: 'blur' }],
  message: [{ required: true, message: '请填写需求说明', trigger: 'blur' }]
}

async function handleSubmit() {
  submitting.value = true
  try {
    const result = await submitContact(formState)
    if (result.code !== 200) {
      throw new Error(result.msg || '提交失败')
    }
    message.success(`提交成功，线索编号：${result.data?.leadNo || '已生成'}`)
    resetForm()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '提交失败，请稍后再试')
  } finally {
    submitting.value = false
  }
}

function resetForm() {
  Object.assign(formState, initialForm)
  formRef.value?.clearValidate()
}
</script>

<style scoped>
.page {
  width: min(1240px, calc(100% - 32px));
  margin: 0 auto;
  padding: 28px 0 40px;
}

.hero,
.panel {
  border-radius: 22px;
  box-shadow: 0 18px 40px rgba(16, 35, 63, 0.08);
}

.hero {
  overflow: hidden;
  background:
    linear-gradient(135deg, rgba(7, 23, 47, 0.92), rgba(22, 119, 255, 0.78)),
    url("@/assets/hero.jpg") center / cover no-repeat;
}

.hero-copy {
  padding: 46px;
  color: #fff;
}

.hero-copy h2 {
  max-width: 720px;
  margin: 16px 0 0;
  font-size: clamp(30px, 4vw, 52px);
  line-height: 1.08;
}

.hero-copy p {
  max-width: 56ch;
  margin: 18px 0 0;
  color: rgba(255, 255, 255, 0.84);
  line-height: 1.8;
}

.section {
  margin-top: 18px;
}

.panel {
  min-height: 100%;
}

.panel :deep(.ant-card-body) {
  padding: 26px;
}

.contact-info h3 {
  margin: 0;
  font-size: 24px;
}

.info-list {
  display: grid;
  gap: 14px;
  margin-top: 22px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 12px;
  color: #41516a;
}

.info-item :deep(svg) {
  color: #1677ff;
  font-size: 18px;
}

.promise h4 {
  margin: 0;
  font-size: 18px;
}

.promise p {
  margin: 10px 0 0;
  color: #66748b;
  line-height: 1.8;
}

.form-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.location-panel {
  margin-top: 18px;
}

.location-heading {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 20px;
}

.location-kicker {
  color: #1677ff;
  font-weight: 600;
}

.location-heading h3 {
  margin: 6px 0 8px;
  font-size: 24px;
}

.location-heading p {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
  color: #66748b;
}

.map-link {
  flex: none;
}

.map-shell {
  position: relative;
  min-height: 380px;
  overflow: hidden;
  border: 1px solid #e8edf4;
  border-radius: 16px;
  background: #f5f8fc;
}

.amap-map {
  position: absolute;
  inset: 0;
  z-index: 1;
  width: 100%;
  height: 380px;
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.2s ease;
}

.amap-map.is-ready {
  opacity: 1;
  pointer-events: auto;
}

.map-fallback {
  display: block;
  width: 100%;
  height: 380px;
  border: 0;
}

.map-fallback-tip {
  position: absolute;
  right: 12px;
  bottom: 12px;
  z-index: 2;
  padding: 6px 10px;
  border-radius: 6px;
  color: #66748b;
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 2px 8px rgba(16, 35, 63, 0.12);
  font-size: 12px;
}

@media (max-width: 760px) {
  .page {
    width: min(100% - 20px, 1240px);
  }

  .hero-copy {
    padding: 28px;
  }

  .location-heading {
    align-items: flex-start;
    flex-direction: column;
  }

  .map-shell,
  .amap-map,
  .map-fallback {
    height: 300px;
    min-height: 300px;
  }
}
</style>
