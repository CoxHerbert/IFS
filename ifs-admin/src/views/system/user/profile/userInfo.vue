<template>
  <a-form ref="userRef" :model="user" :rules="rules" :label-col="{ style: { width: '80px' } }">
    <a-form-item label="用户昵称" name="nickName">
      <a-input v-model:value="user.nickName" :maxlength="30" />
    </a-form-item>
    <a-form-item label="手机号码" name="phonenumber">
      <a-input v-model:value="user.phonenumber" :maxlength="11" />
    </a-form-item>
    <a-form-item label="邮箱" name="email">
      <a-input v-model:value="user.email" :maxlength="50" />
    </a-form-item>
    <a-form-item label="性别" name="sex">
      <a-radio-group v-model:value="user.sex">
        <a-radio value="0">男</a-radio>
        <a-radio value="1">女</a-radio>
      </a-radio-group>
    </a-form-item>
    <a-form-item>
      <a-space>
        <a-button type="primary" @click="submit">保存</a-button>
        <a-button danger @click="close">关闭</a-button>
      </a-space>
    </a-form-item>
  </a-form>
</template>

<script setup>
import { updateUserProfile } from "@/api/system/user";

const props = defineProps({
  user: {
    type: Object
  }
});

const { proxy } = getCurrentInstance();

const rules = ref({
  nickName: [{ required: true, message: "用户昵称不能为空", trigger: "blur" }],
  email: [
    { required: true, message: "邮箱地址不能为空", trigger: "blur" },
    { type: "email", message: "请输入正确的邮箱地址", trigger: ["blur", "change"] }
  ],
  phonenumber: [
    { required: true, message: "手机号码不能为空", trigger: "blur" },
    { pattern: /^1[3-9][0-9]\d{8}$/, message: "请输入正确的手机号码", trigger: "blur" }
  ]
});

function submit() {
  proxy.$refs.userRef.validate().then(() => {
    updateUserProfile(props.user).then(() => {
      proxy.$modal.msgSuccess("修改成功");
    });
  }).catch(() => {});
}

function close() {
  proxy.$tab.closePage();
}
</script>
