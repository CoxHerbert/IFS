<template>
  <div class="login">
    <a-form ref="loginRef" :model="loginForm" :rules="loginRules" class="login-form" layout="vertical">
      <h3 class="title">IFS 后台管理系统</h3>
      <a-form-item name="username">
        <div class="login-logo">
          <img src="../assets/images/loginLogo.png" alt="">
        </div>
        <a-input v-model:value="loginForm.username" autocomplete="off" placeholder="账号">
          <template #prefix><svg-icon icon-class="user" class="input-icon" /></template>
        </a-input>
      </a-form-item>
      <a-form-item name="password">
        <a-input-password v-model:value="loginForm.password" autocomplete="off" placeholder="密码"
          @pressEnter="handleLogin">
          <template #prefix><svg-icon icon-class="password" class="input-icon" /></template>
        </a-input-password>
      </a-form-item>
      <a-form-item v-if="captchaOnOff" name="code">
        <div class="captcha-row">
          <a-input v-model:value="loginForm.code" autocomplete="off" placeholder="验证码" @pressEnter="handleLogin">
            <template #prefix><svg-icon icon-class="validCode" class="input-icon" /></template>
          </a-input>
          <div class="login-code">
            <img :src="codeUrl" class="login-code-img" @click="getCode">
          </div>
        </div>
      </a-form-item>
      <a-checkbox v-model:checked="loginForm.rememberMe" class="remember-checkbox">记住密码</a-checkbox>
      <a-form-item>
        <a-button :loading="loading" type="primary" block @click.prevent="handleLogin">
          {{ loading ? "登录中..." : "登录" }}
        </a-button>
        <div v-if="register" class="register-link">
          <router-link class="link-type" to="/register">立即注册</router-link>
        </div>
      </a-form-item>
    </a-form>
    <div class="login-footer">
      <span>Copyright © IFS All Rights Reserved.</span>
    </div>
  </div>
</template>

<script setup>
import { getCodeImg } from "@/api/login";
import Cookies from "js-cookie";
import { encrypt, decrypt } from "@/utils/jsencrypt";

const store = useStore();
const router = useRouter();
const { proxy } = getCurrentInstance();

const loginForm = ref({
  username: "admin",
  password: "admin123",
  rememberMe: false,
  code: "",
  uuid: ""
});

const loginRules = {
  username: [{ required: true, trigger: "blur", message: "请输入您的账号" }],
  password: [{ required: true, trigger: "blur", message: "请输入您的密码" }],
  code: [{ required: true, trigger: "change", message: "请输入验证码" }]
};

const codeUrl = ref("");
const loading = ref(false);
const captchaOnOff = ref(true);
const register = ref(false);
const redirect = ref(undefined);

function handleLogin() {
  proxy.$refs.loginRef.validate().then(() => {
    loading.value = true;
    if (loginForm.value.rememberMe) {
      Cookies.set("username", loginForm.value.username, { expires: 30 });
      Cookies.set("password", encrypt(loginForm.value.password), { expires: 30 });
      Cookies.set("rememberMe", loginForm.value.rememberMe, { expires: 30 });
    } else {
      Cookies.remove("username");
      Cookies.remove("password");
      Cookies.remove("rememberMe");
    }
    store.dispatch("Login", loginForm.value).then(() => {
      router.push({ path: redirect.value || "/" });
    }).catch(() => {
      loading.value = false;
      if (captchaOnOff.value) {
        getCode();
      }
    });
  }).catch(() => { });
}

function getCode() {
  getCodeImg().then(res => {
    captchaOnOff.value = res.captchaOnOff === undefined ? true : res.captchaOnOff;
    if (captchaOnOff.value) {
      codeUrl.value = res.data.img;
      loginForm.value.uuid = res.data.uuid;
    }
  });
}

function getCookie() {
  const username = Cookies.get("username");
  const password = Cookies.get("password");
  const rememberMe = Cookies.get("rememberMe");
  loginForm.value = {
    ...loginForm.value,
    username: username === undefined ? loginForm.value.username : username,
    password: password === undefined ? loginForm.value.password : decrypt(password),
    rememberMe: rememberMe === undefined ? false : Boolean(rememberMe)
  };
}

getCode();
getCookie();
</script>

<style lang="scss" scoped>
.login-logo {
  position: absolute;
  top: -211px;
  left: -24px;
}

.login {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  background: #efefef;
}

.title {
  margin: 0 auto 30px;
  text-align: center;
  color: #707070;
}

.login-form {
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

.remember-checkbox {
  margin: 0 0 25px;
}

.register-link {
  margin-top: 12px;
  text-align: right;
}

.login-code img {
  width: 100%;
  cursor: pointer;
}

.login-footer {
  height: 40px;
  line-height: 40px;
  position: fixed;
  bottom: 0;
  width: 100%;
  text-align: center;
  color: #fff;
  font-family: Arial;
  font-size: 12px;
  letter-spacing: 1px;
}

.login-code-img {
  height: 38px;
}
</style>
