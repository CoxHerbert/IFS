<template>
  <div class="user-info-head" @click="editCropper"><img :src="options.img" title="点击上传头像" class="img-circle img-lg" /></div>
  <a-modal v-model:open="open" :title="title" width="800px" :footer="null" destroy-on-close @afterOpenChange="handleAfterOpenChange">
    <a-row :gutter="16">
      <a-col :xs="24" :md="12" :style="{ height: '350px' }">
        <vue-cropper
          v-if="visible"
          ref="cropper"
          :img="options.img"
          :info="true"
          :autoCrop="options.autoCrop"
          :autoCropWidth="options.autoCropWidth"
          :autoCropHeight="options.autoCropHeight"
          :fixedBox="options.fixedBox"
          @realTime="realTime"
        />
      </a-col>
      <a-col :xs="24" :md="12" :style="{ height: '350px' }">
        <div class="avatar-upload-preview">
          <img :src="options.previews.url" :style="options.previews.img">
        </div>
      </a-col>
    </a-row>
    <br>
    <a-row :gutter="12" align="middle">
      <a-col :lg="4" :md="4">
        <a-upload :custom-request="requestUpload" :show-upload-list="false" :before-upload="beforeUpload">
          <a-button>选择</a-button>
        </a-upload>
      </a-col>
      <a-col :lg="3" :md="3"><a-button @click="changeScale(1)">放大</a-button></a-col>
      <a-col :lg="3" :md="3"><a-button @click="changeScale(-1)">缩小</a-button></a-col>
      <a-col :lg="3" :md="3"><a-button @click="rotateLeft">左旋</a-button></a-col>
      <a-col :lg="3" :md="3"><a-button @click="rotateRight">右旋</a-button></a-col>
      <a-col :lg="{ span: 4, offset: 4 }" :md="4">
        <a-button type="primary" @click="uploadImg">提交</a-button>
      </a-col>
    </a-row>
  </a-modal>
</template>

<script setup>
import "vue-cropper/dist/index.css";
import { VueCropper } from "vue-cropper";
import { uploadAvatar } from "@/api/system/user";
import { resourceUrl } from "@/utils/resource-url";

defineProps({
  user: {
    type: Object,
    default: () => ({})
  }
});

const store = useStore();
const { proxy } = getCurrentInstance();

const open = ref(false);
const visible = ref(false);
const title = ref("修改头像");

const options = reactive({
  img: store.getters.avatar,
  autoCrop: true,
  autoCropWidth: 200,
  autoCropHeight: 200,
  fixedBox: true,
  previews: {}
});

function editCropper() {
  open.value = true;
}

function handleAfterOpenChange(status) {
  visible.value = status;
  if (!status) {
    closeDialog();
  }
}

function requestUpload({ onSuccess }) {
  onSuccess?.("ok");
}

function rotateLeft() {
  proxy.$refs.cropper.rotateLeft();
}

function rotateRight() {
  proxy.$refs.cropper.rotateRight();
}

function changeScale(num) {
  proxy.$refs.cropper.changeScale(num || 1);
}

function beforeUpload(file) {
  if (file.type.indexOf("image/") === -1) {
    proxy.$modal.msgError("文件格式错误，请上传图片类型文件，如 JPG、PNG。");
    return false;
  }
  const reader = new FileReader();
  reader.readAsDataURL(file);
  reader.onload = () => {
    options.img = reader.result;
  };
  return false;
}

function uploadImg() {
  proxy.$refs.cropper.getCropBlob(data => {
    const formData = new FormData();
    formData.append("avatarfile", data);
    uploadAvatar(formData).then(response => {
      open.value = false;
      options.img = resourceUrl(response.imgUrl);
      store.commit("SET_AVATAR", options.img);
      proxy.$modal.msgSuccess("修改成功");
      visible.value = false;
    });
  });
}

function realTime(data) {
  options.previews = data;
}

function closeDialog() {
  options.img = store.getters.avatar;
  visible.value = false;
}
</script>

<style lang="scss" scoped>
.user-info-head {
  position: relative;
  display: inline-block;
  height: 120px;
}

.user-info-head:hover:after {
  content: "+";
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  color: #eee;
  background: rgba(0, 0, 0, 0.5);
  font-size: 24px;
  font-style: normal;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  cursor: pointer;
  line-height: 110px;
  border-radius: 50%;
}
</style>
