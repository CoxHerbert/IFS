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
              <span>400-888-2026</span>
            </div>
            <div class="info-item">
              <MailOutlined />
              <span>quote@seawaypro.com</span>
            </div>
            <div class="info-item">
              <EnvironmentOutlined />
              <span>上海市徐汇区</span>
            </div>
            <div class="info-item">
              <ClockCircleOutlined />
              <span>工作日 09:00-18:00</span>
            </div>
          </div>

          <a-divider />

          <div class="promise">
            <h4>响应承诺</h4>
            <p>工作时间内提交需求后，我们会优先核对起运港、目的港、货物属性和时效要求，再给出可执行的运输方案。</p>
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
  </main>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import type { FormInstance } from 'ant-design-vue'
import {
  ClockCircleOutlined,
  EnvironmentOutlined,
  MailOutlined,
  PhoneOutlined
} from '@ant-design/icons-vue'
import { submitContact, type ContactPayload } from '@/api/contact'

const formRef = ref<FormInstance>()
const submitting = ref(false)

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

@media (max-width: 760px) {
  .page {
    width: min(100% - 20px, 1240px);
  }

  .hero-copy {
    padding: 28px;
  }
}
</style>
