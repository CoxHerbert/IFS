<template><AuthForm mode="register" :inviter-uid="inviterUid" /></template>
<script setup lang="ts">
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import AuthForm from '../../components/AuthForm.vue'

const inviterUid = ref('')
onLoad(options => { inviterUid.value = String(options?.inviter || '').trim().toUpperCase() })
</script>
