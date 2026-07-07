<template>
  <div class="app-container">
    <a-row :gutter="16">
      <a-col :span="12" class="card-box">
        <a-card title="CPU">
          <table class="info-table">
            <thead>
              <tr>
                <th>属性</th>
                <th>值</th>
              </tr>
            </thead>
            <tbody>
              <tr><td>核心数</td><td>{{ server.cpuNum }}</td></tr>
              <tr><td>线程数</td><td>{{ server.cpuNumThread }}</td></tr>
              <tr><td>用户使用率</td><td>{{ server.cpuUsed }}%</td></tr>
              <tr><td>系统使用率</td><td>{{ server.cpuAvg5 }}%</td></tr>
              <tr><td>当前空闲率</td><td>{{ server.cpuAvg15 }}%</td></tr>
            </tbody>
          </table>
        </a-card>
      </a-col>

      <a-col :span="12" class="card-box">
        <a-card title="内存">
          <table class="info-table">
            <thead>
              <tr>
                <th>属性</th>
                <th>内存</th>
              </tr>
            </thead>
            <tbody>
              <tr><td>总内存</td><td>{{ server.memTotal }}G</td></tr>
              <tr><td>已用内存</td><td>{{ server.memUsed }}G</td></tr>
              <tr><td>GO 使用内存</td><td>{{ server.goUsed }}G</td></tr>
              <tr><td>剩余内存</td><td>{{ server.memFree }}G</td></tr>
              <tr><td>使用率</td><td :class="{ 'text-danger': server.memUsage > 80 }">{{ server.memUsage }}%</td></tr>
            </tbody>
          </table>
        </a-card>
      </a-col>

      <a-col :span="24" class="card-box">
        <a-card title="服务器信息">
          <table class="info-table">
            <tbody>
              <tr>
                <td>服务器名称</td><td>{{ server.sysComputerName }}</td>
                <td>操作系统</td><td>{{ server.sysOsName }}</td>
              </tr>
              <tr>
                <td>服务器 IP</td><td>{{ server.sysComputerIp }}</td>
                <td>系统架构</td><td>{{ server.sysOsArch }}</td>
              </tr>
              <tr>
                <td>语言环境</td><td>{{ server.goName }}</td>
                <td>环境版本</td><td>{{ server.goVersion }}</td>
              </tr>
              <tr>
                <td>启动时间</td><td>{{ server.goStartTime }}</td>
                <td>运行时长</td><td>{{ server.goRunTime }}</td>
              </tr>
              <tr>
                <td>安装路径</td><td colspan="3">{{ server.goHome }}</td>
              </tr>
              <tr>
                <td>项目路径</td><td colspan="3">{{ server.goUserDir }}</td>
              </tr>
            </tbody>
          </table>
        </a-card>
      </a-col>

      <a-col :span="24" class="card-box">
        <a-card title="磁盘状态">
          <table class="info-table">
            <thead>
              <tr>
                <th>盘符路径</th>
                <th>总大小</th>
                <th>可用大小</th>
                <th>已用大小</th>
                <th>已用百分比</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="sysFile in server.diskList || []" :key="sysFile.path">
                <td>{{ sysFile.path }}</td>
                <td>{{ sysFile.total }}</td>
                <td>{{ sysFile.free }}</td>
                <td>{{ sysFile.used }}</td>
                <td :class="{ 'text-danger': sysFile.usedPercent > 80 }">{{ sysFile.usedPercent }}%</td>
              </tr>
            </tbody>
          </table>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { getServer } from "@/api/monitor/server";

const server = ref({});
const { proxy } = getCurrentInstance();

function getList() {
  proxy.$modal.loading("正在加载服务监控数据，请稍候！");
  getServer().then(response => {
    server.value = response.data;
    proxy.$modal.closeLoading();
  });
}

getList();
</script>

<style scoped>
.info-table {
  width: 100%;
  border-collapse: collapse;
}

.info-table th,
.info-table td {
  padding: 12px 16px;
  border: 1px solid #f0f0f0;
}
</style>
