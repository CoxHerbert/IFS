<template>
  <div class="app-container">
    <a-row :gutter="20">
      <a-col :span="6" :xs="24">
        <a-card title="个人信息" class="box-card">
          <div class="text-center">
            <userAvatar :user="state.user" />
          </div>
          <ul class="list-group list-group-striped">
            <li class="list-group-item">
              <svg-icon icon-class="user" />用户名
              <div class="pull-right">{{ state.user.userName }}</div>
            </li>
            <li class="list-group-item">
              <svg-icon icon-class="phone" />手机号码
              <div class="pull-right">{{ state.user.phonenumber }}</div>
            </li>
            <li class="list-group-item">
              <svg-icon icon-class="email" />用户邮箱
              <div class="pull-right">{{ state.user.email }}</div>
            </li>
            <li class="list-group-item">
              <svg-icon icon-class="tree" />所属部门
              <div class="pull-right">{{ state.user.deptName }} / {{ state.postGroup }}</div>
            </li>
            <li class="list-group-item">
              <svg-icon icon-class="peoples" />所属角色
              <div class="pull-right">{{ state.roleGroup }}</div>
            </li>
            <li class="list-group-item">
              <svg-icon icon-class="date" />创建日期
              <div class="pull-right">{{ state.user.createTime }}</div>
            </li>
          </ul>
        </a-card>
      </a-col>
      <a-col :span="18" :xs="24">
        <a-card title="基本资料">
          <a-tabs v-model:activeKey="activeTab">
            <a-tab-pane key="userinfo" tab="基本资料">
              <userInfo :user="state.user" />
            </a-tab-pane>
            <a-tab-pane key="resetPwd" tab="修改密码">
              <resetPwd :user="state.user" />
            </a-tab-pane>
          </a-tabs>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup name="Profile">
import userAvatar from "./userAvatar";
import userInfo from "./userInfo";
import resetPwd from "./resetPwd";
import { getUserProfile } from "@/api/system/user";

const activeTab = ref("userinfo");
const state = reactive({
  user: {},
  roleGroup: {},
  postGroup: {}
});

function getUser() {
  getUserProfile().then(response => {
    state.user = response.data.user;
    state.roleGroup = response.data.roleGroup;
    state.postGroup = response.data.postGroup;
  });
}

getUser();
</script>
