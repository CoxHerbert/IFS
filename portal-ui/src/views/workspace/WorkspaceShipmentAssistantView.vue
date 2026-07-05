<template>
  <main class="assistant-page">
    <section class="hero-card">
      <div>
        <p class="eyebrow">智能出货助手</p>
        <h1>导入 Excel，快速测算整柜和散货</h1>
        <span>支持在线整理货物明细、补体积、估算柜型和散货方数。</span>
      </div>
      <div class="hero-actions">
        <a-space wrap>
          <a-button type="primary" @click="pickFile">导入 Excel</a-button>
          <a-button @click="downloadTemplate">下载模板</a-button>
          <a-button v-if="activeStep === 0" @click="appendRow">新增一行</a-button>
          <a-button v-if="activeStep === 0" danger :disabled="!selectedRows.length" @click="removeSelected">删除选中</a-button>
          <a-button v-if="activeStep === 0" :loading="estimating" @click="handleEstimate">生成出货分析</a-button>
          <a-button v-if="activeStep === 1" @click="activeStep = 0">返回修改</a-button>
          <a-button v-if="activeStep === 1" type="primary" :disabled="!result" @click="downloadReport">下载报告</a-button>
          <a-button v-if="activeStep === 1" :disabled="!result" :loading="creatingPlan" @click="openCreatePlanModal">生成出货计划</a-button>
        </a-space>
        <a-radio-group v-model:value="preferredType" button-style="solid" size="large">
          <a-radio-button value="">自动推荐</a-radio-button>
          <a-radio-button value="20GP">20GP</a-radio-button>
          <a-radio-button value="40GP">40GP</a-radio-button>
          <a-radio-button value="40HQ">40HQ</a-radio-button>
          <a-radio-button value="LCL">散货</a-radio-button>
        </a-radio-group>
      </div>
      <input ref="fileInputRef" type="file" accept=".xlsx,.xls,.csv" class="hidden-input" @change="handleFileChange" />
    </section>

    <a-alert v-if="errorMessage" class="top-gap" type="error" :message="errorMessage" show-icon />

    <section class="step-card top-gap">
      <a-steps :current="activeStep" responsive>
        <a-step title="录入货物明细" description="导入 Excel 或手工维护货物、价格参数" />
        <a-step title="生成出货分析" description="查看推荐方案、成本对比和装柜模拟" />
      </a-steps>
    </section>

    <section v-if="activeStep === 0" class="panel-grid top-gap">
      <article class="panel-card">
        <div class="panel-header">
          <div>
            <h3>货物明细</h3>
            <p>一行代表一种货物；总重量和总体积为该行合计，单箱尺寸用于自动补算体积。</p>
          </div>
          <span>{{ rows.length }} 行</span>
        </div>

        <vxe-table
          ref="tableRef"
          border
          round
          stripe
          height="560"
          :data="rows"
          :checkbox-config="{ highlight: true }"
          :edit-config="{ trigger: 'click', mode: 'cell' }"
          @checkbox-all="syncSelection"
          @checkbox-change="syncSelection"
        >
          <vxe-column type="checkbox" width="56" />
          <vxe-column field="sku" title="SKU" min-width="120" :edit-render="{ name: 'input' }" />
          <vxe-column field="cargoName" title="品名" min-width="180" :edit-render="{ name: 'input' }" />
          <vxe-column field="packageType" title="包装" min-width="100" :edit-render="{ name: 'input' }" />
          <vxe-column field="quantity" title="数量" width="100" :edit-render="{ name: 'input', props: { type: 'number', min: 0 } }" />
          <vxe-column field="cartons" title="箱数" width="100" :edit-render="{ name: 'input', props: { type: 'number', min: 0 } }" />
          <vxe-column field="weightKg" title="总重量(KG)" width="130" :edit-render="{ name: 'input', props: { type: 'number', min: 0, step: '0.01' } }" />
          <vxe-column field="volumeCbm" title="总体积(CBM)" width="130" :edit-render="{ name: 'input', props: { type: 'number', min: 0, step: '0.001' } }" />
          <vxe-column field="lengthCm" title="单箱长(cm)" width="110" :edit-render="{ name: 'input', props: { type: 'number', min: 0, step: '0.01' } }" />
          <vxe-column field="widthCm" title="单箱宽(cm)" width="110" :edit-render="{ name: 'input', props: { type: 'number', min: 0, step: '0.01' } }" />
          <vxe-column field="heightCm" title="单箱高(cm)" width="110" :edit-render="{ name: 'input', props: { type: 'number', min: 0, step: '0.01' } }" />
        </vxe-table>
      </article>

      <article class="side-stack">
        <section class="panel-card">
          <div class="panel-header">
            <div>
              <h3>价格参数</h3>
              <p>用于快速比较散货和整柜成本，可按实际报价调整。</p>
            </div>
          </div>
          <div class="rate-grid">
            <label>
              <span>散货 / CBM</span>
              <a-input-number v-model:value="rates.lclRate" :min="0" :precision="2" />
            </label>
            <label>
              <span>散货最低收费</span>
              <a-input-number v-model:value="rates.lclMinCharge" :min="0" :precision="2" />
            </label>
            <label>
              <span>20GP</span>
              <a-input-number v-model:value="rates.rate20GP" :min="0" :precision="2" />
            </label>
            <label>
              <span>40GP</span>
              <a-input-number v-model:value="rates.rate40GP" :min="0" :precision="2" />
            </label>
            <label>
              <span>40HQ</span>
              <a-input-number v-model:value="rates.rate40HQ" :min="0" :precision="2" />
            </label>
            <label>
              <span>附加费</span>
              <a-input-number v-model:value="rates.extraFees" :min="0" :precision="2" />
            </label>
          </div>
        </section>

        <section class="panel-card">
          <div class="panel-header">
            <div>
              <h3>基础汇总</h3>
              <p>可直接根据当前表格内容实时查看。</p>
            </div>
          </div>
          <div class="metric-grid">
            <div class="metric-item">
              <small>行数</small>
              <strong>{{ localSummary.lineCount }}</strong>
            </div>
            <div class="metric-item">
              <small>总数量</small>
              <strong>{{ localSummary.totalQuantity }}</strong>
            </div>
            <div class="metric-item">
              <small>总箱数</small>
              <strong>{{ localSummary.totalCartons }}</strong>
            </div>
            <div class="metric-item">
              <small>总重量</small>
              <strong>{{ localSummary.totalWeight }} KG</strong>
            </div>
            <div class="metric-item metric-wide">
              <small>总体积</small>
              <strong>{{ localSummary.totalVolume }} CBM</strong>
            </div>
          </div>
        </section>

        <section class="panel-card">
          <div class="panel-header">
            <div>
              <h3>生成出货分析</h3>
              <p>确认货物明细和价格参数后，进入第二步查看报告。</p>
            </div>
          </div>
          <a-button block type="primary" size="large" :loading="estimating" @click="handleEstimate">生成分析报告</a-button>
          <a-button v-if="result" block class="report-shortcut" @click="activeStep = 1">查看上次报告</a-button>
        </section>

        <section class="panel-card">
          <div class="panel-header">
            <div>
              <h3>导入说明</h3>
              <p>Excel 格式不固定也可以导入，系统会按常见表头自动识别。</p>
            </div>
          </div>
          <ul class="tips">
            <li>品名是必填项，空白行会自动忽略。</li>
            <li>总重量、总体积表示该行所有箱子的合计值。</li>
            <li>单箱尺寸支持“长宽高”合并列，例如 60x40x50cm。</li>
            <li>如果只有单箱重量或单箱体积，系统会乘以箱数换算为合计值。</li>
          </ul>
        </section>
      </article>
    </section>

    <section v-if="activeStep === 1" class="report-page top-gap">
      <section class="panel-card">
        <div class="panel-header report-header">
          <div>
            <h3>出货分析报告</h3>
            <p>基于当前货物明细、价格参数、柜型安全容量和装载风险生成。</p>
          </div>
          <a-space wrap>
            <a-button @click="activeStep = 0">返回修改</a-button>
            <a-button type="primary" :disabled="!result" @click="downloadReport">下载报告</a-button>
            <a-button :disabled="!result" :loading="creatingPlan" @click="openCreatePlanModal">生成出货计划</a-button>
          </a-space>
        </div>

        <template v-if="result">
          <div v-if="result.recommendation" class="recommendation-card">
            <small>推荐方案</small>
            <strong>{{ result.recommendation.title }}</strong>
            <span>{{ result.recommendation.reason }}</span>
            <em>置信度 {{ result.recommendation.confidence }} · 风险 {{ result.recommendation.riskLevel }}</em>
          </div>

          <div class="metric-grid report-metrics">
            <div class="metric-item">
              <small>总箱数</small>
              <strong>{{ result.summary.totalCartons }}</strong>
            </div>
            <div class="metric-item">
              <small>总重量</small>
              <strong>{{ result.summary.totalWeight }} KG</strong>
            </div>
            <div class="metric-item">
              <small>总体积</small>
              <strong>{{ result.summary.totalVolume }} CBM</strong>
            </div>
            <div class="metric-item">
              <small>散货成本</small>
              <strong>{{ money(result.lcl.totalCost) }}</strong>
            </div>
          </div>

          <div class="report-section">
            <h4>整柜方案</h4>
            <div class="result-list report-result-grid">
              <div v-for="container in result.containers" :key="container.containerType" class="result-item">
                <strong>{{ container.containerType }} x {{ container.quantity }}</strong>
                <span>安全体积 {{ container.usedVolume }}/{{ container.maxVolume }} CBM</span>
                <span>折算装载体积 {{ container.effectiveVolume }}/{{ container.safeVolume }} CBM</span>
                <span>重量 {{ container.usedWeight }}/{{ container.maxWeight }} KG</span>
                <span>装载率 {{ container.loadRate }}%</span>
                <span>成本 {{ money(container.totalCost) }}</span>
                <span>风险 {{ container.riskLevel }}</span>
                <small>{{ container.remark }}</small>
              </div>
            </div>
          </div>

          <div class="report-section">
            <h4>散货方案</h4>
            <div class="lcl-card" :class="{ recommend: result.lcl.recommended }">
              <strong>{{ result.lcl.recommended ? '建议关注散货方案' : '优先关注整柜方案' }}</strong>
              <span>散货体积 {{ result.lcl.totalVolume }} CBM</span>
              <span>费用 {{ result.lcl.totalVolume }} × {{ money(result.lcl.ratePerCbm) }}，最低 {{ money(result.lcl.minCharge) }}</span>
              <span>散货总成本 {{ money(result.lcl.totalCost) }}</span>
              <small>{{ result.lcl.remark }}</small>
            </div>
          </div>

          <a-alert
            v-if="result.warnings.length"
            class="result-warning"
            type="warning"
            show-icon
            :message="result.warnings[0]"
          />
        </template>

        <a-empty v-else description="请先在第一步生成出货分析" />
      </section>

      <section v-if="createdShipment" class="panel-card created-plan-card">
        <div class="panel-header">
          <div>
            <h3>出货计划已生成</h3>
            <p>计划编号：{{ createdShipment.plan.shipmentNo }}，可前往出货查询继续跟踪。</p>
          </div>
          <router-link to="/customer/shipment">
            <a-button type="primary">查看出货查询</a-button>
          </router-link>
        </div>
      </section>
    </section>

    <section v-if="activeStep === 1 && result?.loadingPlan" class="panel-card top-gap">
      <div class="panel-header">
        <div>
          <h3>模拟装柜视图</h3>
          <p>{{ result.loadingPlan.quantity }} × {{ result.loadingPlan.containerType }}，俯视图用于判断长宽方向空隙，侧视图用于判断堆叠高度。</p>
        </div>
        <span>利用率 {{ result.loadingPlan.utilization }}%</span>
      </div>

      <div class="load-view-grid">
        <div>
          <div class="view-title">俯视图</div>
          <div class="container-top-view" :style="topViewStyle(result.loadingPlan)">
            <button
              v-for="(box, index) in result.loadingPlan.placements"
              :key="`${box.cargoName}-${index}`"
              class="cargo-box"
              :style="topBoxStyle(box, result.loadingPlan)"
              :title="`${box.cargoName} ${box.quantity}件`"
            >
              <span>{{ shortName(box.cargoName) }}</span>
            </button>
          </div>
        </div>

        <div>
          <div class="view-title">侧视图</div>
          <div class="container-side-view" :style="sideViewStyle(result.loadingPlan)">
            <button
              v-for="(box, index) in result.loadingPlan.placements"
              :key="`side-${box.cargoName}-${index}`"
              class="cargo-box side"
              :style="sideBoxStyle(box, result.loadingPlan)"
              :title="`${box.cargoName} ${box.quantity}件，高 ${box.height}cm`"
            >
              <span>{{ box.quantity }}</span>
            </button>
          </div>
        </div>
      </div>

      <div class="legend-list">
        <span v-for="(box, index) in result.loadingPlan.placements.slice(0, 8)" :key="`legend-${index}`">
          <i :style="{ background: box.color }"></i>{{ box.cargoName }} · {{ box.remark }}
        </span>
      </div>
    </section>

    <a-modal
      v-model:visible="createPlanVisible"
      title="根据报告生成出货计划"
      :confirm-loading="creatingPlan"
      ok-text="生成出货计划"
      cancel-text="取消"
      @ok="handleCreatePlan"
    >
      <a-form layout="vertical">
        <a-form-item label="客户订单号">
          <a-input v-model:value="planForm.orderNo" placeholder="可选，如 PO/订单号" />
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="起运港">
              <a-input v-model:value="planForm.pol" placeholder="例如 Shanghai" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="目的港">
              <a-input v-model:value="planForm.pod" placeholder="例如 Los Angeles" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="计划开船">
              <a-date-picker v-model:value="planForm.plannedEtd" value-format="YYYY-MM-DD" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="计划到港">
              <a-date-picker v-model:value="planForm.plannedEta" value-format="YYYY-MM-DD" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="备注">
          <a-textarea v-model:value="planForm.remark" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>
  </main>
</template>

<script setup lang="ts">
import * as XLSX from 'xlsx'
import { computed, nextTick, ref } from 'vue'
import type { VxeTableInstance } from 'vxe-table'
import { message as antMessage } from 'ant-design-vue'
import { parseDimensions } from '@/utils/dimensionParser'
import {
  createWorkspaceShipmentPlan,
  estimateWorkspaceShipment,
  type ShipmentPlanCreated,
  type ShipmentAssistantResult,
  type ShipmentAssistantRow,
  type LoadingPlacement,
  type LoadingPlan,
} from '@/api/workspace/shipmentAssistant'

type RowInput = ShipmentAssistantRow & { id: string }

const fileInputRef = ref<HTMLInputElement>()
const tableRef = ref<VxeTableInstance<RowInput>>()
const activeStep = ref(0)
const estimating = ref(false)
const creatingPlan = ref(false)
const createPlanVisible = ref(false)
const errorMessage = ref('')
const preferredType = ref('')
const result = ref<ShipmentAssistantResult>()
const createdShipment = ref<ShipmentPlanCreated>()
const rows = ref<RowInput[]>([createRow()])
const selectedRows = ref<RowInput[]>([])
const planForm = ref({
  orderNo: '',
  pol: '',
  pod: '',
  plannedEtd: '',
  plannedEta: '',
  remark: '',
})
const rates = ref({
  lclRate: 85,
  lclMinCharge: 300,
  rate20GP: 1800,
  rate40GP: 2600,
  rate40HQ: 2800,
  extraFees: 120,
})

const headerAliases: Record<keyof ShipmentAssistantRow, string[]> = {
  sku: ['sku', '货号', '商品编码'],
  cargoName: ['品名', '货物名称', '货名', 'cargo name', 'product name'],
  packageType: ['包装', '包装类型', 'package', 'package type'],
  quantity: ['数量', '件数', 'qty', 'quantity'],
  cartons: ['箱数', '箱量', 'ctns', 'cartons'],
  weightKg: ['总重量', '总毛重', '合计重量', '重量', '毛重', '净重', 'kg', 'totalweight', 'grossweight', 'weightkg'],
  volumeCbm: ['总体积', '总方数', '合计体积', '体积', '方数', 'cbm', 'totalcbm', 'volume', 'volumecbm'],
  lengthCm: ['单箱长', '长', '长cm', 'length', 'lengthcm'],
  widthCm: ['单箱宽', '宽', '宽cm', 'width', 'widthcm'],
  heightCm: ['单箱高', '高', '高cm', 'height', 'heightcm'],
}

const dimensionAliases = ['尺寸', '规格', '长宽高', '单箱尺寸', '包装尺寸', '外箱尺寸', 'lwh', 'dimension', 'dimensions']
const unitWeightAliases = ['单箱重量', '单件重量', '单重', '每箱重量', 'unitweight', 'weightpercarton']
const unitVolumeAliases = ['单箱体积', '单件体积', '单方', '每箱体积', 'unitcbm', 'cbmpercarton']

const localSummary = computed(() => {
  return rows.value.reduce(
    (summary, row) => {
      const volume = row.volumeCbm || autoVolume(row)
      if (!row.cargoName.trim()) {
        return summary
      }
      summary.lineCount += 1
      summary.totalQuantity += toNumber(row.quantity)
      summary.totalCartons += toNumber(row.cartons)
      summary.totalWeight = round2(summary.totalWeight + toNumber(row.weightKg))
      summary.totalVolume = round2(summary.totalVolume + volume)
      return summary
    },
    { lineCount: 0, totalQuantity: 0, totalCartons: 0, totalWeight: 0, totalVolume: 0 },
  )
})

function createRow(): RowInput {
  return {
    id: `${Date.now()}-${Math.random().toString(16).slice(2)}`,
    sku: '',
    cargoName: '',
    packageType: '',
    quantity: 0,
    cartons: 0,
    weightKg: 0,
    volumeCbm: 0,
    lengthCm: 0,
    widthCm: 0,
    heightCm: 0,
  }
}

function round2(value: number) {
  return Math.round(value * 100) / 100
}

function money(value: number) {
  return `$${round2(value).toLocaleString()}`
}

function toNumber(value: unknown) {
  const num = Number(value || 0)
  return Number.isFinite(num) ? num : 0
}

function normalizeHeader(value: unknown) {
  return String(value || '').trim().toLowerCase().replace(/[\s_/（）().-]+/g, '')
}

function autoVolume(row: ShipmentAssistantRow) {
  if (row.volumeCbm > 0) {
    return round2(row.volumeCbm)
  }
  if (row.lengthCm > 0 && row.widthCm > 0 && row.heightCm > 0 && row.cartons > 0) {
    return round2((row.lengthCm * row.widthCm * row.heightCm * row.cartons) / 1000000)
  }
  return 0
}

function normalizeImportedRow(raw: Record<string, unknown>): RowInput {
  const row = createRow()
  const normalizedEntries = Object.entries(raw).map(([key, value]) => [normalizeHeader(key), value] as const)

  for (const field of Object.keys(headerAliases) as Array<keyof ShipmentAssistantRow>) {
    const hit = findImportValue(normalizedEntries, headerAliases[field])
    if (!hit) {
      continue
    }
    const value = hit
    if (typeof row[field] === 'number') {
      ;(row[field] as number) = toNumber(value)
    } else {
      ;(row[field] as string) = String(value || '').trim()
    }
  }

  const dimension = parseDimensions(findImportValue(normalizedEntries, dimensionAliases))
  if (dimension && (!row.lengthCm || !row.widthCm || !row.heightCm)) {
    row.lengthCm = dimension.length
    row.widthCm = dimension.width
    row.heightCm = dimension.height
  }

  const unitWeight = toNumber(findImportValue(normalizedEntries, unitWeightAliases))
  if (!row.weightKg && unitWeight > 0 && row.cartons > 0) {
    row.weightKg = round2(unitWeight * row.cartons)
  }

  const unitVolume = toNumber(findImportValue(normalizedEntries, unitVolumeAliases))
  if (!row.volumeCbm && unitVolume > 0 && row.cartons > 0) {
    row.volumeCbm = round2(unitVolume * row.cartons)
  }

  row.volumeCbm = autoVolume(row)
  return row
}

function findImportValue(entries: Array<readonly [string, unknown]>, aliases: string[]) {
  const aliasSet = aliases.map(normalizeHeader)
  const exact = entries.find(([key]) => aliasSet.includes(key))
  if (exact) {
    return exact[1]
  }
  const fuzzy = entries.find(([key]) => aliasSet.some((alias) => key.includes(alias) || alias.includes(key)))
  return fuzzy?.[1]
}

function downloadTemplate() {
  const data = [
    {
      SKU: 'A001',
      品名: '示例纸箱',
      包装: '纸箱',
      数量: 120,
      箱数: 120,
      '总重量(KG)': 960,
      '总体积(CBM)': '',
      '单箱尺寸(cm)': '60x40x50',
      说明: '总体积为空时，将按单箱尺寸×箱数自动计算',
    },
  ]
  const worksheet = XLSX.utils.json_to_sheet(data)
  const workbook = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(workbook, worksheet, '货物明细')
  XLSX.writeFile(workbook, '智能出货助手导入模板.xlsx')
}

function pickFile() {
  fileInputRef.value?.click()
}

function appendRow() {
  rows.value.push(createRow())
  nextTick(() => {
    tableRef.value?.setActiveCell(rows.value[rows.value.length - 1], 'cargoName')
  })
}

function syncSelection() {
  selectedRows.value = tableRef.value?.getCheckboxRecords() || []
}

function removeSelected() {
  const selectedIds = new Set(selectedRows.value.map((item) => item.id))
  rows.value = rows.value.filter((item) => !selectedIds.has(item.id))
  if (!rows.value.length) {
    rows.value = [createRow()]
  }
  selectedRows.value = []
}

function fillRowsFromResult(resultRows: ShipmentAssistantRow[]) {
  rows.value = resultRows.map((item) => ({
    id: `${Date.now()}-${Math.random().toString(16).slice(2)}`,
    ...item,
  }))
}

async function handleFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) {
    return
  }

  errorMessage.value = ''
  try {
    const buffer = await file.arrayBuffer()
    const workbook = XLSX.read(buffer, { type: 'array' })
    const sheetName = workbook.SheetNames[0]
    if (!sheetName) {
      throw new Error('Excel 中没有可读取的工作表')
    }
    const sheet = workbook.Sheets[sheetName]
    const jsonRows = XLSX.utils.sheet_to_json<Record<string, unknown>>(sheet, { defval: '' })
    const importedRows = jsonRows.map(normalizeImportedRow).filter((item) => item.cargoName.trim())
    if (!importedRows.length) {
      throw new Error('没有识别到有效的货物明细，请检查表头')
    }
    rows.value = importedRows
    selectedRows.value = []
    result.value = undefined
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Excel 导入失败'
  } finally {
    input.value = ''
  }
}

async function handleEstimate() {
  errorMessage.value = ''
  result.value = undefined
  const payloadRows = rows.value
    .filter((item) => item.cargoName.trim())
    .map((item) => ({
      sku: item.sku.trim(),
      cargoName: item.cargoName.trim(),
      packageType: item.packageType.trim(),
      quantity: toNumber(item.quantity),
      cartons: toNumber(item.cartons),
      weightKg: toNumber(item.weightKg),
      volumeCbm: autoVolume(item),
      lengthCm: toNumber(item.lengthCm),
      widthCm: toNumber(item.widthCm),
      heightCm: toNumber(item.heightCm),
    }))

  if (!payloadRows.length) {
    errorMessage.value = '请先录入至少一条货物明细'
    return
  }

  estimating.value = true
  try {
    const response = await estimateWorkspaceShipment({
      preferredType: preferredType.value,
      ...rates.value,
      cargoList: payloadRows,
    })
    if (response.code !== 200 || !response.data) {
      errorMessage.value = response.msg || '测算失败'
      return
    }
    result.value = response.data
    fillRowsFromResult(response.data.normalizedCargoList)
    activeStep.value = 1
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '测算失败'
  } finally {
    estimating.value = false
  }
}

function shortName(value: string) {
  const text = value || '货物'
  return text.length > 5 ? `${text.slice(0, 5)}…` : text
}

function topViewStyle(plan: LoadingPlan) {
  return {
    aspectRatio: `${plan.internalLengthCm} / ${plan.internalWidthCm}`,
  }
}

function sideViewStyle(plan: LoadingPlan) {
  return {
    aspectRatio: `${plan.internalLengthCm} / ${plan.internalHeightCm}`,
  }
}

function topBoxStyle(box: LoadingPlacement, plan: LoadingPlan) {
  return {
    left: `${(box.x / plan.internalLengthCm) * 100}%`,
    top: `${(box.y / plan.internalWidthCm) * 100}%`,
    width: `${Math.max((box.length / plan.internalLengthCm) * 100, 1.2)}%`,
    height: `${Math.max((box.width / plan.internalWidthCm) * 100, 3)}%`,
    background: box.color,
  }
}

function sideBoxStyle(box: LoadingPlacement, plan: LoadingPlan) {
  return {
    left: `${(box.x / plan.internalLengthCm) * 100}%`,
    bottom: '0%',
    width: `${Math.max((box.length / plan.internalLengthCm) * 100, 1.2)}%`,
    height: `${Math.max((box.height / plan.internalHeightCm) * 100, 4)}%`,
    background: box.color,
  }
}

function downloadReport() {
  if (!result.value) {
    return
  }
  const html = buildReportHtml(result.value)
  const blob = new Blob([html], { type: 'text/html;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `出货分析报告-${new Date().toISOString().slice(0, 10)}.html`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

function openCreatePlanModal() {
  if (!result.value) {
    return
  }
  if (!planForm.value.remark) {
    planForm.value.remark = buildDefaultPlanRemark()
  }
  createPlanVisible.value = true
}

async function handleCreatePlan() {
  if (!result.value) {
    return
  }
  creatingPlan.value = true
  errorMessage.value = ''
  try {
    const response = await createWorkspaceShipmentPlan({
      ...planForm.value,
      preferredType: preferredType.value || (result.value.recommendation?.mode === 'LCL' ? 'LCL' : ''),
      cargoList: result.value.normalizedCargoList,
    })
    if (response.code !== 200 || !response.data) {
      errorMessage.value = response.msg || '生成出货计划失败'
      return
    }
    createdShipment.value = response.data
    createPlanVisible.value = false
    antMessage.success(`出货计划已生成：${response.data.plan.shipmentNo}`)
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '生成出货计划失败'
  } finally {
    creatingPlan.value = false
  }
}

function buildDefaultPlanRemark() {
  const recommendation = result.value?.recommendation
  const title = recommendation?.title || '系统自动生成出货分析'
  return `由智能出货助手根据分析报告生成；推荐方案：${title}`
}

function buildReportHtml(report: ShipmentAssistantResult) {
  const containerRows = report.containers
    .map(
      (item) => `
        <tr>
          <td>${escapeHtml(item.containerType)} x ${item.quantity}</td>
          <td>${item.usedVolume}/${item.maxVolume} CBM</td>
          <td>${item.effectiveVolume}/${item.safeVolume} CBM</td>
          <td>${item.usedWeight}/${item.maxWeight} KG</td>
          <td>${item.loadRate}%</td>
          <td>${money(item.totalCost)}</td>
          <td>${escapeHtml(item.riskLevel)}</td>
        </tr>
      `,
    )
    .join('')
  const cargoRows = report.normalizedCargoList
    .map(
      (item) => `
        <tr>
          <td>${escapeHtml(item.sku)}</td>
          <td>${escapeHtml(item.cargoName)}</td>
          <td>${escapeHtml(item.packageType)}</td>
          <td>${item.cartons}</td>
          <td>${item.weightKg}</td>
          <td>${item.volumeCbm}</td>
          <td>${item.lengthCm} x ${item.widthCm} x ${item.heightCm}</td>
        </tr>
      `,
    )
    .join('')
  const warnings = report.warnings.length ? report.warnings.map((item) => `<li>${escapeHtml(item)}</li>`).join('') : '<li>暂无明显尺寸或装载风险提示。</li>'
  const recommendation = report.recommendation
  return `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8" />
  <title>出货分析报告</title>
  <style>
    body { margin: 0; padding: 32px; font-family: Arial, "Microsoft YaHei", sans-serif; color: #111827; background: #f8fafc; }
    main { max-width: 1080px; margin: 0 auto; background: #fff; padding: 32px; border: 1px solid #e2e8f0; }
    h1, h2 { margin: 0; }
    h1 { font-size: 28px; }
    h2 { margin-top: 28px; font-size: 18px; }
    p { color: #475569; line-height: 1.8; }
    .summary { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; margin-top: 20px; }
    .metric { border: 1px solid #e2e8f0; padding: 14px; background: #f8fafc; }
    .metric small { display: block; color: #64748b; }
    .metric strong { display: block; margin-top: 8px; font-size: 20px; }
    .recommend { margin-top: 20px; border: 1px solid #99f6e4; background: #ecfeff; padding: 16px; }
    table { width: 100%; border-collapse: collapse; margin-top: 12px; }
    th, td { border: 1px solid #e2e8f0; padding: 10px; text-align: left; font-size: 13px; }
    th { background: #f1f5f9; }
    ul { color: #475569; line-height: 1.8; }
    @media print { body { background: #fff; padding: 0; } main { border: 0; } }
  </style>
</head>
<body>
  <main>
    <h1>出货分析报告</h1>
    <p>生成时间：${new Date().toLocaleString()}</p>
    <div class="recommend">
      <strong>${escapeHtml(recommendation?.title || '暂无推荐方案')}</strong>
      <p>${escapeHtml(recommendation?.reason || '')} 置信度：${escapeHtml(recommendation?.confidence || '-')}；风险：${escapeHtml(recommendation?.riskLevel || '-')}</p>
    </div>
    <section class="summary">
      <div class="metric"><small>总箱数</small><strong>${report.summary.totalCartons}</strong></div>
      <div class="metric"><small>总重量</small><strong>${report.summary.totalWeight} KG</strong></div>
      <div class="metric"><small>总体积</small><strong>${report.summary.totalVolume} CBM</strong></div>
      <div class="metric"><small>散货成本</small><strong>${money(report.lcl.totalCost)}</strong></div>
    </section>
    <h2>整柜方案</h2>
    <table>
      <thead><tr><th>柜型</th><th>安全体积</th><th>折算装载体积</th><th>重量</th><th>装载率</th><th>成本</th><th>风险</th></tr></thead>
      <tbody>${containerRows || '<tr><td colspan="7">暂无整柜方案</td></tr>'}</tbody>
    </table>
    <h2>散货方案</h2>
    <p>散货体积 ${report.lcl.totalVolume} CBM；费用 ${report.lcl.totalVolume} × ${money(report.lcl.ratePerCbm)}，最低 ${money(report.lcl.minCharge)}；散货总成本 ${money(report.lcl.totalCost)}。</p>
    <h2>风险提示</h2>
    <ul>${warnings}</ul>
    <h2>货物明细</h2>
    <table>
      <thead><tr><th>SKU</th><th>品名</th><th>包装</th><th>箱数</th><th>总重量(KG)</th><th>总体积(CBM)</th><th>单箱尺寸(cm)</th></tr></thead>
      <tbody>${cargoRows}</tbody>
    </table>
  </main>
</body>
</html>`
}

function escapeHtml(value: unknown) {
  return String(value ?? '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;')
}
</script>

<style scoped>
.assistant-page {
  min-height: 100%;
}

.hero-card,
.panel-card {
  border-radius: 8px;
  border: 1px solid rgba(15, 23, 42, 0.08);
  box-shadow: 0 16px 36px rgba(15, 23, 42, 0.06);
  background: #fff;
}

.hero-card {
  padding: 28px;
  display: grid;
  gap: 18px;
}

.eyebrow,
.metric-item small {
  margin: 0;
  font-size: 11px;
  letter-spacing: 0.08em;
  color: #64748b;
}

.hero-card h1 {
  margin: 12px 0 0;
  font-size: clamp(30px, 4vw, 40px);
  color: #111827;
}

.hero-card span {
  display: block;
  margin-top: 12px;
  color: #475569;
}

.hero-actions {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}

.hidden-input {
  display: none;
}

.top-gap {
  margin-top: 18px;
}

.step-card {
  border-radius: 8px;
  border: 1px solid rgba(15, 23, 42, 0.08);
  background: #fff;
  padding: 22px;
  box-shadow: 0 10px 28px rgba(15, 23, 42, 0.05);
}

.panel-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.65fr) minmax(320px, 0.95fr);
  gap: 18px;
}

.panel-card {
  padding: 22px;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
  align-items: flex-start;
}

.panel-header h3 {
  margin: 0;
  font-size: 20px;
  color: #111827;
}

.panel-header p,
.tips {
  margin: 8px 0 0;
  color: #64748b;
  line-height: 1.8;
}

.side-stack {
  display: grid;
  gap: 18px;
}

.metric-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.metric-item {
  border-radius: 8px;
  background: #f8fafc;
  padding: 16px;
  border: 1px solid rgba(148, 163, 184, 0.14);
}

.metric-item strong {
  display: block;
  margin-top: 8px;
  font-size: 24px;
  color: #111827;
}

.metric-wide {
  grid-column: 1 / -1;
}

.report-page {
  display: grid;
  gap: 18px;
}

.report-header {
  align-items: center;
}

.report-metrics {
  grid-template-columns: repeat(4, minmax(0, 1fr));
  margin-top: 16px;
}

.report-section {
  margin-top: 22px;
}

.report-section h4 {
  margin: 0 0 12px;
  font-size: 18px;
  color: #111827;
}

.report-result-grid {
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
}

.result-list {
  display: grid;
  gap: 12px;
}

.result-item,
.lcl-card,
.recommendation-card {
  border-radius: 8px;
  border: 1px solid rgba(148, 163, 184, 0.16);
  background: #f8fafc;
  padding: 16px;
}

.recommendation-card {
  margin-bottom: 14px;
  background: #ecfeff;
  border-color: rgba(8, 145, 178, 0.24);
}

.recommendation-card small,
.recommendation-card strong,
.recommendation-card span,
.recommendation-card em {
  display: block;
}

.recommendation-card small {
  color: #0f766e;
  font-weight: 700;
}

.recommendation-card strong {
  margin-top: 6px;
  font-size: 22px;
  color: #111827;
}

.recommendation-card span,
.recommendation-card em {
  margin-top: 8px;
  color: #475569;
  font-style: normal;
}

.result-item strong,
.result-item span,
.result-item small,
.lcl-card strong,
.lcl-card span,
.lcl-card small {
  display: block;
}

.result-item span,
.result-item small,
.lcl-card span,
.lcl-card small {
  margin-top: 8px;
  color: #64748b;
}

.lcl-card {
  margin-top: 14px;
}

.lcl-card.recommend {
  background: #eff6ff;
  border-color: rgba(59, 130, 246, 0.2);
}

.result-warning {
  margin-top: 14px;
}

.report-shortcut {
  margin-top: 10px;
}

.rate-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.rate-grid label {
  display: grid;
  gap: 6px;
  color: #475569;
  font-size: 13px;
}

.rate-grid :deep(.ant-input-number) {
  width: 100%;
}

.load-view-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.35fr) minmax(0, 1fr);
  gap: 18px;
}

.view-title {
  margin-bottom: 8px;
  font-weight: 700;
  color: #111827;
}

.container-top-view,
.container-side-view {
  position: relative;
  width: 100%;
  min-height: 180px;
  overflow: hidden;
  border: 2px solid #334155;
  border-radius: 8px;
  background:
    linear-gradient(90deg, rgba(148, 163, 184, 0.12) 1px, transparent 1px),
    linear-gradient(0deg, rgba(148, 163, 184, 0.12) 1px, transparent 1px),
    #f8fafc;
  background-size: 32px 32px;
}

.container-side-view {
  min-height: 150px;
}

.cargo-box {
  position: absolute;
  border: 1px solid rgba(15, 23, 42, 0.2);
  border-radius: 4px;
  color: #fff;
  font-size: 11px;
  line-height: 1.1;
  overflow: hidden;
  cursor: default;
  box-shadow: inset 0 -10px 18px rgba(15, 23, 42, 0.12);
}

.cargo-box span {
  display: block;
  padding: 3px;
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
}

.cargo-box.side span {
  text-align: center;
}

.legend-list {
  display: flex;
  gap: 10px 16px;
  flex-wrap: wrap;
  margin-top: 16px;
  color: #475569;
}

.legend-list span {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
}

.legend-list i {
  width: 10px;
  height: 10px;
  border-radius: 2px;
}

.tips {
  padding-left: 18px;
}

:deep(.vxe-table) {
  --vxe-ui-font-color: #111827;
  --vxe-ui-layout-background-color: #fff;
  --vxe-ui-table-row-hover-background-color: #f8fafc;
  --vxe-ui-table-header-background-color: #f8fafc;
  --vxe-ui-table-border-color: rgba(148, 163, 184, 0.24);
  --vxe-ui-primary-color: #111827;
}

@media (max-width: 1200px) {
  .panel-grid {
    grid-template-columns: 1fr;
  }

  .report-metrics {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .load-view-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 640px) {
  .metric-grid,
  .report-metrics {
    grid-template-columns: 1fr;
  }

  .rate-grid {
    grid-template-columns: 1fr;
  }
}
</style>
