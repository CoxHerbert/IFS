<template>
  <div class="errPage-container">
    <a-button class="pan-back-btn" @click="back">返回</a-button>
    <a-row :gutter="24">
      <a-col :span="12">
        <h1 class="text-jumbo text-ginormous">401 错误!</h1>
        <h2>您没有访问权限</h2>
        <h6>抱歉，您没有访问权限，请不要进行非法操作。您可以返回首页。</h6>
        <ul class="list-unstyled">
          <li class="link-type">
            <router-link to="/">回首页</router-link>
          </li>
        </ul>
      </a-col>
      <a-col :span="12">
        <img :src="errGif" width="313" height="428" alt="Girl has dropped her ice cream.">
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import errImage from "@/assets/401_images/401.gif";

let { proxy } = getCurrentInstance();

const errGif = ref(errImage + "?" + +new Date());

function back() {
  if (proxy.$route.query.noGoBack) {
    proxy.$router.push({ path: "/" });
  } else {
    proxy.$router.go(-1);
  }
}
</script>

<style lang="scss" scoped>
.errPage-container {
  width: 800px;
  max-width: 100%;
  margin: 100px auto;

  .pan-back-btn {
    background: #008489;
    color: #fff;
    border: none !important;
  }

  .text-jumbo {
    font-size: 60px;
    font-weight: 700;
    color: #484848;
  }

  .list-unstyled {
    font-size: 14px;

    li {
      padding-bottom: 5px;
    }

    a {
      color: #008489;
      text-decoration: none;

      &:hover {
        text-decoration: underline;
      }
    }
  }
}
</style>
