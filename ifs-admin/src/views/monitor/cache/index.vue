<template>
  <div class="app-container">
    <a-row :gutter="16">
      <a-col :span="24" class="card-box">
        <a-card title="基本信息">
          <div class="info-table-wrap">
            <table class="info-table">
              <tbody>
                <tr>
                  <td>Redis 版本</td>
                  <td>{{ cache.info?.redis_version }}</td>
                  <td>运行模式</td>
                  <td>{{ cache.info?.redis_mode === "standalone" ? "单机" : "集群" }}</td>
                  <td>端口</td>
                  <td>{{ cache.info?.tcp_port }}</td>
                  <td>客户端数</td>
                  <td>{{ cache.info?.connected_clients }}</td>
                </tr>
                <tr>
                  <td>运行时长(天)</td>
                  <td>{{ cache.info?.uptime_in_days }}</td>
                  <td>使用内存</td>
                  <td>{{ cache.info?.used_memory_human }}</td>
                  <td>使用 CPU</td>
                  <td>{{ Number.parseFloat(cache.info?.used_cpu_user_children || 0).toFixed(2) }}</td>
                  <td>内存配置</td>
                  <td>{{ cache.info?.maxmemory_human }}</td>
                </tr>
                <tr>
                  <td>AOF 是否开启</td>
                  <td>{{ cache.info?.aof_enabled === "0" ? "否" : "是" }}</td>
                  <td>RDB 是否成功</td>
                  <td>{{ cache.info?.rdb_last_bgsave_status }}</td>
                  <td>Key 数量</td>
                  <td>{{ cache.dbSize }}</td>
                  <td>网络入口/出口</td>
                  <td>{{ cache.info?.instantaneous_input_kbps }}kps/{{ cache.info?.instantaneous_output_kbps }}kps</td>
                </tr>
              </tbody>
            </table>
          </div>
        </a-card>
      </a-col>

      <a-col :span="12" class="card-box">
        <a-card title="命令统计">
          <div ref="commandstats" style="height: 420px" />
        </a-card>
      </a-col>

      <a-col :span="12" class="card-box">
        <a-card title="内存信息">
          <div ref="usedmemory" style="height: 420px" />
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { getCache } from "@/api/monitor/cache";
import * as echarts from "echarts";

const cache = ref({});
const commandstats = ref(null);
const usedmemory = ref(null);
const { proxy } = getCurrentInstance();

function getList() {
  proxy.$modal.loading("正在加载缓存监控数据，请稍候！");
  getCache().then(response => {
    proxy.$modal.closeLoading();
    cache.value = response.data;

    const commandstatsIntance = echarts.init(commandstats.value, "macarons");
    commandstatsIntance.setOption({
      tooltip: {
        trigger: "item",
        formatter: "{a} <br/>{b} : {c} ({d}%)"
      },
      series: [
        {
          name: "命令",
          type: "pie",
          roseType: "radius",
          radius: [15, 95],
          center: ["50%", "38%"],
          data: response.data.commandStats,
          animationEasing: "cubicInOut",
          animationDuration: 1000
        }
      ]
    });

    const usedmemoryInstance = echarts.init(usedmemory.value, "macarons");
    usedmemoryInstance.setOption({
      tooltip: {
        formatter: `{b} <br/>{a} : ${cache.value.info.used_memory_human}`
      },
      series: [
        {
          name: "峰值",
          type: "gauge",
          min: 0,
          max: 1000,
          detail: {
            formatter: cache.value.info.used_memory_human
          },
          data: [
            {
              value: Number.parseFloat(cache.value.info.used_memory_human),
              name: "内存消耗"
            }
          ]
        }
      ]
    });
  });
}

getList();
</script>

<style scoped>
.info-table-wrap {
  overflow-x: auto;
}

.info-table {
  width: 100%;
  border-collapse: collapse;
}

.info-table td {
  padding: 12px 16px;
  border: 1px solid #f0f0f0;
}
</style>
