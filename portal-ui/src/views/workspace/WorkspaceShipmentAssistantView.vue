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
          <a-button @click="appendRow">新增一行</a-button>
          <a-button danger :disabled="!selectedRows.length" @click="removeSelected">删除选中</a-button>
          <a-button :loading="estimating" @click="handleEstimate">开始测算</a-button>
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

    <section class="panel-grid top-gap">
      <article class="panel-card">
        <div class="panel-header">
          <div>
            <h3>货物明细</h3>
            <p>可直接点击单元格编辑，支持 Excel 导入后二次整理。</p>
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
          <vxe-column field="weightKg" title="重量(KG)" width="120" :edit-render="{ name: 'input', props: { type: 'number', min: 0, step: '0.01' } }" />
          <vxe-column field="volumeCbm" title="体积(CBM)" width="120" :edit-render="{ name: 'input', props: { type: 'number', min: 0, step: '0.001' } }" />
          <vxe-column field="lengthCm" title="长(cm)" width="100" :edit-render="{ name: 'input', props: { type: 'number', min: 0, step: '0.01' } }" />
          <vxe-column field="widthCm" title="宽(cm)" width="100" :edit-render="{ name: 'input', props: { type: 'number', min: 0, step: '0.01' } }" />
          <vxe-column field="heightCm" title="高(cm)" width="100" :edit-render="{ name: 'input', props: { type: 'number', min: 0, step: '0.01' } }" />
        </vxe-table>
      </article>

      <article class="side-stack">
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
              <h3>测算结果</h3>
              <p>按体积和重量推荐整柜，也给出散货判断。</p>
            </div>
          </div>

          <template v-if="result">
            <div class="result-list">
              <div v-for="container in result.containers" :key="container.containerType" class="result-item">
                <strong>{{ container.containerType }} x {{ container.quantity }}</strong>
                <span>体积 {{ container.usedVolume }}/{{ container.maxVolume }} CBM</span>
                <span>重量 {{ container.usedWeight }}/{{ container.maxWeight }} KG</span>
                <span>装载率 {{ container.loadRate }}%</span>
                <small>{{ container.remark }}</small>
              </div>
            </div>

            <div class="lcl-card" :class="{ recommend: result.lcl.recommended }">
              <strong>{{ result.lcl.recommended ? '建议关注散货方案' : '优先关注整柜方案' }}</strong>
              <span>散货体积 {{ result.lcl.totalVolume }} CBM</span>
              <small>{{ result.lcl.remark }}</small>
            </div>
          </template>

          <a-empty v-else description="导入或录入数据后开始测算" />
        </section>

        <section class="panel-card">
          <div class="panel-header">
            <div>
              <h3>导入说明</h3>
              <p>支持中英文表头，未填写体积时会用长宽高和箱数自动换算。</p>
            </div>
          </div>
          <ul class="tips">
            <li>品名是必填项，空白行会自动忽略。</li>
            <li>支持表头：品名、SKU、箱数、重量、体积、长、宽、高。</li>
            <li>如果 Excel 里只有尺寸和箱数，系统会自动补算 CBM。</li>
          </ul>
        </section>
      </article>
    </section>
  </main>
</template>

<script setup lang="ts">
import * as XLSX from 'xlsx'
import { computed, nextTick, ref } from 'vue'
import type { VxeTableInstance } from 'vxe-table'
import {
  estimateWorkspaceShipment,
  type ShipmentAssistantResult,
  type ShipmentAssistantRow,
} from '@/api/workspace/shipmentAssistant'

type RowInput = ShipmentAssistantRow & { id: string }

const fileInputRef = ref<HTMLInputElement>()
const tableRef = ref<VxeTableInstance<RowInput>>()
const estimating = ref(false)
const errorMessage = ref('')
const preferredType = ref('')
const result = ref<ShipmentAssistantResult>()
const rows = ref<RowInput[]>([createRow()])
const selectedRows = ref<RowInput[]>([])

const headerAliases: Record<keyof ShipmentAssistantRow, string[]> = {
  sku: ['sku', '货号', '商品编码'],
  cargoName: ['品名', '货物名称', '货名', 'cargo name', 'product name'],
  packageType: ['包装', '包装类型', 'package', 'package type'],
  quantity: ['数量', '件数', 'qty', 'quantity'],
  cartons: ['箱数', '箱量', 'ctns', 'cartons'],
  weightKg: ['重量', '毛重', '净重', 'kg', 'weight', 'weightkg'],
  volumeCbm: ['体积', '方数', 'cbm', 'volume', 'volumecbm'],
  lengthCm: ['长', '长cm', 'length', 'lengthcm'],
  widthCm: ['宽', '宽cm', 'width', 'widthcm'],
  heightCm: ['高', '高cm', 'height', 'heightcm'],
}

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

function toNumber(value: unknown) {
  const num = Number(value || 0)
  return Number.isFinite(num) ? num : 0
}

function normalizeHeader(value: unknown) {
  return String(value || '').trim().toLowerCase().replace(/\s+/g, '')
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
    const hit = normalizedEntries.find(([key]) => headerAliases[field].some((alias) => normalizeHeader(alias) === key))
    if (!hit) {
      continue
    }
    const value = hit[1]
    if (typeof row[field] === 'number') {
      ;(row[field] as number) = toNumber(value)
    } else {
      ;(row[field] as string) = String(value || '').trim()
    }
  }

  row.volumeCbm = autoVolume(row)
  return row
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
      cargoList: payloadRows,
    })
    if (response.code !== 200 || !response.data) {
      errorMessage.value = response.msg || '测算失败'
      return
    }
    result.value = response.data
    fillRowsFromResult(response.data.normalizedCargoList)
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '测算失败'
  } finally {
    estimating.value = false
  }
}
</script>

<style scoped>
.assistant-page {
  min-height: 100%;
}

.hero-card,
.panel-card {
  border-radius: 22px;
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
  border-radius: 16px;
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

.result-list {
  display: grid;
  gap: 12px;
}

.result-item,
.lcl-card {
  border-radius: 16px;
  border: 1px solid rgba(148, 163, 184, 0.16);
  background: #f8fafc;
  padding: 16px;
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
}

@media (max-width: 640px) {
  .metric-grid {
    grid-template-columns: 1fr;
  }
}
</style>
