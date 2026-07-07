<template>
  <div class="register">
    <a-form ref="registerRef" :model="registerForm" :rules="registerRules" class="register-form" layout="vertical">
      <h3 class="title">IFS 后台管理系统</h3>
      <a-form-item name="username">
        <a-input v-model:value="registerForm.username" autocomplete="off" placeholder="账号">
          <template #prefix><svg-icon icon-class="user" class="input-icon" /></template>
        </a-input>
      </a-form-item>
      <a-form-item name="password">
        <a-input-password v-model:value="registerForm.password" autocomplete="off" placeholder="密码"
          @pressEnter="handleRegister">
          <template #prefix><svg-icon icon-class="password" class="input-icon" /></template>
        </a-input-password>
      </a-form-item>
      <a-form-item name="confirmPassword">
        <a-input-password v-model:value="registerForm.confirmPassword" autocomplete="off" placeholder="确认密码"
          @pressEnter="handleRegister">
          <template #prefix><svg-icon icon-class="password" class="input-icon" /></template>
        </a-input-password>
      </a-form-item>
      <a-form-item v-if="captchaOnOff" name="code">
        <div class="captcha-row">
          <a-input v-model:value="registerForm.code" autocomplete="off" placeholder="验证码" @pressEnter="handleRegister">
            <template #prefix><svg-icon icon-class="validCode" class="input-icon" /></template>
          </a-input>
          <div class="register-code">
            <img :src="codeUrl" class="register-code-img" @click="getCode" />
          </div>
        </div>
      </a-form-item>
      <a-form-item>
        <a-button :loading="loading" type="primary" block @click.prevent="handleRegister">
          {{ loading ? "注册中..." : "注册" }}
        </a-button>
        <div class="register-link">
          <router-link class="link-type" to="/login">使用已有账户登录</router-link>
        </div>
      </a-form-item>
    </a-form>
    <div class="register-footer">
      <span>Copyright © 2018-2021 ruoyi.vip All Rights Reserved.</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Modal } from "ant-design-vue";
import { getCurrentInstance, ref } from "vue";
import { useRouter } from "vue-router";
import { getCodeImg, register } from "@/api/login";

const router = useRouter();
const { proxy } = getCurrentInstance() as any;

const registerForm = ref({
  username: "",
  password: "",
  confirmPassword: "",
  code: "",
  uuid: ""
});

const equalToPassword = (_rule: unknown, value: string, callback: (error?: Error) => void) => {
  if (registerForm.value.password !== value) {
    callback(new Error("两次输入的密码不一致"));
    return;
  }
  callback();
};

const registerRules = {
  username: [
    { required: true, trigger: "blur", message: "请输入您的账号" },
    { min: 2, max: 20, message: "用户账号长度必须介于 2 和 20 之间", trigger: "blur" }
  ],
  password: [
    { required: true, trigger: "blur", message: "请输入您的密码" },
    { min: 5, max: 20, message: "用户密码长度必须介于 5 和 20 之间", trigger: "blur" }
  ],
  confirmPassword: [
    { required: true, trigger: "blur", message: "请再次输入您的密码" },
    { required: true, validator: equalToPassword, trigger: "blur" }
  ],
  code: [{ required: true, trigger: "change", message: "请输入验证码" }]
};

const codeUrl = ref("");
const loading = ref(false);
const captchaOnOff = ref(true);

function handleRegister() {
  proxy.$refs.registerRef.validate().then(() => {
    loading.value = true;
    register(registerForm.value)
      .then(() => {
        const username = registerForm.value.username;
        Modal.success({
          title: "系统提示",
          content: `恭喜你，您的账号 ${username} 注册成功。`,
          onOk: () => router.push("/login")
        });
      })
      .catch(() => {
        loading.value = false;
        if (captchaOnOff.value) {
          getCode();
        }
      });
  }).catch(() => { });
}

function getCode() {
  getCodeImg().then((res: any) => {
    captchaOnOff.value = res.captchaOnOff === undefined ? true : res.captchaOnOff;
    if (captchaOnOff.value) {
      codeUrl.value = "data:image/gif;base64," + res.img;
      registerForm.value.uuid = res.uuid;
    }
  });
}

getCode();
</script>

<style lang="scss" scoped>
.register {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  background-image: url("../assets/images/login-background.jpg");
  background-size: cover;
}

.title {
  margin: 0 auto 30px;
  text-align: center;
  color: #707070;
}

.register-form {
  border-radius: 6px;
  background: #ffffff;
  width: 400px;
  padding: 25px 25px 5px;
}

.input-icon {
  width: 14px;
  margin-left: 2px;
}

.captcha-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 120px;
  gap: 12px;
  width: 100%;
}

.register-link {
  margin-top: 12px;
  text-align: right;
}

.register-code img {
  width: 100%;
  cursor: pointer;
}

.register-footer {
  height: 40px;
  line-height: 40px;
  position: fixed;
  bottom: 0;
  width: 100%;
  text-align: center;
  color: #fff;
  font-family: Arial, sans-serif;
  font-size: 12px;
  letter-spacing: 1px;
}

.register-code-img {
  height: 38px;
}
</style>
